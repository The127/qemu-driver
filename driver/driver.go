package driver

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gwenya/qemu-driver/cmdBuilder"
	"github.com/gwenya/qemu-driver/devices"
	"github.com/gwenya/qemu-driver/devices/chardev"
	"github.com/gwenya/qemu-driver/devices/pcie"
	"github.com/gwenya/qemu-driver/devices/storage"
	"github.com/gwenya/qemu-driver/machine"
	"github.com/gwenya/qemu-driver/qmp"
	"github.com/gwenya/qemu-driver/util"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

type StorageFilename string
type RuntimeFilename string

const (
	QemuPidFileName     StorageFilename = "qemu.pid"
	RootDiskFileName    StorageFilename = "rootdisk.img"
	ConfigFileName      StorageFilename = "qemu.conf"
	QemuStdErrFileName  StorageFilename = "stderr.log"
	QemuStdOutFileName  StorageFilename = "stdout.log"
	FirmwareFileName    StorageFilename = "firmware.fd"
	NvramFileName       StorageFilename = "nvram.fd"
	CloudInitIsoFile    StorageFilename = "cloud-init.iso"
	CreatedFlagFileName StorageFilename = "created"
)

const (
	QemuQmpSocketFileName RuntimeFilename = "qmp.sock"
	ConsoleSocketFileName RuntimeFilename = "console.sock"
)

type Driver interface {
	Create() error
	Start() error
	Stop() error
	Reboot() error
	GetState() Status
	Scale(cpuCount uint32, memory uint64, disk uint64) error
	AttachNetworkInterface(net NetworkInterface) error
	DetachNetworkInterface(id uuid.UUID) error
	AttachVolume(volume Volume) error
	DetachVolume(id uuid.UUID) error
}

type driver struct {
	config    MachineConfiguration
	qemuPath  string
	mu        sync.Mutex
	qemuPidFd int
	mon       qmp.Monitor
}

func New(qemuPath string, config MachineConfiguration) (Driver, error) {
	d := &driver{
		config:    config,
		qemuPath:  qemuPath,
		qemuPidFd: -1,
	}

	pidFilePath := d.storagePath(QemuPidFileName)
	pidFileBytes, err := os.ReadFile(pidFilePath)
	if os.IsNotExist(err) {
		return d, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to open pid file %q: %w", pidFilePath, err)
	}

	pidInt64, err := strconv.ParseInt(strings.TrimSpace(string(pidFileBytes)), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pid file content %q: %w", string(pidFileBytes), err)
	}

	pid := int(pidInt64)

	pidfd, err := unix.PidfdOpen(pid, unix.PIDFD_NONBLOCK)
	if errors.Is(err, unix.ESRCH) {
		err := os.Remove(pidFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove old pid file %q: %w", pidFilePath, err)
		}

		return d, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to open process handle for pid %d: %w", pid, err)
	}

	cmdlineBytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return nil, fmt.Errorf("failed to read cmdline for pid %d: %w", pid, err)
	}

	cmdline := string(cmdlineBytes)

	if !strings.Contains(cmdline, config.Id.String()) {
		err := os.Remove(pidFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove old pid file %q: %w", pidFilePath, err)
		}

		return d, nil
	}

	// check if the process is still running, to prevent TOCTOU on the cmdline compare
	err = unix.PidfdSendSignal(pidfd, unix.Signal(0), nil, 0)
	if errors.Is(err, unix.ENOENT) {
		err := os.Remove(pidFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove old pid file %q: %w", pidFilePath, err)
		}

		return d, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to signal qemu process: %w", err)
	}

	d.qemuPidFd = pidfd
	err = d.startWatcher()
	if err != nil {
		return nil, fmt.Errorf("creating qemu process watcher: %w", err)
	}

	return d, nil
}

func (d *driver) startWatcher() error {
	pidfd := d.qemuPidFd

	epfd, err := syscall.EpollCreate(1)
	if err != nil {
		return fmt.Errorf("creating epoll: %w", err)
	}

	err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, pidfd, &syscall.EpollEvent{
		Events: syscall.EPOLLIN,
	})
	if err != nil {
		return fmt.Errorf("configuring epoll: %w", err)
	}

	go func() {
		events := make([]syscall.EpollEvent, 1)

		for {
			ready, err := syscall.EpollWait(epfd, events, 1000)
			if err != nil {
				fmt.Printf("epoll_wait failed: %v", err) // TODO: use logger
				time.Sleep(time.Second * 1000)
			}

			if ready > 0 {
				break
			}
		}

		d.mu.Lock()
		d.qemuPidFd = -1
		if d.mon != nil {
			err := d.mon.Disconnect()
			if err != nil {
				fmt.Printf("warning: closing qmp socket failed: %v\n", err)
			}
			d.mon = nil
		}
		d.signalExit()
		defer d.mu.Unlock()
	}()

	return nil
}

func (d *driver) signalExit() {

}

// Create copies over the rootdisk from the image specified in the configuration, as well as the firmware and nvram files
// This operation effectively destroys any existing root disk and resets the nvram.
func (d *driver) Create() error {
	if d.qemuPidFd != -1 {
		return fmt.Errorf("VM is running")
	}

	rootDiskPath := d.storagePath(RootDiskFileName)
	err := util.RemoveIfExists(rootDiskPath)
	if err != nil {
		return fmt.Errorf("deleting existing root disk: %w", err)
	}

	cmd := exec.Command("qemu-img", "convert", "-f", "qcow2", "-O", "raw", d.config.ImageSourcePath, rootDiskPath)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("converting image (%s): %w", cmd, err)
	}

	var diskKiB uint64

	if d.config.DiskSize%4096 != 0 {
		diskKiB = ((d.config.DiskSize / 4096) + 1) * 4
	} else {
		diskKiB = d.config.DiskSize / 1024
	}

	cmd = exec.Command("qemu-img", "resize", "-f", "raw", rootDiskPath, fmt.Sprintf("%dK", diskKiB))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("resizing image (%s): %w", cmd, err)
	}

	firmwarePath := d.storagePath(FirmwareFileName)
	err = util.RemoveIfExists(firmwarePath)
	if err != nil {
		return fmt.Errorf("deleting existing firmware file: %w", err)
	}

	err = util.CopyFile(d.config.FirmwareSourcePath, firmwarePath)
	if err != nil {
		return fmt.Errorf("copying firmware file: %w", err)
	}

	nvramPath := d.storagePath(NvramFileName)
	err = util.RemoveIfExists(nvramPath)

	if err != nil {
		return fmt.Errorf("deleting existing nvram file: %w", err)
	}

	err = util.CopyFile(d.config.NvramSourcePath, nvramPath)
	if err != nil {
		return fmt.Errorf("copying firmware file: %w", err)
	}

	fmt.Printf("cloud init: %v", d.config.CloudInit)
	if (d.config.CloudInit != CloudInitData{}) {
		fmt.Printf("creating cloud init volume")
		err := d.config.CloudInit.CreateIso(d.storagePath(CloudInitIsoFile))
		if err != nil {
			return fmt.Errorf("creating cloud init iso: %w", err)
		}
	}

	err = os.WriteFile(d.storagePath(CreatedFlagFileName), make([]byte, 0), 0o644)
	if err != nil {
		return fmt.Errorf("creating flag file: %w", err)
	}

	return nil
}

func (d *driver) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.qemuPidFd != -1 {
		return fmt.Errorf("VM is already running")
	}

	configFilePath := d.storagePath(ConfigFileName)
	err := util.RemoveIfExists(configFilePath)
	if err != nil {
		return fmt.Errorf("removing old config file: %w", err)
	}

	builder := cmdBuilder.New(d.qemuPath)
	defer builder.CloseFds()

	builder.AddArgs(
		"-S",
		"-uuid", d.config.Id.String(),
		"-nographic",
		//"-display", "gtk",
		"-nodefaults",
		"-no-user-config",
		"-serial", "chardev:console",
		"-readconfig", configFilePath,
		"-machine", "q35",
		"-pidfile", d.storagePath(QemuPidFileName),
	)

	err = util.RemoveIfExists(d.runtimePath(QemuQmpSocketFileName))
	if err != nil {
		return fmt.Errorf("removing old qmp socket: %w", err)
	}

	err = util.RemoveIfExists(d.runtimePath(ConsoleSocketFileName))
	if err != nil {
		return fmt.Errorf("removing old console socket: %w", err)
	}

	desc := machine.Description{}

	desc.Cpu(1, int(d.config.CpuCount), int(d.config.CpuCount))

	memoryMiB := int(d.config.MemorySize / 1024 / 1024)

	if d.config.MemorySize%(1024*1024) != 0 {
		memoryMiB += 1
	}

	desc.Memory(memoryMiB, 64*1024) // TODO: configurable max size

	desc.Monitor(d.runtimePath(QemuQmpSocketFileName))

	desc.Pcie().AddDevice(pcie.NewBalloon("balloon"))
	desc.Pcie().AddDevice(pcie.NewKeyboard("keyboard"))
	desc.Pcie().AddDevice(pcie.NewTablet("tablet"))
	desc.Pcie().AddDevice(pcie.NewVga("vga", pcie.StdVga))

	desc.AddChardev(chardev.NewRingbuf("console-ringbuf", 4096))

	consoleSocketListener, err := net.Listen("unix", d.runtimePath(ConsoleSocketFileName))
	if err != nil {
		return fmt.Errorf("creating listener for console socket: %w", err)
	}

	err = os.Chmod(d.runtimePath(ConsoleSocketFileName), 0777)
	if err != nil {
		_ = consoleSocketListener.Close()
		return fmt.Errorf("changing permissions on console socket path: %w", err)
	}

	consoleSocketFile, err := consoleSocketListener.(*net.UnixListener).File()
	if err != nil {
		return fmt.Errorf("getting fd from console socket listener: %w", err)
	}

	desc.AddChardev(chardev.NewSocket("console-socket", chardev.SocketOpts{
		Unix: chardev.SocketOptsUnix{
			Fd: builder.AddFd(consoleSocketFile),
		},
		Server: true,
		Wait:   false,
	}))

	desc.AddChardev(chardev.NewHub("console", "console-ringbuf", "console-socket"))

	desc.Scsi().AddDisk(storage.NewImageDrive("rootdisk", d.storagePath(RootDiskFileName)))

	if (d.config.CloudInit != CloudInitData{}) {
		desc.Scsi().AddDisk(storage.NewCdromDrive("cloud-init-cidata", d.storagePath(CloudInitIsoFile)))
	}

	for _, volume := range d.config.Volumes {
		switch v := volume.(type) {
		case CephVolume:
			_ = v
			desc.Scsi().AddDisk(storage.NewRbdDrive()) // TODO: volume config
		}
	}

	for _, networkInterface := range d.config.NetworkInterfaces {
		switch n := networkInterface.(type) {
		case tapNetworkInterface:
			desc.Pcie().AddDevice(pcie.NewTapNetworkDevice("tap-"+n.name, n.name, n.macAddress))
		case physicalNetworkInterface:
			panic("not supported")
			//desc.Pcie().AddDevice(pcie.NewPhysicalNetworkDevice("phys-"+n.Id.String(), n.Name))
		}
	}

	if d.config.VsockCid != 0 {
		vsockFile, err := openVsock(d.config.VsockCid)
		if err != nil {
			return fmt.Errorf("allocating vsock cid %d: %w", d.config.VsockCid, err)
		}

		vsockFd := builder.AddFd(vsockFile)
		desc.Pcie().AddDevice(pcie.NewVsock("vsock", d.config.VsockCid, vsockFd))
	}

	config, hotplugDevices := desc.BuildConfig()

	err = os.WriteFile(configFilePath, []byte(config.ToString()), 0o644)
	if err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	stdout, err := os.Create(d.storagePath(QemuStdOutFileName))
	if err != nil {
		return fmt.Errorf("creating stdout log file: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer stdout.Close()

	stderr, err := os.Create(d.storagePath(QemuStdErrFileName))
	if err != nil {
		return fmt.Errorf("creating stderr log file: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer stderr.Close()

	builder.ConnectStderr(stderr)
	builder.ConnectStdout(stdout)

	builder.SetSession(true)

	var pidfd int
	builder.SetPidFdReceiver(&pidfd)

	cmd := builder.Build()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting qemu process (%s): %w", cmd, err)
	}

	d.qemuPidFd = pidfd
	err = d.startWatcher()
	if err != nil {
		return fmt.Errorf("starting qemu process watcher: %w", err)
	}

	_, err = d.connectMonitor()
	if err != nil {
		return err
	}

	return d.postStartConfig(hotplugDevices)
}

func (d *driver) postStartConfig(devices []devices.HotplugDevice) error {
	mon, err := d.connectMonitor()
	if err != nil {
		return err
	}

	var hotplugErrors []error

	for _, device := range devices {
		err := device.Plug(mon)
		if err != nil {
			hotplugErrors = append(hotplugErrors, err)
		}
	}

	if hotplugErrors != nil {
		return fmt.Errorf("hotplugging: %w", errors.Join(hotplugErrors...))
	}

	err = mon.Continue()
	if err != nil {
		return fmt.Errorf("resuming VM execution: %w", err)
	}

	return nil
}

func (d *driver) connectMonitor() (qmp.Monitor, error) {
	if d.mon != nil {
		_, err := d.mon.Status()
		if err == nil {
			return d.mon, nil
		}

		fmt.Println("connection to qmp socket is broken, reconnecting")
		_ = d.mon.Disconnect()
	}

	var mon qmp.Monitor
	var err error

	for {
		mon, err = qmp.Connect(d.runtimePath(QemuQmpSocketFileName))
		if errors.Is(err, unix.ENOENT) || errors.Is(err, unix.ECONNREFUSED) {
			fmt.Println("qmp socket is not ready yet, retrying in a second")
			time.Sleep(time.Second * 1)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("connecting to qmp monitor: %w", err)
		}

		break
	}

	d.mon = mon
	return mon, nil
}

func (d *driver) storagePath(name StorageFilename) string {
	return path.Join(d.config.StorageDirectory, string(name))
}

func (d *driver) runtimePath(name RuntimeFilename) string {
	return path.Join(d.config.RuntimeDirectory, string(name))
}

func (d *driver) Stop() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	mon, err := d.connectMonitor()
	if err != nil {
		return err
	}

	err = mon.Quit()
	if err != nil {
		return fmt.Errorf("stopping VM: %w", err)
	}

	err = d.mon.Disconnect()
	if err != nil {
		fmt.Printf("closing qmp: %w", err)
	}

	d.mon = nil

	return nil
}

func (d *driver) Reboot() error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) Scale(cpuCount uint32, memory uint64, disk uint64) error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) AttachNetworkInterface(net NetworkInterface) error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) DetachNetworkInterface(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) AttachVolume(volume Volume) error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) DetachVolume(id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (d *driver) GetState() Status {
	// TODO: proper state logic for starting, stopping and restarting states
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.qemuPidFd == -1 {
		isCreated, err := util.FileExists(d.storagePath(CreatedFlagFileName))
		if err != nil {
			fmt.Printf("failed to check creation flag existence: %v\n", err)
			return Unknown
		}

		if !isCreated {
			return Uninitialized
		}

		return Stopped
	}

	mon, err := d.connectMonitor()
	if err != nil {
		fmt.Printf("failed to connect qmp monitor: %v\n", err)
		return Unknown
	}

	qemuStatus, err := mon.Status()
	if err != nil {
		fmt.Printf("failed to query qemu status: %v\n", err)
		return Unknown
	}

	return mapQemuStatus(qemuStatus)
}

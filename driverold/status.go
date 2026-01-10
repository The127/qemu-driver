package driver

import "github.com/gwenya/qemu-driver/qmp"

type Status string

const (
	Uninitialized Status = "uninitialized"
	Starting      Status = "starting"
	Stopping      Status = "stopping"
	Restarting    Status = "restarting"
	Running       Status = "running"
	Stopped       Status = "stopped"
	Paused        Status = "paused"
	PausedSystem  Status = "paused-system" // the VM is in paused state but cannot be resumed by the user, not sure yet if it will be needed or be replaced with specific migration etc states
	Error         Status = "error"
	Unknown       Status = "unknown"
)

func mapQemuStatus(state qmp.RunState) Status {
	switch state {
	case qmp.RunStateDebug:
		return Unknown
	case qmp.RunStateFinishMigrate:
		return PausedSystem
	case qmp.RunStateInmigrate:
		return PausedSystem
	case qmp.RunStateInternalError:
		return Error
	case qmp.RunStateIOError:
		return Unknown // we don't configure devices to pause on IO errors so this state should not happen
	case qmp.RunStatePaused:
		return Paused
	case qmp.RunStatePostmigrate:
		return PausedSystem
	case qmp.RunStatePrelaunch:
		return PausedSystem
	case qmp.RunStateRestoreVm:
		return PausedSystem
	case qmp.RunStateRunning:
		return Running
	case qmp.RunStateSaveVm:
		return PausedSystem
	case qmp.RunStateShutdown:
		return Unknown // we don't start VMs with -no-shutdown so this state should not happen
	case qmp.RunStateSuspended:
		return Unknown // we don't enable S3 so this state should not happen
	case qmp.RunStateWatchdog:
		return Unknown // we don't use watchdog so this state should not happen
	case qmp.RunStatePanicked:
		return Running // TODO: not sure what should happen here, but from the perspective of a cloud provider a VM in which the OS is panicked is still running
	case qmp.RunStateColo:
		return Unknown // we don't use colo so this state should not happen

	default:
		return Unknown
	}
}

package qmp

type RunState string

const (
	RunStateDebug         RunState = "debug"
	RunStateFinishMigrate RunState = "finish-migrate"
	RunStateInmigrate     RunState = "inmigrate"
	RunStateInternalError RunState = "internal-error"
	RunStateIOError       RunState = "io-error"
	RunStatePaused        RunState = "paused"
	RunStatePostmigrate   RunState = "postmigrate"
	RunStatePrelaunch     RunState = "prelaunch"
	RunStateRestoreVm     RunState = "restore-vm"
	RunStateRunning       RunState = "running"
	RunStateSaveVm        RunState = "save-vm"
	RunStateShutdown      RunState = "shutdown"
	RunStateSuspended     RunState = "suspended"
	RunStateWatchdog      RunState = "watchdog"
	RunStatePanicked      RunState = "guest-panicked"
	RunStateColo          RunState = "colo"
)

package shutdown

type ShutdownManager struct {
	CleanupFuncs []func() error
}

func NewShutDownManager() *ShutdownManager {
	return &ShutdownManager{
		CleanupFuncs: make([]func() error, 0),
	}
}

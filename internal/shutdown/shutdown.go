package shutdown

import "context"

type ShutdownManager struct {
	CleanupFuncs            []func() error
	CleanupFuncsWithContext []func(ctx context.Context) error
}

func NewShutDownManager() *ShutdownManager {
	return &ShutdownManager{
		CleanupFuncs:            make([]func() error, 0),
		CleanupFuncsWithContext: make([]func(ctx context.Context) error, 0),
	}
}

package os

import (
	"container/list"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	mutex         = sync.RWMutex{}
	hooks         = list.New()
	signalChannel = make(chan os.Signal, 1)
)

// HookFn a signal handler
type HookFn = func(os.Signal)

func init() {
	startListening()
}

// RegisterShutdownHook registers a hook for system call signals SIGTERM and SIGKILL
func RegisterShutdownHook(hook HookFn) {
	mutex.Lock()
	defer mutex.Unlock()

	hooks.PushBack(hook)
}

// NewSignalHook returns a hook that only invokes the specified function if the received signal matches the specified one.
// Supported signals are SIGTERM and SIGKILL
func NewSignalHook(filter os.Signal, f func()) HookFn {
	return func(sig os.Signal) {
		if sig == filter {
			f()
		}
	}
}

func startListening() {
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-signalChannel

		callHooks(sig)
	}()
}

func clearHooks() {
	mutex.Lock()
	defer mutex.Unlock()

	hooks.Init()
}

func callHooks(sig os.Signal) {
	mutex.RLock()
	defer mutex.RUnlock()

	for e := hooks.Front(); e != nil; e = e.Next() {
		e.Value.(HookFn)(sig)
	}
}

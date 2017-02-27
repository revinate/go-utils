package helper

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ContextFunc func(context.Context)
type Func func()

func InitApp(inits ...Func) {
	for _, init := range inits {
		callback := init
		callback()
	}
}

func RunApp(rootCtx context.Context, runners ...ContextFunc) {
	var exit = make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)

	ctx, cancel := context.WithCancel(rootCtx)
	wg := sync.WaitGroup{}
	for _, runner := range runners {
		callback := runner
		wg.Add(1)
		go func() {
			callback(ctx)
			wg.Done()
		}()

	}

	<-exit    // Wait for Kill Signal
	cancel()  // Tell go-routines that we are done
	wg.Wait() // Wait for go-routines gracefully quit
}

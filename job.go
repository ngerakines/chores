package chores

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
)

// Job is neat.
type Job interface {
	Run(context.Context) error
	Shutdown()
	String() string
}

// SignalError is neat.
type SignalError struct {
	Signal os.Signal
}

type signalJob struct {
	ctx     context.Context
	cancel  context.CancelFunc
	signals []os.Signal
}

// Run is neat.
func Run(ctx context.Context, logger *log.Logger, job Job) (func() error, func(error)) {
	return func() error {
			logger.Println("starting job", job)
			return ignoreCanceled(job.Run(ctx))
		}, func(error) {
			logger.Println("stopping job", job)
			job.Shutdown()
		}
}

// NewSignalJob is neat.
func NewSignalJob(signals ...os.Signal) Job {
	return &signalJob{
		signals: signals,
	}
}

func ignoreCanceled(err error) error {
	if err == nil || err == context.Canceled {
		return nil
	}
	return err
}

func (j *signalJob) String() string {
	return fmt.Sprintf("signal watcher (%s)", j.signals)
}

func (j *signalJob) Run(ctx context.Context) error {
	j.ctx, j.cancel = context.WithCancel(ctx)
	defer j.cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, j.signals...)
	defer signal.Stop(c)

	select {
	case sig := <-c:
		return SignalError{Signal: sig}
	case <-j.ctx.Done():
		return j.ctx.Err()
	}
}

func (j *signalJob) Shutdown() {
	j.cancel()
	<-j.ctx.Done()
}

func (e SignalError) Error() string {
	return e.Signal.String()
}

// Is is neat.
func (e SignalError) Is(target error) bool {
	t, ok := target.(SignalError)
	if !ok {
		return false
	}
	return e.Signal == t.Signal
}

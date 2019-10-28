package runner

import (
	"os/signal"
	"time"
	"os"
	"syscall"
	"errors"
)

var (
	ErrInterrupt = errors.New("program is interrupted")
	ErrTimeout   = errors.New("program running timeout")
)

type TaskFunc func(int)

type Runner struct {
	timeout   <-chan time.Time // 超时
	tasks     []TaskFunc
	complete  chan error     // 完成通知
	interrupt chan os.Signal // 中断执行
}

func New(d time.Duration) *Runner {
	return &Runner{
		timeout:   time.After(d),
		complete:  make(chan error),
		interrupt: make(chan os.Signal, 1),
	}
}

func (r *Runner) AddTasks(task ...TaskFunc) {
	r.tasks = append(r.tasks, task...)
}

func (r *Runner) run() error {
	for i, f := range r.tasks {
		if r.isInerrupt() {
			return ErrInterrupt
		}
		if f == nil {
			continue
		}
		f(i + 1)
	}
	return nil
}

func (r *Runner) isInerrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

func (r *Runner) Start() error {
	signal.Notify(r.interrupt, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		r.complete <- r.run()
		// r.run()
	}()
	select {
	case <-r.timeout:
		return ErrTimeout
	case err := <-r.complete:
		return err
	}
}

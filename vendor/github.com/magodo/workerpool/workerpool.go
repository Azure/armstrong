package workerpool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type task func() (interface{}, error)

type ResultHandler func(interface{}) error

type Result struct {
	Value interface{}
	Error error
}

var ErrSkipTask = errors.New("skip the task")

type WorkPool interface {
	Run(ResultHandler)
	Done() error
	AddTask(task) bool
}

type workPool struct {
	size      int
	taskQueue chan task

	stopAddingTaskDone chan interface{}
	addingTaskStopped  bool

	resultCh chan Result
	done     chan interface{}
	stopOnce sync.Once
	errorCh  chan error
	stopCh   chan interface{}
}

var _ WorkPool = &workPool{}

// NewWorkPool creates a new worker pool which has "size" numbers of workers.
func NewWorkPool(size int) *workPool {
	return &workPool{
		size:      size,
		taskQueue: make(chan task, 1),
		resultCh:  make(chan Result, 1),
		done:      make(chan interface{}),
		errorCh:   make(chan error, 1),
		stopCh:    make(chan interface{}),
	}
}

// Run launches a number of workers (determined by the work pool size) in their own goroutines to run the task.
// Besides, it will launch a separate result handler (can be nil) which consumes the task results.
//
// If any task hits an error, or the result handler hits an error. The workers will stop handling tasks.
// Especially, if there are undergoing tasks running when the error occurs, those tasks will be handled, and the error
// will be appended together, if any.
func (wp *workPool) Run(h ResultHandler) {
	closeChs := make([]chan interface{}, 0, wp.size)
	for i := 0; i < wp.size; i++ {
		closeChs = append(closeChs, make(chan interface{}))
		go func(ch chan interface{}) {
			for t := range wp.taskQueue {
				select {
				case <-wp.stopCh:
					// Throw away the task here. The reason not to use "brea" here is because the AddTask()
					// might sending a new task to the queue, while all the workers are stopped (if using break).
					// That results into the AddTask() hang.
					continue
				default:
					v, err := t()
					if err != nil {
						err = fmt.Errorf("task error: %w", err)
					}
					wp.resultCh <- Result{v, err}
				}
			}
			close(ch)
		}(closeChs[i])
	}

	go func() {
		for _, ch := range closeChs {
			<-ch
		}
		wp.done <- struct{}{}
	}()

	var once sync.Once
	go func() {
		var result error
		for res := range wp.resultCh {
			err := res.Error
			if err == nil && h != nil {
				err = h(res.Value)
				if err != nil {
					err = fmt.Errorf("task handler error: %w", err)
				}
			}
			if err != nil {
				if errors.Is(err, ErrSkipTask) {
					continue
				}
				once.Do(func() {
					close(wp.stopCh)
				})
				result = multierror.Append(result, err)
			}
		}

		wp.errorCh <- result
	}()
}

// AddTask adds new task to the worker pool.
// Users shouldn't call it after calling Done().
// It returns false if the worker pool is stopped due to any error occured in the task or the result handler.
func (wp *workPool) AddTask(t task) bool {
	select {
	case <-wp.stopCh:
		return false
	case wp.taskQueue <- t:
		return true
	}
}

// Done stopps the worker pool and return the error indicating any error occured in the task or the result handler.
// Users should always call only once for this method.
func (wp *workPool) Done() error {
	close(wp.taskQueue)
	<-wp.done
	close(wp.resultCh)
	return <-wp.errorCh
}

package dealConcurrent

import (
	"time"
	"sync"
)

func doTask() {
	//模拟耗时操作
	time.Sleep(200 * time.Millisecond)
	wg.Done()
}

//模拟http接口 每次请求抽象为job
func handle() {
	job := Job{}
	JobQueue <- job
}

var (
	MaxWorker = 1000
	MaxQueue  = 200000
	wg        sync.WaitGroup
)

type Worker struct {
	quit chan bool
}

func NewWorker() Worker {
	return Worker{
		quit: make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			select {
			case <-JobQueue:
				doTask()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

type Job struct {
}

var JobQueue chan Job = make(chan Job, MaxQueue)

type Dispatcher struct {
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) Run() {
	//开始工作
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker()
		worker.Start()
	}
}

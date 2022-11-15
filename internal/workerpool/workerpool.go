package workerpool

import "sync"

type WorkerPool struct {
	jobQueue         chan func()
	maxWorkers       uint8
	workers          []worker
	availableWorkers chan chan func()
	closed           bool
	once             sync.Once
}

func NewWP(nWorkers uint8) *WorkerPool {
	if nWorkers == 0 {
		nWorkers = 1
	}
	wp := WorkerPool{
		jobQueue:         make(chan func(), 1),
		maxWorkers:       nWorkers,
		workers:          make([]worker, 0, nWorkers),
		availableWorkers: make(chan chan func(), nWorkers),
	}

	return &wp
}

func (wp *WorkerPool) Run() {
	for i := uint8(0); i < wp.maxWorkers; i++ {
		worker := newWorker(wp.availableWorkers)
		wp.workers = append(wp.workers, worker)
		worker.start()
	}

	go wp.dispatch()
}

func (wp *WorkerPool) dispatch() {
	for {
		job, ok := <-wp.jobQueue
		if !ok {
			return
		}
		go func(job func()) {
			worker := <-wp.availableWorkers
			worker <- job
		}(job)
	}
}

func (wp *WorkerPool) Stop() {
	wp.once.Do(func() {
		close(wp.jobQueue)
		wp.closed = true
		for _, w := range wp.workers {
			w.stop()
		}
	})
}

func (wp *WorkerPool) Do(job func()) {
	if !wp.closed {
		wp.jobQueue <- job
	}
}

type worker struct {
	jobCh chan func()
	pool  chan chan func()
	done  chan struct{}
}

func newWorker(workerPool chan chan func()) worker {
	return worker{
		jobCh: make(chan func()),
		pool:  workerPool,
		done:  make(chan struct{}),
	}
}

func (w worker) start() {
	go func() {
		for {
			w.registerAvailableWorker()

			select {
			case job := <-w.jobCh:
				job()
			case <-w.done:
				return
			}
		}
	}()
}

func (w worker) registerAvailableWorker() {
	w.pool <- w.jobCh
}

func (w worker) stop() {
	close(w.done)
}

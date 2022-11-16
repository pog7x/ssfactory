package workerpool

import (
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestWorkerPool_OneWorkerMultipleJobs(t *testing.T) {
	t.Parallel()

	wp := NewWP(1)
	wp.Run()
	defer wp.Stop()

	var wg sync.WaitGroup
	msgNum := 5
	for i := 0; i < msgNum; i++ {
		wg.Add(1)
		wp.Do(func() {
			x := 100 + 200
			_ = x
			wg.Done()
		})
	}

	wg.Wait()
}

func TestWorkerPool_MultipleWorkersMultipleJobs(t *testing.T) {
	t.Parallel()

	wp := NewWP(30)
	wp.Run()
	defer wp.Stop()

	var wg sync.WaitGroup
	msgNum := 8
	for i := 0; i < msgNum; i++ {
		wg.Add(1)
		wp.Do(func() {
			x := 100 + 200
			_ = x
			wg.Done()
		})
	}

	wg.Wait()
}

func TestWorkerPool_ConcurrentIncomingJobs(t *testing.T) {
	t.Parallel()

	wp := NewWP(30)
	wp.Run()
	defer wp.Stop()

	var wg sync.WaitGroup
	jobNum := 20
	wg.Add(jobNum)
	for jn := 0; jn < jobNum; jn++ {
		go func(i int) {
			wp.Do(func() {
				x := 100 + 200
				_ = x
				wg.Done()
			})
		}(jn)
	}

	wg.Wait()
}

func TestWorkerPool_Stop(t *testing.T) {
	t.Parallel()

	wp := NewWP(10)
	wp.Run()
	defer wp.Stop()

	var wg sync.WaitGroup
	jobNum := 20
	wg.Add(jobNum)
	done := make(chan struct{})
	for i := 0; i < jobNum; i++ {
		go func(i int) {
			wp.Do(func() {
				<-done
				wg.Done()
			})
		}(i)
	}

	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	wg.Wait()
}

func TestWorkerPool_OverflowJobs(t *testing.T) {
	t.Parallel()

	wp := NewWP(2)
	wp.Run()
	defer wp.Stop()

	var wg sync.WaitGroup
	msgNum := 128
	for i := 0; i < msgNum; i++ {
		wg.Add(1)
		wp.Do(func() {
			x := 100 + 200
			_ = x
			wg.Done()
		})
	}

	wg.Wait()
}

func BenchmarkWorkerPool_Do2Workers(b *testing.B) {
	benchDoWorkers(2, b)
}

func BenchmarkWorkerPool_Do4Workers(b *testing.B) {
	benchDoWorkers(4, b)
}

func BenchmarkWorkerPool_Do8Workers(b *testing.B) {
	benchDoWorkers(8, b)
}

func BenchmarkWorkerPool_Do16Workers(b *testing.B) {
	benchDoWorkers(16, b)
}

func BenchmarkWorkerPool_Do32Workers(b *testing.B) {
	benchDoWorkers(32, b)
}

func BenchmarkWorkerPool_Do64Workers(b *testing.B) {
	benchDoWorkers(64, b)
}

func BenchmarkWorkerPool_Do128Workers(b *testing.B) {
	benchDoWorkers(128, b)
}

func BenchmarkWorkerPool_Do256Workers(b *testing.B) {
	benchDoWorkers(256, b)
}

func BenchmarkWorkerPool_Do512Workers(b *testing.B) {
	benchDoWorkers(512, b)
}

func BenchmarkWorkerPool_Do1024Workers(b *testing.B) {
	benchDoWorkers(1024, b)
}

func BenchmarkWorkerPool_Do2048Workers(b *testing.B) {
	benchDoWorkers(2048, b)
}

func benchDoWorkers(n int, b *testing.B) {
	nWorkers := n
	var wg sync.WaitGroup
	wg.Add(b.N * nWorkers)
	b.ResetTimer()
	b.ReportAllocs()
	wp := NewWP(nWorkers)
	wp.Run()
	defer wp.Stop()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			wp.Do(func() {
				wg.Done()
			})
		}
	}
	wg.Wait()
}

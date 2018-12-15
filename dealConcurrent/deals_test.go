package dealConcurrent

import (
	"testing"
	"runtime"
)

var cpuNum = runtime.NumCPU()

func Benchmark_handle(b *testing.B) {
	b.StopTimer()
	b.StartTimer()

	b.SetParallelism(cpuNum)
	runtime.GOMAXPROCS(cpuNum)

	d := NewDispatcher()
	d.Run()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		handle()
	}
	wg.Wait()
}

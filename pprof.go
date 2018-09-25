package nigiri

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

func StartProfile(prefix string) {
	f, err := os.Create(prefix + "cpu.prof")
	if err != nil {
		log.Panicln("can't create cpu profile", prefix, err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Panicln("can't start CPU profile ", err)
	}

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	tf, err := os.Create(prefix + "trace.out")
	if err != nil {
		log.Panicln("can't create trace profile", prefix, err)
	}
	trace.Start(tf)
}

func heap(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		return
	}
	defer f.Close()
	runtime.GC()
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		return
	}
}

func StopProfile(prefix string) {
	mutex, err := os.Create(prefix + "mutex.prof")
	if err != nil {
		return
	}
	defer mutex.Close()

	block, err := os.Create(prefix + "block.prof")
	if err != nil {
		return
	}
	defer block.Close()

	pprof.Lookup("mutex").WriteTo(mutex, 1)
	pprof.Lookup("block").WriteTo(block, 1)

	pprof.StopCPUProfile()

	heap(prefix + "mem.prof")
	trace.Stop()
}

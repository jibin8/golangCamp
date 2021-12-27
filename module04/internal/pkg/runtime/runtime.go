package runtime

import (
	"os"
	"runtime"
)

func MustInit() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func GetPid() int {
	return os.Getpid()
}

package main

import (
	"log"
	"os"
	"runtime/pprof"

	"github.com/TypicalAM/goread/cmd/goread"
)

var (
	version = "dev"
)

func main() {
	f, err := os.Create("test.prof")
	if err != nil {
		log.Fatal(err)
	}

	if err = pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()
	goread.SetVersion(version)
	goread.Execute()
}

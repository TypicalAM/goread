package main

import (
	"github.com/TypicalAM/goread/cmd/goread"
)

var (
	version = "dev"
)

func main() {
	goread.SetVersion(version)
	goread.Execute()
}

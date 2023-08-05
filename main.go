package main

import (
	"github.com/TypicalAM/goread/cmd/goread"
)

var version = "v1.5.1"

func main() {
	goread.SetVersion(version)
	goread.Execute()
}

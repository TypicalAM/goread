package main

import (
	"embed"

	"github.com/TypicalAM/goread/cmd"
)

//go:embed internal/test/example
var exampleFiles embed.FS

func main() {
	cmd.SetVersion("v1.7.1")
	cmd.SetExampleFiles(exampleFiles)
	cmd.Execute()
}

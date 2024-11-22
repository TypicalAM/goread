package main

import (
	"embed"

	"github.com/TypicalAM/goread/cmd"
)

//go:embed internal/test/example
var exampleFiles embed.FS

func main() {
	cmd.SetVersion("v1.6.6")
	cmd.SetExampleFiles(exampleFiles)
	cmd.Execute()
}

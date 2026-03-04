package main

import (
	gotenbergcmd "github.com/gotenberg/gotenberg/v8/cmd"
	_ "github.com/gotenberg/gotenberg/v8/pkg/standard"

	_ "github.com/PortailPro/gotenberg-gs/pkg/modules/ghostscript"
)

func main() {
	gotenbergcmd.Run()
}

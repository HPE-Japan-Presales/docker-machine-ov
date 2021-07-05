package main

import (
	"os"
	"path"

	"github.com/docker/machine/libmachine/drivers/plugin"
	//	"github.com/fideltak/docker-machine-driver-ov/driver"
	"github.com/HPE-Japan-Presales/docker-machine-driver-ov/driver"
	"github.com/urfave/cli"
)

var version string

var (
	helpTpl = `Driver for Docker-machine and Rancher.
<Version>
{{.Version}}
<Authors>
{{if .Author}}{{.Author}}{{end}}
`
	authorsTpl = `
	Tak (taku.kimura@hpe.com)
	Suguru (suguru.inoue@hpe.com)
	Kazuki (kazuki.otomo@hpe.com)
	`
)

func main() {
	cli.AppHelpTemplate = helpTpl
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Author = authorsTpl
	app.Version = version
	app.Action = func(c *cli.Context) {
		plugin.RegisterDriver(driver.NewDriver("", ""))
	}
	app.Run(os.Args)
}

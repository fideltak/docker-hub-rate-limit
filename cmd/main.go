package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/fideltak/docker-hub-rate-limit/internal/apps"
	"github.com/fideltak/docker-hub-rate-limit/internal/global"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	appName  = "docker-prs"
	appUsage = "docker-hub Pull Rate Status"
)

var (
	version     = "Development"
	appUsageTxt = fmt.Sprintf("%s [global-options] [command] [options]", appName)
	globalFlags = global.Flags
	cmds        = apps.Cmds
)

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "version",
	}

	app := &cli.App{
		Name:      appName,
		Usage:     appUsage,
		Flags:     globalFlags,
		UsageText: appUsageTxt,
		Commands:  cmds,
		Before: func(c *cli.Context) error {
			global.EnableDebug(c.Bool("debug"))
			return nil
		},
		Version: version,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package global

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

//Usage description
var (
	debugFlag = &cli.BoolFlag{
		Name:     "debug",
		Usage:    "for debug",
		Required: false,
	}
)

func EnableDebug(flag bool) {
	if flag {
		log.SetLevel(log.DebugLevel)
	}
}

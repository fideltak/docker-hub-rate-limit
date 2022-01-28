package global

import (
	"github.com/urfave/cli/v2"
)

//Register options
var (
	Flags = []cli.Flag{
		debugFlag,
	}
)

package apps

import (
	"github.com/fideltak/docker-hub-rate-limit/internal/apps/show"
	"github.com/urfave/cli/v2"
)

var (
	Cmds = []*cli.Command{
		show.Cmd,
	}
)

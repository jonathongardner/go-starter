package cli

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var helloCommand = &cli.Command{
	Name:      "hello",
	Aliases:   []string{"b"},
	Usage:     "Say hello to someone",
	ArgsUsage: "[who]",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		who := cmd.Args().Get(0)
		if who == "" {
			who = "world"
		}
		log.Infof("Hello %v", who)

		return nil
	},
}

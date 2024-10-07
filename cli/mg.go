package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/jonathongardner/go-starter/routines"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

type greeting struct {
	name  string
	count int
}

func (g *greeting) Run(rc *routines.Controller) error {
	t := (2 * g.count) + 2
	time.Sleep(time.Duration(t) * time.Second)
	log.Infof("Hello %v (%v)", g.name, t)
	return nil
}

type waiting struct{}

func (w *waiting) Run(rc *routines.Controller) error {
	count := 0
	for {
		select {
		case <-rc.IsDone():
			return nil
		default:
			if count > 8 {
				return fmt.Errorf("to many people!")
			}
			time.Sleep(1 * time.Second)
			log.Info("Still waiting...")
			count++
		}
	}
}

var mgCommand = &cli.Command{
	Name:      "many-greetings",
	Aliases:   []string{"m"},
	Usage:     "say many hellos",
	ArgsUsage: "[whos]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "how-many",
			Value:   5,
			Usage:   "How many to say hello",
			Sources: cli.EnvVars("START_HOW_MANY"),
			Action: func(ctx context.Context, cmd *cli.Command, v int64) error {
				if 0 >= v {
					return fmt.Errorf("flag number to say hello %v must be greater than 0", v)
				}
				return nil
			},
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		if cmd.NArg() == 0 {
			return fmt.Errorf("must pass someone to talk to")
		}

		routineController := routines.NewController()

		routineController.GoBackground(&waiting{})

		i := 0
		for i < cmd.NArg() {
			routineController.Go(&greeting{name: cmd.Args().Get(i), count: i})
			i++
		}

		return routineController.IsFinish()
	},
}

package app

import (
	"github.com/urfave/cli/v2"
)

func New() (app *cli.App) {
	app = &cli.App{
		Name:      "gst",
		Usage:     "Automatically upload task",
		UsageText: "gst [global options] path bucket",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "size",
				Aliases: []string{"s"},
				Usage:   "filter file size[B,K,M,G]",
				Value:   "101G",
			},
			&cli.UintFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "scan interval minutes",
				Value:   30,
			},
			&cli.StringFlag{
				Name:    "ext",
				Aliases: []string{"e"},
				Usage:   "files filter suffix",
				Value:   "gz",
			},
		},
	}
	return
}

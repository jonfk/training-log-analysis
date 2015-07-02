package main

import (
	"github.com/codegangsta/cli"
	"os"
	// "training-log/projections"
)

func main() {
	app := cli.NewApp()
	app.Name = "training-log: data-exporter"
	app.Usage = "Exports training-logs to various formats"
	app.Authors = []cli.Author{cli.Author{Name: "Jonathan D Fok", Email: ""}}

	app.Commands = []cli.Command{
		cli.Command{
			Name:        "csv",
			Usage:       "Export to csv",
			Description: "Exports training-logs to csv",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "",
					Usage: "Path to directory containing training-logs",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "target",
					Usage: "Path to directory to write output. If directory does not exist, it is created",
				},
			},
			Action: ExportCSV,
		},
		cli.Command{
			Name:        "json",
			Usage:       "Export to json",
			Description: `Exports training-logs to json format`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Value: "",
					Usage: "Path to directory containing training-logs",
				},
				cli.StringFlag{
					Name:  "output, o",
					Value: "target",
					Usage: "Path to directory to write output. If directory does not exist, it is created",
				},
			},
			Action: ExportJson,
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}

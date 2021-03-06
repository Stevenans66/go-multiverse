package main

import (
	"path/filepath"

	"github.com/multiverse-vcs/go-multiverse/rpc"
	"github.com/urfave/cli/v2"
)

var importCommand = &cli.Command{
	Action:    importAction,
	Name:      "import",
	Usage:     "Import a repo",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Usage:   "Repo type",
			Value:   "git",
		},
		&cli.StringFlag{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "Repo url",
		},
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{"d"},
			Usage:   "Repo directory",
		},
	},
}

func importAction(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	dir, err := filepath.Abs(c.String("dir"))
	if err != nil {
		return err
	}

	client, err := rpc.NewClient()
	if err != nil {
		return err
	}

	args := rpc.ImportArgs{
		Name: c.Args().Get(0),
		Type: c.String("type"),
		URL:  c.String("url"),
		Dir:  dir,
	}

	var reply rpc.ImportReply
	if err := client.Call("Service.Import", &args, &reply); err != nil {
		return err
	}

	return nil
}

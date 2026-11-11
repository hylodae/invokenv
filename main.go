package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/hylodae/invokenv/internal/command"
	"github.com/hylodae/invokenv/internal/envfile"
)

func cmdRoot(ctx context.Context, cmd *cli.Command) error {
	// Get current path
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	invokenvVars := []string{}

	// Check if secrets file exists
	log.Debugf("Looking for secrets file in `%s`", path)
	if envfile.Exists(path) {
		log.Debugf("Parsing secrets file in `%s`", path)
		envFile, err := envfile.Parse(path)
		if err != nil {
			return err
		}

		invokenvVars, err = envFile.FetchValues()
		if err != nil {
			return err
		}
	} else {
		log.Debugf("No secrets file found in `%s`", path)
	}

	subCommandEnv := os.Environ()
	subCommandEnv = append(subCommandEnv, invokenvVars...)

	subCommand := []string{cmd.Args().Get(0)}
	subCommand = append(subCommand, cmd.Args().Tail()...)

	err = command.RunCommand(subCommand, subCommandEnv)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetOutput(os.Stdout)

	cmd := &cli.Command{
		Name:         "invokenv",
		Usage:        "Build environment variables from external sources",
		StopOnNthArg: new(1),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "error",
				Usage: "Set invokenv's log level",
				Action: func(ctx context.Context, cmd *cli.Command, value string) error {
					logLevel, err := log.ParseLevel(value)
					if err != nil {
						return err
					}
					log.SetLevel(logLevel)
					return nil
				},
			},
		},
		Action: cmdRoot,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

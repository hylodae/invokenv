package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/hylodae/invokenv/internal/command"
	"github.com/hylodae/invokenv/internal/envfile"
)

var (
	Version   = "<filled by ldflags>"
	BuildTime = "<filled by ldflags>"
)

func cmdRoot(ctx context.Context, cmd *cli.Command) (err error) {
	envFilePath := cmd.String("envfile")

	if envFilePath == "" {
		// Get current path
		path, err := os.Getwd()
		if err != nil {
			return err
		}

		envFilePath = fmt.Sprintf("%s/%s", path, envfile.EnvFilename)
	}

	invokenvVars := []string{}

	// Check if secrets file exists
	log.Debugf("Looking for secrets file in `%s`", envFilePath)
	if envfile.Exists(envFilePath) {
		log.Debugf("Parsing secrets file in `%s`", envFilePath)
		envFile, err := envfile.Parse(envFilePath)
		if err != nil {
			return err
		}

		invokenvVars, err = envFile.FetchValues()
		if err != nil {
			return err
		}
	} else {
		log.Debugf("No secrets file found in `%s`", envFilePath)
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

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
	}

	cmd := &cli.Command{
		Name:         "invokenv",
		Version:      fmt.Sprintf("%s (%s)", Version, BuildTime),
		Usage:        "Build environment variables from external sources",
		ArgsUsage:    "SUBCOMMAND",
		StopOnNthArg: new(1),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "error",
				Usage: "set invokenv's log level",
				Action: func(ctx context.Context, cmd *cli.Command, value string) error {
					logLevel, err := log.ParseLevel(value)
					if err != nil {
						return err
					}
					log.SetLevel(logLevel)
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "envfile",
				Aliases: []string{"f"},
				Usage:   "file to use as invokenv config file",
			},
		},
		Action: cmdRoot,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

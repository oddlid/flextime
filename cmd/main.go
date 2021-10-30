package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Variables meant to be set in Makefile and passed to the linker
var (
	VERSION    string
	BUILD_DATE string
	COMMIT_ID  string
)

func getBuildVersion() string {
	var v, c string
	if VERSION != "" {
		v = VERSION
	} else {
		v = "v-UNDEF"
	}
	if COMMIT_ID != "" {
		c = COMMIT_ID
	} else {
		c = "00000000"
	}
	return fmt.Sprintf("%s_%s", v, c)
}

func getCompiledDate() time.Time {
	var compiledTime time.Time
	if BUILD_DATE != "" {
		t, err := time.Parse(time.RFC3339, BUILD_DATE)
		if err != nil {
			log.Debug().
				Str("BUILD_DATE", BUILD_DATE).
				Msg("Invalid date format, using time.Now()")
			compiledTime = time.Now()
		} else {
			compiledTime = t
		}
	} else {
		log.Debug().Msg("BUILD_DATE not set, using time.Now()")
		compiledTime = time.Now()
	}
	return compiledTime
}

func main() {
	app := cli.NewApp()
	app.Name = "flextime"
	app.Usage = "Track flextime +/-"
	app.Copyright = "(C) 2021 Odd Eivind Ebbesen"
	app.Compiled = getCompiledDate()
	app.Version = getBuildVersion()
	app.Authors = []*cli.Author{
		{
			Name:  "Odd E. Ebbesen",
			Email: "oddebb@gmail.com",
		},
	}
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		{
			Name:    "",
			Aliases: []string{""},
			Usage:   "",
			Action:  nil,
			Flags:   []cli.Flag{
				//
			},
		},
		{
			Name:    "",
			Aliases: []string{""},
			Usage:   "",
			Action:  nil,
			Flags:   []cli.Flag{
				//
			},
		},
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Aliases: []string{"l"},
			Value:   "info",
			Usage:   "Log `level` (options: debug, info, warn, error, fatal, panic)",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "Set log-level to debug",
			EnvVars: []string{"DEBUG"},
		},
	}

	app.Before = func(c *cli.Context) error {
		zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999-07:00"
		if c.Bool("debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else if c.IsSet("log-level") {
			level, err := zerolog.ParseLevel(c.String("log-level"))
			if err != nil {
				log.Error().Err(err).Send()
			} else {
				zerolog.SetGlobalLevel(level)
			}
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
		return nil
	}

	app.Run(os.Args)
}

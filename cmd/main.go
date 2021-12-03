package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oddlid/flextime/flex"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// Variables meant to be set in Makefile and passed to the linker
var (
	Version   string
	BuildDate string
	CommitID  string
)

func getBuildVersion() string {
	var v, c string
	if Version != "" {
		v = Version
	} else {
		v = "v-UNDEF"
	}
	if CommitID != "" {
		c = CommitID
	} else {
		c = "00000000"
	}
	return fmt.Sprintf("%s_%s", v, c)
}

func getCompiledDate() time.Time {
	var compiledTime time.Time
	if BuildDate != "" {
		t, err := time.Parse(time.RFC3339, BuildDate)
		if err != nil {
			log.Debug().
				Str("BuildDate", BuildDate).
				Msg("Invalid date format, using time.Now()")
			compiledTime = time.Now()
		} else {
			compiledTime = t
		}
	} else {
		log.Debug().Msg("BuildDate not set, using time.Now()")
		compiledTime = time.Now()
	}
	return compiledTime
}

//func getDBFile(c *cli.Context) error {
//	fileName := c.String("file")
//	if fileName == "" {
//		log.Debug().Msg("Empty filename")
//	}
//	_, err := os.Stat(fileName)
//	if err != nil {
//		log.Error().Err(err).Send()
//	}
//	return nil
//}

func entryPointAdd(c *cli.Context) error {
	log.Debug().Msg("In entryPointAdd")
	if c.Bool("debug") {
		log.Debug().Msg("We have access to global flags even in subcommands")
	}
	//	if err := getDBFile(c); err != nil {
	//		return err
	//	}

	date := c.Timestamp("date")
	log.Debug().Msgf("Date: %s", date)
	amount := c.Duration("amount")
	log.Debug().Msgf("Amount: %.0f minutes", amount.Minutes())
	return nil
}

func main() {
	app := &cli.App{
		Name:                 "flextime",
		Usage:                "Track flextime +/-",
		Copyright:            "(C) 2021 Odd Eivind Ebbesen",
		Compiled:             getCompiledDate(),
		Version:              getBuildVersion(),
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{
				Name:  "Odd E. Ebbesen",
				Email: "oddebb@gmail.com",
			},
		},
		Before: func(c *cli.Context) error {
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
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				EnvVars: []string{"FLEXTIME_FILE"},
				Usage:   "JSON `file` to load/save data from",
			},
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
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"set"},
				Usage:   "Add or set flex time for a given customer",
				Action:  entryPointAdd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "customer",
						Aliases: []string{"c"},
						Usage:   "The customer name for whom to add flex",
					},
					&cli.TimestampFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "Date (YYYY-MM-DD) to add flex for",
						Layout:  flex.ShortDateFormat,
					},
					&cli.DurationFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Usage:   "Amount of flex time",
					},
					&cli.BoolFlag{
						Name:    "overwrite",
						Aliases: []string{"o"},
						Usage:   "Overwrite if matching entry already exists",
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List recorded flex time",
				Action:  nil,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "customer",
						Aliases: []string{"c"},
						Usage:   "The customer for whom to list flex time",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "List flex for all customers",
					},
					&cli.TimestampFlag{
						Name:    "date",
						Aliases: []string{"d"},
						Usage:   "List entries for this specific date",
						Layout:  flex.ShortDateFormat,
					},
					&cli.TimestampFlag{
						Name:    "from",
						Aliases: []string{"f"},
						Usage:   "List entries starting from this date",
						Layout:  flex.ShortDateFormat,
					},
					&cli.TimestampFlag{
						Name:    "to",
						Aliases: []string{"t"},
						Usage:   "List entries up to this date",
						Layout:  flex.ShortDateFormat,
					},
				},
			},
			{
				Name:    "delete",
				Aliases: []string{""},
				Usage:   "Delete flex entries",
				Action:  nil,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "customer",
						Aliases: []string{"c"},
						Usage:   "Customer from whom to delete flex entries",
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "Delete matching entries from all customers",
					},
				},
			},
			//{
			//	Name:    "",
			//	Aliases: []string{""},
			//	Usage:   "",
			//	Action:  nil,
			//	Flags:   []cli.Flag{
			//		//
			//	},
			//},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

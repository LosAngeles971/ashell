package command

import (
	"fmt"
	"io"
	"os"

	"github.com/LosAngeles971/ashell/business"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootDescription = `A wrapper of Linux shell with audit logs.`
)

var loglevel string
var logfile string

var rootCmd = &cobra.Command{
	Use:   "ashell",
	Short: "ashell",
	Long:  rootDescription,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("log level ( %s )\n", loglevel)
		switch loglevel {
		case "trace":
			log.SetLevel(log.TraceLevel)
		case "debug":
			log.SetLevel(log.DebugLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
		if len(logfile) > 0 {
			if f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err == nil {
				log.SetOutput(io.MultiWriter(os.Stdout, f))
			} else {
				panic(err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		w := business.New()
		w.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}

func initConfig() {
}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.AddCommand(briscola.BriscolaCmd)
	// rootCmd.PersistentFlags().StringVar(&loglevel, "log", "info", "log level = info|debug|trace")
	// rootCmd.PersistentFlags().StringVar(&logfile, "logfile", "", "log file")
}

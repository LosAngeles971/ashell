package command

import (
	"fmt"
	"os"

	"github.com/LosAngeles971/s-h-entinel/business"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootDescription = `A wrapper of Linux shell with audit logs.`
	defaultLogfile  = "/var/log/s-h-entinel.log"
)

var loglevel string
var logfilename string
var logFile *os.File
var dumpFilename string
var jsonFilename string
var dumpLogger bool
var jsonLogger bool

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
		if logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644); err == nil {
			log.SetOutput(logFile)
		} else {
			panic(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		w := business.New(
			business.WithDumpLogger(dumpLogger, dumpFilename),
			business.WithJsonLogger(jsonLogger, jsonFilename),
		)
		w.Start()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logFile.Close()
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
	rootCmd.PersistentFlags().StringVar(&loglevel, "log", "info", "log level = info|debug|trace")
	rootCmd.PersistentFlags().StringVar(&logfilename, "logfile", defaultLogfile, "log file")
	rootCmd.PersistentFlags().StringVar(&dumpFilename, "dumpfile", "", "target file of dump logger")
	rootCmd.PersistentFlags().StringVar(&jsonFilename, "jsonfile", "", "target file of json logger")
	rootCmd.PersistentFlags().BoolVar(&dumpLogger, "dumpfile", false, "target file of dump logger")
	rootCmd.PersistentFlags().BoolVar(&jsonLogger, "jsonfile", false, "target file of json logger")
}

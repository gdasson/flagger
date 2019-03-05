package main

import (
	"flag"
	"github.com/knative/pkg/signals"
	"github.com/stefanprodan/flagger/pkg/loadtester"
	"github.com/stefanprodan/flagger/pkg/logging"
	"go.uber.org/zap"
	"log"
	"time"
)

var VERSION = "0.1.0"
var (
	logLevel          string
	port              string
	timeout           time.Duration
	logCmdOutput      bool
	zapReplaceGlobals bool
	zapEncoding       string
)

func init() {
	flag.StringVar(&logLevel, "log-level", "debug", "Log level can be: debug, info, warning, error.")
	flag.StringVar(&port, "port", "9090", "Port to listen on.")
	flag.DurationVar(&timeout, "timeout", time.Hour, "Command exec timeout.")
	flag.BoolVar(&logCmdOutput, "log-cmd-output", true, "Log command output to stderr")
	flag.BoolVar(&zapReplaceGlobals, "zap-replace-globals", false, "Whether to change the logging level of the global zap logger.")
	flag.StringVar(&zapEncoding, "zap-encoding", "json", "Zap logger encoding.")
}

func main() {
	flag.Parse()

	logger, err := logging.NewLoggerWithEncoding(logLevel, zapEncoding)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	if zapReplaceGlobals {
		zap.ReplaceGlobals(logger.Desugar())
	}

	defer logger.Sync()

	stopCh := signals.SetupSignalHandler()

	taskRunner := loadtester.NewTaskRunner(logger, timeout, logCmdOutput)

	go taskRunner.Start(100*time.Millisecond, stopCh)

	logger.Infof("Starting load tester v%s API on port %s", VERSION, port)
	loadtester.ListenAndServe(port, time.Minute, logger, taskRunner, stopCh)
}

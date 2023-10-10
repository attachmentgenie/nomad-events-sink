package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/nomad/api"
	"github.com/mr-karan/nomad-events-sink/pkg/stream"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Create a new context which gets cancelled upon receiving `SIGINT`/`SIGTERM`.
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// Initialise and load the config.
	ko, err := initConfig("config.toml", "NOMAD_EVENTS_SINK_")
	if err != nil {
		fmt.Println("error initialising config", err)
		os.Exit(1)
	}

	var (
		log    = initLogger(ko)
		sink   = initSink(ko, log)
		stream = initStream(ctx, ko, log, func(e api.Event, meta stream.Meta) {
			sink.Add(e)
		})
		opts = initOpts(ko)
	)

	// Initialise a new instance of app.
	app := App{
		log:    log,
		sink:   sink,
		stream: stream,
		opts:   opts,
	}

	// Start an instance of app.
	app.log.WithFields(logrus.Fields{"version": version, "commit": commit, "date": date}).Info("booting nomad events collector")
	app.Start(ctx)
}

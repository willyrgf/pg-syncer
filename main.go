package main

import (
	"context"
	"flag"
	"os"

	"github.com/cryp-com-br/pg-syncer/config"
	"github.com/cryp-com-br/pg-syncer/repository"
	"github.com/cryp-com-br/pg-syncer/syncer"
	log "github.com/sirupsen/logrus"
)

func init() {
	dev := flag.Bool("dev", false, "Set the environment to dev.")
	trace := flag.Bool("trace", false, "Enable trace.")
	flag.Parse()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if *dev {
		log.SetLevel(log.DebugLevel)
	}

	if *trace {
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	}

}

func main() {
	config := config.New()
	config.Load()
	log.Debugf("main(): config=%+v", config)

	repo := repository.New(config.GetRepositoryURL())
	defer repo.Close()
	log.Debugf("main(): repo=%+v", repo)

	syncersConfig := syncer.NewConfig()
	syncersConfig.Load()
	log.Debugf("main(): config=%+v", syncersConfig)

	syncers := syncer.New(syncersConfig)
	defer syncers.Close(syncersConfig)
	log.Debugf("main(): syncers=%+v", syncers)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := syncer.Start(ctx, syncers, syncersConfig); err != nil {
		log.Fatalf("main(): syncer.Start() error=%w", err)
	}
}

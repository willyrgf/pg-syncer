package syncer

import (
	"context"

	"github.com/cryp-com-br/pg-syncer/repository"
	log "github.com/sirupsen/logrus"
)

// SyncMode represents the types of syncs
type SyncMode string

const (
	// FullSync truncate the table and copy all result of query to then
	FullSync SyncMode = "fullsync"
	// OnlyDiff sync only the diff data between the source and destination
	OnlyDiff SyncMode = "onlydiff"
	// PartialSync copy all result from query without truncate the destination
	PartialSync SyncMode = "partialsync"
)

// Syncer is a struct of a sync with repository
type Syncer struct {
	Name string
	Repo repository.Repository
}

// Syncers is a map to Syncer
type Syncers map[string]Syncer

// New create the map Syncers
func New(c *Config) Syncers {
	var syncer Syncer
	syncers := make(Syncers)
	for repoName, repo := range c.Repositories {
		log.Debugf("syncer.New(c): repoName=%s; repo=%+v", repoName, repo)
		syncer = Syncer{
			Name: repoName,
			Repo: repository.New(repo.URL),
		}
		syncers[syncer.Name] = syncer
	}

	return syncers
}

// Close finish the Syncers
func (s Syncers) Close(c *Config) {
	for repoName := range c.Repositories {
		s[repoName].Repo.Close()
	}
}

// Start run all services syncers
func Start(ctx context.Context, s Syncers, c *Config) error {
	log.Debugf("Start(ctx, s, c): %+v, %+v, %+v", ctx, s, c)

	for syncerConfigName, syncerConfig := range c.Syncers {
		// skip disabled syncer access
		if !syncerConfig.Enabled {
			continue
		}

		scheduler, err := syncerConfig.GetScheduler()
		if err != nil {
			log.Errorf("syncer.Start(ctx, s, c): syncerConfig.GetScheduler() err=%w", err)
			return err
		}
		log.Debugf("syncer.Start(ctx, s, c): syncerConfigName=%+v, syncerConfig=%+v, scheduler=%+v, err=%w", syncerConfigName, syncerConfig, scheduler, err)

		service := NewService(s[syncerConfig.SourceRepository], s[syncerConfig.DestinationRepository], syncerConfig)
		_, err = scheduler.Do(service.Run, ctx)
		if err != nil {
			log.Debugf("syncer.Start(): error on create job; err=%w", err)
			return err
		}

		scheduler.StartBlocking()

	}

	return nil
}

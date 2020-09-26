package syncer

import (
	"github.com/cryp-com-br/pg-syncer/repository"
	log "github.com/sirupsen/logrus"
)

// Syncer is a struct of a sync with repository
type Syncer struct {
	Name string
	Repo repository.Repository
}

// Syncers is a map to Syncer
type Syncers map[string]Syncer

// New create the map Syncers
func New(c *Config) *Syncers {
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

	return &syncers
}

// Close finish the Syncers
func (s *Syncers) Close(c *Config) {
	for repoName := range c.Repositories {
		(*s)[repoName].Repo.Close()
	}
}

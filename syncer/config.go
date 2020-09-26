package syncer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cryp-com-br/pg-syncer/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFile        string
	defaultConfigFile = "./syncer.toml"
)

// SyncersAccess informations
type SyncersAccess struct {
	SourceRepository      string `mapstructure:"source_repository"`
	SourceDb              string `mapstructure:"source_db"`
	SourceQuery           string `mapstructure:"source_query"`
	DestinationRepository string `mapstructure:"destination_repository"`
	DestinationDb         string `mapstructure:"destination_db"`
	DestinationTable      string `mapstructure:"destination_table"`
	Periodicity           string `mapstructure:"periodicity"`
	OnlyDiff              bool   `mapstructure:"only_diff"`
	CleanDestinationTable bool   `mapstructure:"clean_destination_table"`
}

// SyncersConf map to SyncersAccess confs with hash like a id
type SyncersConf map[string]SyncersAccess

// RepositoryConfig a config to repository
type RepositoryConfig struct {
	Name string
	URL  string
}

// Repositories map to RepositoryConfig with name like a id
type Repositories map[string]RepositoryConfig

// Config basic config
type Config struct {
	Syncers      SyncersConf
	Repositories Repositories
}

func getEnvConfig(env string) (cfg string) {
	cfg = os.Getenv(env)
	return
}

func getDefaultConfig(file string) (fileConfig string) {
	fileConfig = defaultConfigFile
	if file != "" {
		fileConfig = file
	}

	_, err := os.Stat(fileConfig)
	if err != nil {
		fileConfig = ""
	}

	return
}

func viperCfg() {
	configFile = getDefaultConfig(getEnvConfig("CONF"))
	dir, file := filepath.Split(configFile)
	file = strings.TrimSuffix(file, filepath.Ext(file))
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType("toml")
}

// parse Config configs
func parse(cfg *Config) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("config.Parse(): error=%w", err)
		return
	}

	repos := make(Repositories)
	r := RepositoryConfig{
		Name: "repository_1",
		URL:  viper.GetString("repository_1.url"),
	}
	repos[r.Name] = r

	r = RepositoryConfig{
		Name: "repository_2",
		URL:  viper.GetString("repository_2.url"),
	}
	repos[r.Name] = r

	cfg.Repositories = repos

	var sa []SyncersAccess
	err = viper.UnmarshalKey("syncers.access", &sa)
	if err != nil {
		log.Errorf("config.Parse(): error=%w", err)
		return
	}
	syncers := make(SyncersConf)
	for _, s := range sa {
		syncers[helpers.ToSha256(s)] = s
	}

	cfg.Syncers = syncers

	return
}

// Load configuration
func (c *Config) Load() {
	viperCfg()

	if err := parse(c); err != nil {
		log.Fatalf("config.Load(): Parse(ConfigConf): %w", err)
	}

	log.Debug("config.Load(): configuration loaded")
	log.Debugf("config.ConfigConf=%+v", c)
}

// NewConfig initialize the basic config
func NewConfig() *Config {
	return &Config{}
}

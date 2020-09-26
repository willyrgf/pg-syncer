package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFile        string
	defaultConfigFile = "./config.toml"
)

// RepositoryConfig a config to repository
type RepositoryConfig struct {
	URL string
}

// Config basic config
type Config struct {
	Repository *RepositoryConfig
	Dev        bool
	Trace      bool
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
	viper.SetDefault("config.dev", true)
	viper.SetDefault("config.trace", false)
}

// parse Config configs
func parse(cfg *Config) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("config.Parse(): error=%w", err)
		return
	}

	cfg.Dev = viper.GetBool("config.dev")
	cfg.Trace = viper.GetBool("config.trace")
	cfg.Repository.URL = viper.GetString("repository.url")

	return
}

func logConfig(cfg *Config) {
	log.SetReportCaller(false)
	log.SetLevel(log.InfoLevel)

	if cfg.Dev {
		log.SetLevel(log.DebugLevel)
		log.Debug("init(): dev environment")
	}

	if cfg.Trace {
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
		log.Debug("init(): trace enabled")
	}
}

// New initialize the basic config
func New() *Config {
	return &Config{
		Repository: &RepositoryConfig{},
	}
}

// Load configuration
func (c *Config) Load() {
	viperCfg()

	if err := parse(c); err != nil {
		log.Fatalf("config.Load(): Parse(ConfigConf): %w", err)
	}

	logConfig(c)

	log.Debug("config.Load(): configuration loaded")
	log.Debugf("config.ConfigConf=%+v", c)
}

// GetRepositoryURL get url from config
func (c *Config) GetRepositoryURL() string {
	return c.Repository.URL
}

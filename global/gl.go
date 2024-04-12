package global

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"toes/internal/utils"
)

const (
	// RecommendedName defines the default project name.
	RecommendedName = "toes"
	LogTmFmt        = "2006-01-02 15:04:05"
)

var (
	CfgFile     string
	Cache       *cache.Cache
	RedisClient *redis.Client
	Ctx         = context.Background() // redis 使用的
	Cfg         *Config
)

var (
	mu          sync.Mutex
	defaultName = "apiserver.yaml"

	// RecommendedEnvPrefix defines the ENV prefix used by all service.
	RecommendedEnvPrefix = strings.ToUpper(RecommendedName)
)

// 最先进行初始化的
func InitConfig() {
	initConfig(CfgFile)
}

func initConfig(cfgpath string) {
	if cfgpath != "" {
		viper.SetConfigFile(cfgpath)
	} else {
		// 获取目录
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Add `$HOME/<RecommendedHomeDir>` & `.`
		viper.AddConfigPath(filepath.Join("/etc", RecommendedName))
		viper.AddConfigPath(filepath.Join(home, "."+RecommendedName))
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(".", "conf"))

		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read viper configuration file", "err", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Println("config unmarshal err", "err", err)
	}

	// Print using config file.
	log.Println("Using config file", "file", viper.ConfigFileUsed())

	// Watch config file
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("config file updated")
		if err := viper.Unmarshal(&Cfg); err != nil {
			log.Println("config unmarshal err", "err", err)

			return
		}

		// 暂时日志不更新
	})
}

func InitLocalCache() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
	// Cache.Set("foo", "bar", cache.DefaultExpiration)
}

func InitRedis() {

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Cfg.Redis.Host,
		Password: utils.DecryptInternalValue(Cfg.Seckey.Basekey, Cfg.Redis.Password, "redis"),
		Username: Cfg.Redis.Username,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Println("RedisClient.Ping error")
		cobra.CheckErr(err)
	}
}

package cache

import (
	"fmt"

	"github.com/K1ender/MemeWhisper/internal/config"
	"github.com/bradfitz/gomemcache/memcache"
)

func MustInit(cfg *config.Config) *memcache.Client {
	mc := memcache.New(fmt.Sprintf("%s:%d", cfg.Memcached.Host, cfg.Memcached.Port))

	if err := mc.Ping(); err != nil {
		panic(err)
	}

	return mc
}

package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	L2cache l2cache
)

type l2cache struct {
	Lcache     *cache.Cache
	Ltime      int
	RedisCache *redis.Client
	Rtime      int
	Pre        string
}

func InitL2cache(pre string, ltime, rtime int) {
	if ltime == 0 {
		ltime = 100 // 本地 cache 默认100s
	}
	if rtime == 0 {
		rtime = 200 // redis 缓存 200s
	}
	L2cache.RedisCache = RedisClient
	L2cache.Lcache = Cache
}

func (l *l2cache) Set(key, value string) {
	l.Lcache.Set(l.Pre+key, value, time.Duration(l.Ltime)*time.Second)
	l.RedisCache.Set(Ctx, l.Pre+key, value, time.Duration(l.Rtime)*time.Second)
}

func (l *l2cache) Get(key string) (value string, err error) {
	if x, found := l.Lcache.Get(l.Pre + "foo"); found {
		value = x.(string)
		return
	}
	value, err = l.RedisCache.Get(Ctx, l.Pre+key).Result()
	if err != nil {
		l.Lcache.Set(l.Pre+key, value, time.Duration(l.Ltime)*time.Second)
	}
	return
}

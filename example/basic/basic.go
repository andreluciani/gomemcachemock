package basic

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheInterface interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
}

func SetAndGet(mc MemcacheInterface, item *memcache.Item) (*memcache.Item, error) {
	if err := mc.Set(item); err != nil {
		return nil, err
	}
	return mc.Get(item.Key)
}

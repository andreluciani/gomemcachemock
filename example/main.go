package main

import (
	"fmt"

	"github.com/andreluciani/gomemcachemock/example/basic"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	mc := memcache.New("localhost:11211")
	defer mc.Close()

	item := &memcache.Item{
		Key:   "some-key",
		Value: []byte("my value"),
	}

	result, err := basic.SetAndGet(mc, item)
	if err != nil {
		fmt.Printf("An error occurred: %s", err)
	}
	fmt.Printf("Item with key: %s and value: %s retrieved.\n", result.Key, string(result.Value))
}

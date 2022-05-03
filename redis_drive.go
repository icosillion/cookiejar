package cookiejar

import (
	"context"

	redis "github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type RedisDrive struct {
	client     *redis.Client
	entries    map[string]map[string]entry
	namespaces string
}

func (r *RedisDrive) Set(key string, val map[string]entry) {
	r.entries[key] = val
	r.saveEntries(key)
}

func (r *RedisDrive) Get(key string) map[string]entry {
	return r.entries[key]
}

func (r *RedisDrive) Delete(key string) {
	delete(r.entries, key)
}

func (r *RedisDrive) saveEntries(k string) error {
	v, err := json.MarshalToString(r.entries[k])
	if err != nil {
		return err
	}

	err = r.client.HSet(context.TODO(), k, v).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisDrive) readEntries() {
	keys, err := r.client.HKeys(context.TODO(), r.namespaces).Result()
	if err != nil {
		return
	}

	for _, k := range keys {
		b, err := r.client.HGet(context.TODO(), r.namespaces, k).Bytes()
		if err != nil {
			continue
		}

		e := make(map[string]entry)
		if err := json.Unmarshal(b, &e); err != nil {
			// resolve fail
		}

		r.entries[k] = e
	}
}

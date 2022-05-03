package cookiejar

import (
	"net/http"

	"crypto/sha1"
	"encoding/hex"

	redis "github.com/go-redis/redis/v8"
)

type Storage interface {
	Set(string, map[string]entry)
	Get(string) map[string]entry
	Delete(string)
}

func NewFileJar(filename string, o *Options) (http.CookieJar, error) {
	store := &FileDrive{
		filename: filename,
		entries:  make(map[string]map[string]entry),
	}
	store.readEntries()

	return New(store, o)
}

func NewEntriesJar(o *Options) (http.CookieJar, error) {
	store := &EntriesDrive{
		entries: make(map[string]map[string]entry),
	}
	return New(store, o)
}

func NewRedisJar(namespaces string, client *redis.Client, o *Options) (http.CookieJar, error) {
	if namespaces == "" {
		namespaces = "cookiejar"
	}

	r := sha1.Sum([]byte(namespaces))
	namespaces = hex.EncodeToString(r[:])

	store := &RedisDrive{
		client:     client,
		namespaces: namespaces,
		entries:    make(map[string]map[string]entry),
	}
	store.readEntries()
	return New(store, o)
}

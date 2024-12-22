package pokeapi

import (
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type Client struct {
	Cache  *pokecache.Cache
	Client http.Client
}

func NewClient(interval time.Duration) Client {
	return Client{
		Cache:  pokecache.NewCache(interval),
		Client: http.Client{},
	}
}

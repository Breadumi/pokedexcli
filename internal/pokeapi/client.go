package pokeapi

import (
	"net/http"
	"pokedexcli/internal/pokecache"
	"time"
)

type Client struct {
	cache  *pokecache.Cache
	client http.Client
}

func NewClient(interval time.Duration) Client {
	return Client{
		cache:  pokecache.NewCache(interval),
		client: http.Client{},
	}
}

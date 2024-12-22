package main

import (
	"pokedexcli/internal/pokeapi"
	"sync"
	"time"
)

func main() {
	pokeClient := pokeapi.NewClient(time.Minute * 5)
	con := &config{
		fullClient: pokeClient,
		Pokedex: dex{
			Entries: make(map[string]Pokemon),
			mu:      &sync.Mutex{},
		},
	}
	startRepl(con)

}

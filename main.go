package main

import (
	"pokedexcli/internal/pokeapi"
	"time"
)

func main() {
	pokeClient := pokeapi.NewClient(time.Minute * 5)
	con := &config{
		fullClient: pokeClient,
	}
	startRepl(con)

}

package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pokedexcli/internal/constants"
)

func (c *Client) Locations(url *string) (locs LocationBatch, err error) {

	// Client struct in client.go
	//c.client - http.Client{} struct

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		fmt.Printf("Error generating request: %v", err)
		return LocationBatch{}, err
	}

	if v, ok := c.cache.Get(*url); ok {
		err = json.Unmarshal(v, &locs)
		if err != nil {
			return LocationBatch{}, err
		}
		return locs, nil

	}

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Error getting response: %v", err)
		return LocationBatch{}, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationBatch{}, err
	}

	c.cache.Add(*url, data)

	if err = json.Unmarshal(data, &locs); err != nil {
		return LocationBatch{}, err
	}

	return locs, nil
}

func (c *Client) GetPokemon(tag string) (pokemon []string, err error) {

	url := constants.BaseURL + "location-area/" + tag
	var area Area

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error generating request: %v", err)
		return pokemon, err
	}

	if v, ok := c.cache.Get(url); ok {
		err = json.Unmarshal(v, &area)
		if err != nil {
			return pokemon, err
		}
		for _, encounter := range area.PokemonEncounters {
			pokemon = append(pokemon, encounter.Pokemon.Name)
		}
		return pokemon, nil

	}

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Error getting response: %v", err)
		return pokemon, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return pokemon, err
	}

	c.cache.Add(url, data)

	if err = json.Unmarshal(data, &area); err != nil {
		return pokemon, err
	}

	for _, encounter := range area.PokemonEncounters {
		pokemon = append(pokemon, encounter.Pokemon.Name)
	}

	return pokemon, nil
}

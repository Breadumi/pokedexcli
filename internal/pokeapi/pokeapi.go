package pokeapi

import (
	"fmt"
)

const (
	baseURL = "https://pokeapi.co/api/v2"
)

type locationBatch struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func LocationsNext() error {
	fmt.Println("Hello World!")
	return nil
}

package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationBatch struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func Locations(url *string) (locs LocationBatch, err error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		fmt.Printf("Error getting response: %v", err)
		return LocationBatch{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting repsonse: %v", err)
		return LocationBatch{}, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&locs); err != nil {
		return LocationBatch{}, err
	}

	return locs, nil
}

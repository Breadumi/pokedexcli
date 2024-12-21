package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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

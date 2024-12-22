package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"sync"

	"pokedexcli/internal/constants"
	"pokedexcli/internal/pokeapi"
)

func cleanInput(text string) []string {

	text = strings.ToLower(text)
	splitWords := strings.Fields(text)

	return splitWords
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next       *string
	Prev       *string
	fullClient pokeapi.Client
	args       []string
	Pokedex    dex
}

type dex struct {
	Entries map[string]Pokemon
	mu      *sync.Mutex
}

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  Statistics
	Type   []string
}

type Statistics struct {
	HP    int
	ATK   int
	DEF   int
	SpATK int
	SpDEF int
	SPD   int
}

func startRepl(con *config) {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		con.args = cleanInput(scanner.Text())
		if len(con.args) == 0 {
			continue
		}

		command, ok := getCommands()[con.args[0]]
		if ok {
			err := command.callback(con)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Print("Unknown command - type 'help' for more options\n")
		}

	}
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next page of 20 areas",
			callback:    commandNext,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous page of 20 areas",
			callback:    commandPrev,
		},
		"explore": {
			name:        "explore",
			description: "explore <region> - Lists the Pokemon in <region>",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch <pokemon> - Catch a <pokemon>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect <pokemon> - Inspect <pokemon>'s stats if caught",
			callback:    commandInspect,
		},
	}
	return commands
}

func commandExit(con *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(con *config) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func commandNext(con *config) error {
	url := constants.BaseURL + "location-area?offset=0&limit=20"
	if con.Next == nil {
		con.Next = &url
	}
	locs, err := con.fullClient.Locations(con.Next)
	if err != nil {
		fmt.Println("Error A")
		return err
	}
	con.Next = locs.Next
	con.Prev = locs.Previous

	for _, result := range locs.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandPrev(con *config) error {

	if con.Prev == nil {
		fmt.Println("You're on the first page!")
		return nil
	}

	locs, err := con.fullClient.Locations(con.Prev)
	if err != nil {
		return err
	}
	con.Next = locs.Next
	con.Prev = locs.Previous

	for _, result := range locs.Results {
		fmt.Println(result.Name)
	}

	return nil
}

func commandExplore(con *config) error {
	if len(con.args) != 2 {
		return errors.New("expected 1 argument for command `explore`: type help for more info")
	}
	pokemon, err := con.fullClient.GetPokemonList(con.args[1])
	if err != nil {
		return err
	}
	if len(pokemon) < 1 {
		fmt.Println("No Pokemon found!")
	}
	fmt.Printf("Exploring %s...\n", con.args[1])
	fmt.Println("Found Pokemon:")
	for _, pkmn := range pokemon {
		fmt.Printf(" - %s\n", pkmn)
	}
	return nil
}

func commandCatch(con *config) error {
	if len(con.args) != 2 {
		return errors.New("expected 1 argument for command `catch`: type help for more info")
	}
	name := con.args[1]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	baseExp, err := con.fullClient.GetBaseExp(con.args[1])
	if err != nil {
		return err
	}

	var threshold float64

	if baseExp < 200 {
		threshold = 0.9
	} else if baseExp < 400 {
		threshold = 0.8
	} else {
		threshold = 0.7
	}

	if rand.Float64() < threshold {
		fmt.Printf("%s was caught!\n", name)
		con.AddPokemon(name)
	} else {
		fmt.Printf("%s escaped!\n", name)
	}

	return nil
}

func commandInspect(con *config) error {
	if len(con.args) != 2 {
		return errors.New("expected 1 argument for command `inspect`: type help for more info")
	}
	name := con.args[1]

	pokemon := con.Pokedex.Entries[name]

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf("  -hp: %v\n", pokemon.Stats.HP)
	fmt.Printf("  -attack: %v\n", pokemon.Stats.ATK)
	fmt.Printf("  -defense: %v\n", pokemon.Stats.DEF)
	fmt.Printf("  -special-attack: %v\n", pokemon.Stats.SpATK)
	fmt.Printf("  -special-defense: %v\n", pokemon.Stats.SpDEF)
	fmt.Printf("  -speed: %v\n", pokemon.Stats.SPD)
	fmt.Println("Types:")
	for _, t := range pokemon.Type {
		fmt.Printf("  - %s\n", t)
	}
	return nil
}

func (c *config) AddPokemon(tag string) error {
	url := constants.BaseURL + "pokemon/" + tag
	v, ok := c.fullClient.Cache.Get(url)
	if !ok {
		_, err := c.fullClient.GetBaseExp(tag) // this will put it in the cache
		if err != nil {
			return err
		}
		v, ok = c.fullClient.Cache.Get(url)
		if !ok {
			return errors.New("fatal cache error: unable to cache Pokemon information")
		}
	}

	var pokemon pokeapi.PokemonFull
	err := json.Unmarshal(v, &pokemon)
	if err != nil {
		return err
	}

	c.Pokedex.mu.Lock()
	defer c.Pokedex.mu.Unlock()

	var types []string

	for _, t := range pokemon.Types {
		types = append(types, t.Type.Name)
	}

	tmp := Pokemon{
		Name:   tag,
		Height: pokemon.Height,
		Weight: pokemon.Weight,
		Stats: Statistics{
			HP:    0,
			ATK:   0,
			DEF:   0,
			SpATK: 0,
			SpDEF: 0,
			SPD:   0,
		},
		Type: types,
	}

	for _, stat := range pokemon.Stats {
		switch stat.Stat.Name {
		case "hp":
			tmp.Stats.HP = stat.BaseStat
		case "attack":
			tmp.Stats.ATK = stat.BaseStat
		case "defense":
			tmp.Stats.DEF = stat.BaseStat
		case "special-attack":
			tmp.Stats.SpATK = stat.BaseStat
		case "special-defense":
			tmp.Stats.SpDEF = stat.BaseStat
		case "speed":
			tmp.Stats.SPD = stat.BaseStat
		default:
			return errors.New("invalid stat lookup, possibly unexpected json format")
		}
	}

	c.Pokedex.Entries[tag] = tmp

	return nil

}

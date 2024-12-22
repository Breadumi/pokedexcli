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
	callback    func(*config, ...string) error
}

type config struct {
	Next       *string
	Prev       *string
	fullClient pokeapi.Client
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
		args := []string{}
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		if len(words) > 1 {
			args = words[1:]
		}

		command, ok := getCommands()[words[0]]
		if ok {
			err := command.callback(con, args...)
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
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback:    commandPokedex,
		},
	}
	return commands
}

func commandExit(con *config, args ...string) error {
	if len(args) > 0 {
		return errors.New("too many arguments: expected 0")
	}
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(con *config, args ...string) error {
	if len(args) > 0 {
		return errors.New("too many arguments: expected 0")
	}
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func commandNext(con *config, args ...string) error {
	if len(args) > 0 {
		return errors.New("too many arguments: expected 0")
	}
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

func commandPrev(con *config, args ...string) error {

	if len(args) > 0 {
		return errors.New("too many arguments: expected 0")
	}

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

func commandExplore(con *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("expected 1 argument for command `explore`: type help for more info")
	}
	pokemon, err := con.fullClient.GetPokemonList(args[0])
	if err != nil {
		return err
	}
	if len(pokemon) < 1 {
		fmt.Println("No Pokemon found!")
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, pkmn := range pokemon {
		fmt.Printf(" - %s\n", pkmn)
	}
	return nil
}

func commandCatch(con *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("expected 1 argument for command `catch`: type help for more info")
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	baseExp, err := con.fullClient.GetBaseExp(args[0])
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

func commandInspect(con *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("expected 1 argument for command `inspect`: type help for more info")
	}

	pokemon := con.Pokedex.Entries[args[0]]

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

func commandPokedex(con *config, args ...string) error {
	if len(args) > 0 {
		return errors.New("too many arguments: expected 0")
	}
	pokemap := &con.Pokedex.Entries
	if *pokemap == nil || len(*pokemap) == 0 {
		fmt.Println("It's empty! Go catch some Pokemon!")
	} else {
		for k := range *pokemap {
			fmt.Printf(" - %s\n", k)
		}
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

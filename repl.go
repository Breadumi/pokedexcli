package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"pokedexcli/internal/pokeapi"
)

const (
	baseURL = "https://pokeapi.co/api/v2/"
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
	Next *string
	Prev *string
}

func startRepl() {

	scanner := bufio.NewScanner(os.Stdin)
	con := &config{}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command, ok := getCommands()[words[0]]
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
	url := baseURL + "/" + "location-area"
	if con.Next == nil {
		con.Next = &url
	}
	locs, err := pokeapi.Locations(con.Next)
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

	locs, err := pokeapi.Locations(con.Prev)
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

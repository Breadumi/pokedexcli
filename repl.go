package main

import (
	"fmt"
	"os"
	"strings"

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
	callback    func() error
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
			description: "Displays a batch of 20 areas",
			callback:    commandNext,
		},
	}
	return commands
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, cmd := range getCommands() {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func commandNext() error {
	pokeapi.LocationsNext()
	return nil
}

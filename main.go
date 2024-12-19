package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command, ok := getCommands()[words[0]]
		if ok {
			err := command.callback()
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Print("Unknown command\n")
		}

	}

}

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/heavysider/pokedexcli/internal/pokechance"
	"github.com/heavysider/pokedexcli/internal/pokewrap"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	next     string
	previous string
}

var commands map[string]cliCommand
var pokedex map[string]pokewrap.PokemonResponse

func main() {
	pokedex = map[string]pokewrap.PokemonResponse{}
	fmt.Print("Pokedex > ")
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Prints help",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a specific location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get a description of a previously caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Display all caught pokemons",
			callback:    commandPokedex,
		},
	}
	conf := config{
		next:     pokewrap.URL,
		previous: "",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		args := cleanInput(scanner.Text())
		command := args[0]
		params := args[1:]
		foundCommand, ok := commands[command]
		if ok {
			foundCommand.callback(&conf, params)
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	result := []string{}
	for _, v := range strings.Fields(text) {
		result = append(result, strings.ToLower(v))
	}
	return result
}

func commandExit(conf *config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, params []string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

	`)
	for _, v := range commands {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func commandMap(conf *config, params []string) error {
	if conf.next == "" {
		fmt.Println("You're on the last page!")
	}
	res, err := pokewrap.MapLocations(conf.next)
	if err != nil {
		return err
	}
	conf.next = res.Next
	conf.previous = res.Previous
	for _, result := range res.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapB(conf *config, params []string) error {
	if conf.previous == "" {
		fmt.Println("You're on the first page!")
	}
	fmt.Println(conf.previous)
	res, err := pokewrap.MapLocations(conf.previous)
	if err != nil {
		return err
	}
	conf.next = res.Next
	conf.previous = res.Previous
	for _, result := range res.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(conf *config, params []string) error {
	fmt.Printf("Exploring %v...\n", params[0])
	res, err := pokewrap.ExploreLocation(params[0])
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, result := range res.PokemonEncounters {
		fmt.Println(" -", result.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *config, params []string) error {
	pokemon, err := pokewrap.GetPokemon(params[0])
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", params[0])
	catchChance := pokechance.CalculateCaptureChance(pokemon.BaseExperience)
	roll := rand.Float64()
	if roll <= catchChance {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(conf *config, params []string) error {
	pokemon := params[0]
	found, ok := pokedex[pokemon]
	if !ok {
		fmt.Printf("you did not catch %v yet\n", pokemon)
		return nil
	}
	fmt.Printf(
		`
Name: %v
Height: %v
Weight: %v
Stats:
`, found.Name, found.Height, found.Weight)
	for _, stat := range found.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeP := range found.Types {
		fmt.Printf("  -%v\n", typeP.Type.Name)
	}
	return nil
}

func commandPokedex(conf *config, params []string) error {
	if len(pokedex) == 0 {
		fmt.Println("You don't have any pokemons yet!")
		return nil
	}
	fmt.Println("Your pokedex:")
	for name, _ := range pokedex {
		fmt.Printf("  - %v\n", name)
	}
	return nil
}

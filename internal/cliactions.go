package internal

import (
	"errors"
	"fmt"
	"os"
)

type CommandContext struct {
	Pokedex    map[string]Pokemon
	Pagination *MapPagination
	PokeCache  *Cache
	Query      []string
}

type CommandFunc interface {
	Execute(ctx *CommandContext)
}

type MapPagination struct {
	Next string
	Prev string
}

type CliCommand struct {
	name        string
	description string
	Callback    CommandFunc
}

type HelpCommand struct{}

func (h HelpCommand) Execute(ctx *CommandContext) {
	ShowHelp()
}

type ExitCommand struct{}

func (e ExitCommand) Execute(ctx *CommandContext) {
	os.Exit(0)
}

type ExploreCommand struct{}

func (e ExploreCommand) Execute(ctx *CommandContext) {
	query := ctx.Query

	if query == nil || len(query) < 1 {
		fmt.Println("Please enter a location name")
		return
	}

	fmt.Println("Exploring " + query[0] + "...")

	link := "https://pokeapi.co/api/v2/location-area/" + query[0]

	rawdata, exists := ctx.PokeCache.Get(link)

	if !exists {
		rawdata = getPokemonAPIData(link)
		ctx.PokeCache.Add(link, rawdata)
	}

	handleExploreData(rawdata)
}

type MapCommand struct{}

func (m MapCommand) Execute(ctx *CommandContext) {
	if ctx.Pagination.Next == "" {
		fmt.Println("No next locations")
		return
	}

	rawdata, exists := ctx.PokeCache.Get(ctx.Pagination.Next)

	if !exists {
		rawdata = getPokemonAPIData(ctx.Pagination.Next)
		ctx.PokeCache.Add(ctx.Pagination.Next, rawdata)
	}

	handleLocationData(rawdata, ctx.Pagination, false)
}

type MapBCommand struct{}

func (mb MapBCommand) Execute(ctx *CommandContext) {
	if ctx.Pagination.Next == "" {
		fmt.Println("No next locations")
		return
	}

	rawdata, exists := ctx.PokeCache.Get(ctx.Pagination.Next)

	if !exists {
		rawdata = getPokemonAPIData(ctx.Pagination.Next)
		ctx.PokeCache.Add(ctx.Pagination.Next, rawdata)
	}

	handleLocationData(rawdata, ctx.Pagination, true)
}

type CatchCommand struct{}

func (c CatchCommand) Execute(ctx *CommandContext) {
	query := ctx.Query

	if query == nil || len(query) < 1 {
		fmt.Println("Please enter a Pokemon name")
		return
	}

	link := "https://pokeapi.co/api/v2/pokemon/" + query[0]
	rawdata, exists := ctx.PokeCache.Get(link)

	if !exists {
		rawdata = getPokemonAPIData(link)
		ctx.PokeCache.Add(link, rawdata)
	}

	pokemon, caught := CatchPokemon(rawdata)
	if !caught {
		return
	}
	ctx.Pokedex[pokemon.Name] = pokemon
}

type InspectCommand struct{}

func (c InspectCommand) Execute(ctx *CommandContext) {
	query := ctx.Query[0]

	if query == "" {
		fmt.Println("Please enter a Pokemon name")
		return
	}

	pokemonJson, ok := ctx.Pokedex[query]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return
	}

	fmt.Println("Name: " + pokemonJson.Name)
	fmt.Printf("Height: %d\n", pokemonJson.Height)
	fmt.Printf("Weight: %d\n", pokemonJson.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonJson.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, type_ := range pokemonJson.Types {
		fmt.Printf("  - %s\n", type_.Type.Name)
	}
}

type PokedexCommand struct{}

func (c PokedexCommand) Execute(ctx *CommandContext) {
	pokedex := ctx.Pokedex

	fmt.Println("Your Pokedex:")
	for _, pokemon := range pokedex {
		fmt.Println(" - " + pokemon.Name)
	}
}

func ShowHelp() {
	fmt.Println("Welcome to the Pokedex!\nUsage:")

	for _, command := range cliMap {
		fmt.Printf("  %s: %s\n", command.name, command.description)
	}
}

var cliMap = map[string]CliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		Callback:    HelpCommand{},
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		Callback:    ExitCommand{},
	},
	"map": {
		name:        "map",
		description: "Get the next 20 map locations",
		Callback:    MapCommand{},
	},
	"mapb": {
		name:        "mapb",
		description: "Get the previous 20 map locations",
		Callback:    MapBCommand{},
	},
	"explore": {
		name:        "explore",
		description: "Get all the pokemon in one location",
		Callback:    ExploreCommand{},
	},
	"catch": {
		name:        "catch",
		description: "Catch a pokemon",
		Callback:    CatchCommand{},
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a pokemon",
		Callback:    InspectCommand{},
	},
	"pokedex": {
		name:        "pokedex",
		description: "View your pokedex",
		Callback:    PokedexCommand{},
	},
}

func ExecuteCommand(key string, ctx *CommandContext, params []string) error {
	command, exists := cliMap[key]
	if !exists {
		return errors.New("invalid command key")
	}

	ctx.Query = params

	command.Callback.Execute(ctx)
	ctx.Query = nil

	return nil
}

package main

import (
	"bufio"
	"fmt"
	"github.com/sachamorgese/pokedexcli/internal"
	"os"
	"strings"
)

func main() {
	internal.ShowHelp()
	scanner := bufio.NewScanner(os.Stdin)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	pagination := internal.MapPagination{
		Next: "https://pokeapi.co/api/v2/location-area/",
	}

	context := internal.CommandContext{
		Pagination: &pagination,
		PokeCache:  internal.NewCache(5),
		Pokedex:    make(map[string]internal.Pokemon),
	}

	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		input_substring := strings.Fields(input)

		internal.ExecuteCommand(input_substring[0], &context, input_substring[1:])
	}
}

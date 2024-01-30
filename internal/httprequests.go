package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

func getPokemonAPIData(link string) []byte {
	resp, err := http.Get(link)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()
	return body
}

func handleLocationData(rawdata []byte, pagination *MapPagination, isPrev bool) {
	var locationJson MapLocationResponse
	err := json.Unmarshal(rawdata, &locationJson)

	if err != nil {
		fmt.Println(err)
		return
	}

	if isPrev {
		pagination.Next = pagination.Prev

		if locationJson.Previous != nil {
			pagination.Prev = *locationJson.Previous
		} else {
			pagination.Prev = ""
		}
	} else {
		pagination.Prev = pagination.Next

		if locationJson.Next != nil {
			pagination.Next = *locationJson.Next
		} else {
			pagination.Next = ""
		}
	}

	for _, location := range locationJson.Results {
		fmt.Println(location.Name)
	}
}

func handleExploreData(rawdata []byte) {
	var locationJson LocationResponse
	err := json.Unmarshal(rawdata, &locationJson)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationJson.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
}

func catchChance(baseExp int) float64 {
	if baseExp >= 340 {
		return 0.05 // 5% chance for base experience of 340 or more
	}
	return 1 - float64(baseExp)/400.0
}

func tryCatchPokemon(baseExp int) bool {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return rnd.Float64() < catchChance(baseExp)
}

func CatchPokemon(rawdata []byte) (Pokemon, bool) {
	var pokemonJson Pokemon
	err := json.Unmarshal(rawdata, &pokemonJson)

	if err != nil {
		fmt.Println(err)
		return Pokemon{}, false
	}

	fmt.Println("Throwing a Pokeball at " + pokemonJson.Name + "...")
	results := tryCatchPokemon(pokemonJson.BaseExperience)

	if results {
		fmt.Println(pokemonJson.Name + " was caught!")
		return pokemonJson, true
	} else {
		fmt.Println(pokemonJson.Name + " escaped!")
		return Pokemon{}, false
	}
}

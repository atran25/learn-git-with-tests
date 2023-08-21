package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Character struct {
	CharacterId   string `json:"characterId"`
	CharacterName string `json:"characterName"`
	Fame          int    `json:"fame"`
	ServerId      string `json:"serverId`
	JobName       string `json:"jobName"`
}

type Response struct {
	Data []Character `json:"rows"`
}

const (
	baseURL = "https://api.dfoneople.com/df/servers/cain/characters?characterName={LETTER_COMBINATION}&wordType=full&apikey=f7duulCyzOcAdt3jKEoUyNYrBGlFIAhm&limit=200"
	maxRequests = 450
	rateLimit = time.Second / maxRequests
)

func main() {
	characters := "abcdefghijklmnopqrstuvwxyz0123456789"
	combinationLength := 2
	combinations := generateCombinations(characters, combinationLength)
	fmt.Println(len(combinations))

	var wg sync.WaitGroup
	wg.Add(len(combinations))

	semaphore := make(chan struct{}, maxRequests)

	for _, combination := range combinations {
		go func(combination string) {
			wg.Add(1)
			semaphore <- struct{}{}

			defer func() {
				<- semaphore
				wg.Done()
			}()

			// Get URL for that combination
			urlWithCombo := strings.Replace(baseURL, "{LETTER_COMBINATION}", combination, 1)
			fmt.Println(urlWithCombo)

			res, err := http.Get(urlWithCombo)

			if err != nil {
				fmt.Printf("Combination %s - API request failed %s\n", "ab", err)
				return
			}

			var response Response
			decoder := json.NewDecoder(res.Body)
			if err := decoder.Decode(&response); err != nil {
				fmt.Printf("Combination %s - JSON decoding error: %s\n", "ab", err)
			}

			var charactersData []Character
			charactersData = append(charactersData, response.Data...)
			for _, character := range charactersData {
				fmt.Printf("CharacterId: %s, CharacterName: %s, Fame: %d, ServerId: %s, JobName: %s\n", character.CharacterId, character.CharacterName, character.Fame, character.ServerId, character.JobName)
			}

			wg.Done()
		}(combination)
		time.Sleep(rateLimit)
	}

	wg.Wait()
}

func generateCombinations(characters string, length int) []string {
	var combinations []string

	if length == 0 {
		return []string{""}
	}

	for _, char := range characters {
		subCombinations := generateCombinations(characters, length-1)
		for _, subCombination := range subCombinations {
			combinations = append(combinations, string(char)+subCombination)
		}
	}

	return combinations
}

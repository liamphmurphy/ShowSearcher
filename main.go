package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	APIKey string
}

type Movie struct {
	Results []struct {
		Title string `json:"title"`
	} `json:"results"`
}

func LoadConfig() *Config {
	var configFile Config
	_, confErr := toml.DecodeFile("config.toml", &configFile)
	if confErr != nil {
		panic(confErr.Error())
	}

	return &Config{
		APIKey: configFile.APIKey,
	}

}

func main() {
	config := LoadConfig()
	apiKey := config.APIKey

	// Prepare reader for user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of a movie: ")
	// Actually get user input
	userMovie, _ := reader.ReadString('\n')
	fmt.Println("Searching for: " + userMovie)

	client := &http.Client{}
	requestURL := "https://api.themoviedb.org/3/search/movie"
	req, err := http.NewRequest("GET", requestURL, nil)
	// Use Query to properly encode URL values
	query := req.URL.Query()
	query.Add("api_key", apiKey)
	query.Add("query", userMovie)
	req.URL.RawQuery = query.Encode()

	fmt.Println(req)
	//req, err := http.NewRequest("GET", "http://www.omdbapi.com/?t=Blade+Runner", nil)
	if err != nil {
		panic(err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Println(resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	m := Movie{}
	json.Unmarshal(body, &m)

	fmt.Println(m)

}

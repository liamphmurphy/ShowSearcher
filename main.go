package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ApiKey string
}

type Movie struct {
	Title       string `json:"Title"`
	Year        string `json:"Year"`
	Rating      string `json:"Rated"`
	ReleaseDate string `json:"Released"`
	Runtime     string `json:"Runtime"`
	Genre       string `json:"Genre"`
	Director    string `json:"Director"`
	Writer      string `json:"Writer"`
	Actors      string `json:"Actors"`
}

func LoadConfig() *Config {
	var configFile Config
	_, confErr := toml.DecodeFile("config.toml", &configFile)
	if confErr != nil {
		panic(confErr.Error())
	}

	return &Config{
		ApiKey: configFile.ApiKey,
	}

}

func main() {
	config := LoadConfig()
	apiKey := config.ApiKey

	// Prepare reader for user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of a movie: ")
	// Actually get user input
	userMovie, _ := reader.ReadString('\n')
	fmt.Println("Searching for: " + userMovie)
	// Replace \n with empty value to not mess up http request.
	userMovie = strings.Replace(userMovie, "\n", "", -1)
	// Replace spaces in user input with + for url compatibility
	userMovie = strings.Replace(userMovie, " ", "+", -1)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://www.omdbapi.com/?t="+userMovie+"&apikey="+apiKey, nil)
	//req.Header.Set("apikey", "1c20b2b2")
	//req, err := http.NewRequest("GET", "http://www.omdbapi.com/?t=Blade+Runner", nil)
	if err != nil {
		panic(err.Error())
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	m := Movie{}
	json.Unmarshal(body, &m)

	fmt.Println(m)

	fmt.Println(m.Title)
}

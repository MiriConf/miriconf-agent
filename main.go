package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Hostname struct {
	Hostname string `json:"hostname"`
}

func main() {
	serverToken := os.Getenv("MIRICONF_TOKEN")
	if serverToken == "" {
		log.Fatal("miriconf server token is not specified, set with MIRICONF_TOKEN environment variable")
	}

	data := strings.Split(serverToken, ".")

	tokenDec, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		panic(err)
	}

	var hostname Hostname
	json.Unmarshal(tokenDec, &hostname)

	client := http.Client{}
	response, err := http.NewRequest(http.MethodGet, "http://"+hostname.Hostname+"/api/v1/teams/list", nil)
	if err != nil {
		print(err)
	}

	response.Header = http.Header{
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + serverToken},
	}

	teams, err := client.Do(response)
	if err != nil {
		panic(err)
	}

	defer teams.Body.Close()

	body, err := ioutil.ReadAll(teams.Body)
	if err != nil {
		print(err)
	}

	fmt.Println(string(body))
}

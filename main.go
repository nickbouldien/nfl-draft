package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "encoding/json"

	"github.com/joho/godotenv"
)

// Player : struct for players
type Player struct {
	id       int
	name     string
	school   string
	position string
	year     Year
	drafted  bool
}

// Year : enumerating college class
type Year int

const (
	none      Year = 0
	freshman  Year = 1
	sophomore Year = 2
	junior    Year = 3
	senior    Year = 4
)

var players []Player

var draftedPlayers []Player

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testEnvVar := os.Getenv("TEST")
	fmt.Print("testEnvVar: ", testEnvVar)

	p := Player{id: 1, name: "nick", school: "Tennessee", position: "wr", year: junior, drafted: true}

	// players.append(p)
	players = append(players, p)

	// fmt.Print("nick player: ", p)
	// fmt.Print(" players: ", players)

	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// resp, err := http.Get("http://example.com/")
	// if err != nil {
	//     log.Fatal("failed to get")
	// }
	// defer resp.Body.Close()
	// fmt.Print(resp)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Print("ok")
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Print("test route")
}

func draftPlayer(w http.ResponseWriter, req *http.Request) {
	fmt.Print("draft route")
	if req.Method == http.MethodPost {
		id := req.FormValue("id")
		fmt.Print("id: ", id)
	}
}

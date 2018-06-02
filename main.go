package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Year : enumerating college class
type Year int

const (
	none      Year = 0
	freshman  Year = 1
	sophomore Year = 2
	junior    Year = 3
	senior    Year = 4
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

var players []Player

var draftedPlayers []Player

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testEnvVar := os.Getenv("TEST")
	fmt.Println("testEnvVar: ", testEnvVar)

	allPlayers, err := createPlayerDB("players.txt")
	println("allPlayers: ", allPlayers)

	p := Player{id: 1, name: "nick", school: "Tennessee", position: "wr", year: junior, drafted: true}

	players = append(players, p)
	// fmt.Print(" players: ", players)

	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.HandleFunc("/draft", draftPlayer)

	// resp, err := http.Get("http://www.example.com/")
	// if err != nil {
	// 	log.Fatal("failed to get: ")
	// }
	// defer resp.Body.Close()
	// fmt.Print(resp)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func draftPlayer(w http.ResponseWriter, req *http.Request) {
	fmt.Println("draft route")
	if req.Method == http.MethodPost {
		id := req.FormValue("id")
		fmt.Print("id: ", id)
	}
}

func createPlayerDB(fileName string) ([]Player, error) {
	fmt.Println("createPlayerDB called")
	var err error
	return players, err
}

//  props to https://github.com/mlabouardy/movies-restapi for the below functions
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

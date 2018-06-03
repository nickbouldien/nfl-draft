package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Year : enumerating college class
type Year int

const (
	none           Year = 0
	sophomore      Year = 2
	junior         Year = 3
	senior         Year = 4
	redshirtSenior Year = 5
)

// Player : struct for players
type Player struct {
	id       int64
	name     string
	school   string
	position string
	year     Year
	drafted  bool
}

var players []Player
var allPlayers []Player

var draftedPlayers []Player

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testEnvVar := os.Getenv("TEST")
	fmt.Println("testEnvVar: ", testEnvVar)

	csvFile, e := os.Open("files/players.csv")
	if e != nil {
		log.Fatal(e)
		return
	}

	defer csvFile.Close()
	fmt.Println("csvFile: ", csvFile)

	r := csv.NewReader(csvFile)
	// lines, err := r.ReadAll()
	// if err != nil {
	// 	log.Fatalf("error reading all lines: %v", err)
	// }

	for {
		line, error := r.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal("Error parsing csv: ", error)
		}

		Line00, err := strconv.ParseInt(line[0], 10, 64)
		intermediateYear, _ := strconv.ParseInt(line[4], 10, 64)
		Line04 := Year(intermediateYear)
		Line05, err := strconv.ParseBool(line[5])
		if err != nil {
			log.Fatal(err)
		}

		player := Player{
			id:       Line00,
			name:     line[1],
			school:   line[2],
			position: line[3],
			year:     Line04,
			drafted:  Line05,
		}
		// fmt.Println("player: ", "", " ", player)
		allPlayers = append(allPlayers, player)
	}

	println("allPlayers: ", allPlayers)

	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.HandleFunc("/player/", player)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	// io.WriteString(w, "hola, mundo!")
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func player(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		fmt.Println("trying to draft")
		id := req.FormValue("id")
		fmt.Print("id: ", id)
	}
	fmt.Println("GET - player")
	path := req.URL.Path
	parts := strings.Split(path, "/")
	id := parts[2]
	playerID, _ := strconv.ParseInt(id, 10, 64)

	fmt.Println("path: ", path)
	fmt.Println("parts: ", parts)
	fmt.Println("len(parts): ", len(parts))
	fmt.Println("playerID: ", playerID)

	if len(parts) > 3 {
		log.Fatal("You can only get info for one player at a time")
	}

	foundPlayer := allPlayers[playerID]
	fmt.Println("foundPlayer: ", foundPlayer)

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

// resp, err := http.Get("http://www.example.com/")
// if err != nil {
// 	log.Fatal("failed to get: ")
// }
// defer resp.Body.Close()
// fmt.Print(resp)

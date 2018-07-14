package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// // Year : enumerating college class
// type Year int

// TODO: change this. should be a map from year(int) to class(string)
// const (
// 	none           Year = 0
// 	sophomore      Year = 2
// 	junior         Year = 3
// 	senior         Year = 4
// 	redshirtSenior Year = 5
// )

// Player : struct for players
type Player struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" sql:"full_name"`
	School   string `json:"school"`
	Position string `json:"position"`
	// Year     Year   `json:"year"`
	Drafted bool `json:"drafted"`
}

var players []Player
var allPlayers []Player
var draftedPlayerIDs []string

var draftedPlayers []int64

func (p Player) String() string {
	return fmt.Sprintf("Player<ID=%d Name=%q>", p.ID, p.Name)
}

func main() {
	fmt.Println("Starting server...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sqlUser := os.Getenv("SQL_USER")
	sqlDbName := os.Getenv("DB_NAME")
	sqlPW := os.Getenv("SQL_PW")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", sqlUser, sqlPW, sqlDbName)
	fmt.Println("connStr: ", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InitStore(&dbStore{db: db})

	fs := http.StripPrefix("/files", http.FileServer(http.Dir("./files")))
	http.Handle("/files/", fs)

	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.HandleFunc("/players/", playerHandler) // TODO: add param for non-drafted players
	http.HandleFunc("/scouting", scoutingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL
	path := url.Path
	pattern, _ := regexp.Compile(`/players/(\d+)`)
	matches := pattern.FindStringSubmatch(path)
	fmt.Println("matches: ", matches)

	// TODO: implement sorting
	// query := url.Query()
	// sort := query.Get("sort")
	// dir := query.Get("dir")
	// fmt.Printf("sort: %s , dir: %s \n", sort, dir)

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			panic(err)
		}

		f := req.Form
		id := f.Get("id")

		fmt.Println("trying to draft player: ", id)

		playerID, _ := strconv.ParseInt(id, 10, 64)

		foundPlayer := allPlayers[playerID]

		fmt.Println("foundPlayer: ", foundPlayer)

		hasBeenDrafted := foundPlayer.Drafted
		fmt.Println("hasBeenDrafted: ", hasBeenDrafted)
		if hasBeenDrafted {
			msg := "player has already been drafted. pick another player"
			respondWithError(w, http.StatusNotFound, msg)
			return
		}

		fmt.Println("foundPlayer: ", foundPlayer)

		msg := "Congrats, you have successfully drafted player: " + id
		respondWithJSON(w, http.StatusOK, msg)
		return
	}

	if len(matches) == 0 {
		// return all players since they are not looking for specific player
		fmt.Println("returning all players")

		players, err := store.Players()

		if err != nil {
			msg := "Could not retrieve players"
			respondWithError(w, http.StatusNotFound, msg)
			return
		}

		respondWithJSON(w, http.StatusOK, players)
		return
	}

	strID := matches[1]
	id, _ := strconv.Atoi(strID)
	fmt.Println("GET - player: ", id)
	fmt.Println("path: ", path)

	// TODO: error handling for not finding player ??
	foundPlayer, err := store.Player(id)
	if err != nil {
		msg := "Could not retrieve player with id: " + strID
		respondWithError(w, http.StatusNotFound, msg)
		return
	}

	respondWithJSON(w, http.StatusOK, foundPlayer)
	return
}

func scoutingHandler(w http.ResponseWriter, req *http.Request) {
	players, err := store.Scout()

	if err != nil {
		msg := "Could not retrieve undrafted players"
		respondWithError(w, http.StatusNotFound, msg)
		return
	}

	respondWithJSON(w, http.StatusOK, players)
	return
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

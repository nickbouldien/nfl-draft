package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Year : enumerating college class
type Year int

// ID : id (index) for the player
// type ID int64

// TODO: change this. should be a map from year(int) to class(string)
const (
	none           Year = 0
	sophomore      Year = 2
	junior         Year = 3
	senior         Year = 4
	redshirtSenior Year = 5
)

// Player : struct for players
type Player struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" sql:"full_name"`
	School   string `json:"school"`
	Position string `json:"position"`
	Year     Year   `json:"year"`
	Drafted  bool   `json:"drafted"`
}

var players []Player
var allPlayers []Player
var draftedPlayerIDs []string

var draftedPlayers []int64

type dbStore struct {
	db *sql.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	testEnvVar := os.Getenv("TEST")
	fmt.Println("testEnvVar: ", testEnvVar)
	sqlUser := os.Getenv("SQL_USER")
	sqlDbName := os.Getenv("DB_NAME")
	// sqlPW := os.Getenv("SQL_PW")

	// connStr := os.Getenv("DB_CONN_STRING")
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", sqlUser, sqlDbName)
	fmt.Println("connStr: ", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// db, err := connectToDB()
	rows, err := db.Query("SELECT id, full_name FROM players")
	if err != nil {
		fmt.Println("Failed to run query", err)
		return
	}

	for rows.Next() {
		player := &Player{}

		if err := rows.Scan(&player.ID, &player.Name, &player.School, &player.Position, &player.Year, &player.Drafted); err != nil {
			log.Fatal(err)
		}

		// birds = append(birds, bird)
		fmt.Println("player: ", player)
	}

	csvFile, err := os.Open("files/players.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer csvFile.Close()
	fmt.Println("csvFile: ", csvFile)

	r := csv.NewReader(csvFile)

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error parsing csv: ", err)
		}

		Line00, err := strconv.ParseInt(line[0], 10, 64)
		intermediateYear, _ := strconv.ParseInt(line[4], 10, 64)
		Line04 := Year(intermediateYear)
		Line05, err := strconv.ParseBool(line[5])
		if err != nil {
			log.Fatal(err)
		}

		player := Player{
			ID:       Line00,
			Name:     line[1],
			School:   line[2],
			Position: line[3],
			Year:     Line04,
			Drafted:  Line05,
		}
		allPlayers = append(allPlayers, player)
	}

	println("allPlayers: ", allPlayers)

	// create file (if it doesn't exist)
	file, err := os.Create("drafted_players")
	if err != nil {
		log.Fatal("Error creating 'drafted_players file': ", err)
	}
	defer file.Close()

	fs := http.StripPrefix("/files", http.FileServer(http.Dir("./files")))
	http.Handle("/files/", fs)

	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)
	http.HandleFunc("/players/", player) // TODO: add param for non-drafted players
	// http.HandleFunc("/scouting", scouting)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func connectToDB() #sql.DB {
// 	db, err := sql.Open("postgres", { connection string })
// 	if err != nil {
// 		log.Fatal("Could not connect to the database: ", err)
// 	}

// 	return db
// }

func index(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func storeDraftedPlayer(id string) {
	// file, err := os.Create("drafted_players")
	// if err != nil {
	// 	log.Fatal("Error creating 'drafted_players file': ", err)
	// }
	// defer file.Close()
	file, err := os.Open("drafted_players")
	if err != nil {
		log.Fatal("Error loading 'drafted_players' file: ", err)
	}

	file.WriteString(id)
	fmt.Println("drafted player ids 2: ", draftedPlayerIDs)
}

func loadDraftedPlayers() {
	file, err := os.Open("drafted_players")
	if err != nil {
		log.Fatal("Error loading 'drafted_players' file: ", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		draftedPlayerIDs = append(draftedPlayerIDs, scanner.Text())
	}
	// return draftedPlayerIDs
}

func player(w http.ResponseWriter, req *http.Request) {
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
		// id := req.FormValue("id")

		fmt.Println("trying to draft player: ", id)
		// fmt.Println("theID: ", theID)

		playerID, _ := strconv.ParseInt(id, 10, 64)

		if playerID > 32 { // obviously this is only used since I'm not using real IDs (hashes), just indexes
			msg := "PlayerIDs are between 1 and 32 (inclusive)"
			respondWithError(w, http.StatusNotFound, msg)
			return
		}
		foundPlayer := allPlayers[playerID]

		fmt.Println("foundPlayer: ", foundPlayer)

		hasBeenDrafted := foundPlayer.Drafted
		fmt.Println("hasBeenDrafted: ", hasBeenDrafted)
		if hasBeenDrafted {
			msg := "player has already been drafted. pick another player"
			respondWithError(w, http.StatusNotFound, msg)
			return
		}

		foundPlayer.Drafted = true
		// TODO: persist this (write to csv??) so that player can't be drafted again
		storeDraftedPlayer(id)
		fmt.Println("drafted player ids: ", draftedPlayerIDs)

		fmt.Println("foundPlayer: ", foundPlayer)

		draftedPlayers = append(draftedPlayers, playerID)
		fmt.Println("new draftedPlayers: ", draftedPlayers)

		msg := "Congrats, you have successfully drafted player: " + id
		respondWithJSON(w, http.StatusOK, msg)
		return
	}

	if len(matches) == 0 {
		// return all players since they are not looking for specific player
		fmt.Println("returning all players")

		respondWithJSON(w, http.StatusOK, allPlayers)
		return
	}

	id, _ := strconv.Atoi(matches[1])
	fmt.Println("GET - player: ", id)

	fmt.Println("path: ", path)

	if id > 32 { // obviously this is only used since I'm not using real IDs (hashes), just indeces
		msg := "PlayerIDs are between 1 and 32 (inclusive)"
		respondWithError(w, http.StatusNotFound, msg)
		return
	}

	// TODO: error handling for not finding player ??
	foundPlayer := allPlayers[id-1]
	fmt.Println("foundPlayer: ", foundPlayer)

	respondWithJSON(w, http.StatusOK, foundPlayer)
	return // is this needed nb???
}

// func createPlayerDB(fileName string) ([]Player, error) {
// 	fmt.Println("createPlayerDB called")
// 	var err error
// 	return players, err
// }

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

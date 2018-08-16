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
	http.HandleFunc("/players/", playerHandler) // TODO: add param to get non-drafted players?? (get rid of /scouting route)
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

	// TODO: implement sorting
	// query := url.Query()
	// sort := query.Get("sort")
	// dir := query.Get("dir")
	// fmt.Printf("sort: %s , dir: %s \n", sort, dir)

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Fatal(err) // TODO: best way to deal with this???
		}

		f := req.Form
		id := f.Get("id")

		fmt.Println("trying to draft player: ", id)

		playerID, _ := strconv.ParseInt(id, 10, 64)

		pID, err := store.DraftPlayer(int(playerID))

		switch {
		case err == AlreadyDraftedErr:
			msg := (AlreadyDraftedErr).Error()
			fmt.Printf("Already drafted err. %s", msg)
			respondWithError(w, http.StatusNotFound, msg)
			return
		case pID == 0:
			fmt.Println("did not find player")
			msg := "Player with id: " + strconv.FormatInt(playerID, 10) + " does not exist."
			respondWithJSON(w, http.StatusNotFound, msg)
			return
		case err != nil:
			msg := "Could not draft player with id: " + strconv.FormatInt(playerID, 10)
			respondWithError(w, http.StatusNotFound, msg)
			return
		default:
			msg := "Congrats, you have successfully drafted player: " + strconv.FormatInt(playerID, 10)
			respondWithJSON(w, http.StatusOK, msg)
			return
		}
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

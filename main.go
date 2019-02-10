package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nfl_draft/utils"
	"os"
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
	Drafted bool `json:"drafted"`
}

func (p Player) String() string {
	return fmt.Sprintf("Player<ID=%d Name=%q>", p.ID, p.Name)
}

func main() {
	log.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sqlUser := os.Getenv("SQL_USER")
	sqlDbName := os.Getenv("DB_NAME")
	sqlPW := os.Getenv("SQL_PW")

	var connStr string
	if sqlPW != "" {
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", sqlUser, sqlPW, sqlDbName)
	} else {
		connStr = fmt.Sprintf("user=%s dbname=%s sslmode=disable", sqlUser, sqlDbName)
	}
	log.Println("sqlUser: ", sqlUser)
	log.Println("sqlDbName: ", sqlDbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	InitStore(&dbStore{db: db})

	fs := http.StripPrefix("/files", http.FileServer(http.Dir("./files")))
	http.Handle("/files/", fs)

	http.HandleFunc("/players/", playerHandler) // TODO: add param to get non-drafted players?? (get rid of /scouting route)
	http.HandleFunc("/", index)
	http.HandleFunc("/test", test)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	utils.EnableCors(&w)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"index": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	utils.EnableCors(&w)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	utils.EnableCors(&w)
	_, rest := utils.ShiftPath(req.URL.Path) // to take out "players" from the path

	head, _ := utils.ShiftPath(rest)

	var pID string
	if head != "" {
		id, err := strconv.Atoi(head)
		if err == nil {
			pID = strconv.Itoa(id)
			head = "id" // hacky solution ...
		}
		switch head {
		case "reset":
			num, err := store.Reset()
			if err != nil {
				msg := "Could not reset the players to be undrafted"
				utils.RespondWithError(w, http.StatusInternalServerError, msg)
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, num)
			return
		case "id":
			switch req.Method {
			case "GET":
				log.Println("getting player with id: ", pID)
				p, err := store.Player(id)
				if err != nil {
					msg := "Could not retrieve player with id: " + pID
					utils.RespondWithError(w, http.StatusNotFound, msg)
					return
				}
				utils.RespondWithJSON(w, http.StatusOK, p)
				return
			case "POST":
				log.Println("trying to draft player:", pID)
				playerID, err := store.DraftPlayer(id)
				switch {
				case err == AlreadyDraftedErr:
					msg := (AlreadyDraftedErr).Error()
					log.Printf("Already drafted err. %s", msg)
					utils.RespondWithError(w, http.StatusNotFound, msg)
					return
				case playerID == 0:
					log.Println("did not find player")
					msg := "Player with id: " + pID + " does not exist."
					utils.RespondWithJSON(w, http.StatusNotFound, msg)
					return
				case err != nil:
					msg := "Could not draft player with id: " + pID
					utils.RespondWithError(w, http.StatusNotFound, msg)
					return
				default:
					msg := "Congrats, you have successfully drafted player: " + pID
					utils.RespondWithJSON(w, http.StatusOK, msg)
					return
				}
			default:
				msg := "Method not allowed"
				utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
				return
			}
		default:
			utils.RespondWithError(w, http.StatusNotFound, "Not found")
			return
		}
	}

	switch req.Method {
	case "GET":
		// TODO: implement sorting
		// query := url.Query()
		// sort := query.Get("sort")
		// dir := query.Get("dir")
		// fmt.Printf("sort: %s , dir: %s \n", sort, dir)
		log.Println("returning all players")
		players, err := store.Players()
		if err != nil {
			msg := "Could not retrieve players"
			utils.RespondWithError(w, http.StatusNotFound, msg)
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, players)
		return
	case "POST":
		// TODO - make route to create new players
		msg := "Desired method not yet implemented"
		utils.RespondWithError(w, http.StatusNotImplemented, msg)
		return
	default:
		msg := "Method not allowed"
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}
}

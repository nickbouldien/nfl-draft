package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nfl_draft/utils"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Player : struct for players
type Player struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" sql:"full_name"`
	School   string `json:"school"`
	Position string `json:"position"`
	Drafted  bool   `json:"drafted"`
}

// Team : struct for teams
type Team struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Conference string `json:"conference"`
	Division   string `json:"division"`
	DraftOrder int64  `json:"draftOrder" sql:"draft_order"`
}

func (t Team) String() string {
	return fmt.Sprintf("Team<ID=%d Name=%s Conference=%s Division=%s DraftOrder=%d>", t.ID, t.Name, t.Conference, t.Division, t.DraftOrder)
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

	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/test", test).Methods("GET")

	//http.HandleFunc("/players/", playerHandler) // TODO: add param to get non-drafted players?? (get rid of /scouting route)
	r.HandleFunc("/players", playerHandler).Methods("GET", "POST")
	s1 := r.PathPrefix("/players").Subrouter()
	s1.HandleFunc("/{id:[0-9]+}", playerDetail).Methods("GET", "POST")

	r.HandleFunc("/teams", teamHandler).Methods("GET", "POST")
	s2 := r.PathPrefix("/teams").Subrouter()
	s2.HandleFunc("/{id:[0-9]+}", teamDetail).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func index(w http.ResponseWriter, _ *http.Request) {
	utils.EnableCors(&w)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"index": "success"})
}

func test(w http.ResponseWriter, _ *http.Request) {
	utils.EnableCors(&w)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func teamDetail(w http.ResponseWriter, req *http.Request) {
	reqID := mux.Vars(req)["id"]

	var teamID = -1

	id, err := strconv.Atoi(reqID)
	if err == nil {
		teamID = id
	}

	if teamID == -1 {
		msg := "Could not retrieve team with id: " + reqID
		utils.RespondWithError(w, http.StatusNotFound, msg)
		return
	}

	log.Println("returning team with id ", teamID)
	team, err := store.Team(teamID)
	if err != nil {
		msg := "Could not retrieve team with id: " + reqID
		utils.RespondWithError(w, http.StatusNotFound, msg)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, team)
	return
}

func teamHandler(w http.ResponseWriter, req *http.Request) {
	utils.EnableCors(&w)

	switch req.Method {
	case "GET":
		log.Println("returning all teams")
		teams, err := store.Teams()
		if err != nil {
			msg := "Could not retrieve teams"
			utils.RespondWithError(w, http.StatusNotFound, msg)
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, teams)
		return
	case "POST":
		// TODO - make route to create new team
		msg := "Desired method not yet implemented"
		utils.RespondWithError(w, http.StatusNotImplemented, msg)
		return
	default:
		msg := "Method not allowed"
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}
}

func playerDetail(w http.ResponseWriter, req *http.Request) {
	reqID := mux.Vars(req)["id"]

	var pID = -1

	id, err := strconv.Atoi(reqID)
	if err == nil {
		pID = id
	}

	if pID == -1 {
		msg := "Could not retrieve team with id: " + reqID
		utils.RespondWithError(w, http.StatusNotFound, msg)
		return
	}

	switch req.Method {
	case "GET":
		log.Println("getting player with id: ", pID)
		p, err := store.Player(id)
		if err != nil {
			msg := "Could not retrieve player with id: " + reqID
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
			msg := fmt.Sprintf("Player with id: %s does not exist", reqID)
			utils.RespondWithJSON(w, http.StatusNotFound, msg)
			return
		case err != nil:
			msg := "Could not draft player with id: " + reqID
			utils.RespondWithError(w, http.StatusNotFound, msg)
			return
		default:
			msg := "Congrats, you have successfully drafted player: " + reqID
			utils.RespondWithJSON(w, http.StatusOK, msg)
			return
		}
	default:
		msg := "Method not allowed"
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	utils.EnableCors(&w)
	_, rest := utils.ShiftPath(req.URL.Path) // to take out "players" from the path

	head, _ := utils.ShiftPath(rest)

	if head != "" {
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

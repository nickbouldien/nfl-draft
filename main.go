package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

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
	fmt.Println("Starting server...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sqlUser := os.Getenv("SQL_USER")
	sqlDbName := os.Getenv("DB_NAME")
	//sqlPW := os.Getenv("SQL_PW")

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", sqlUser, sqlDbName)
	//if sqlPW != "" {
	//	connStr += fmt.Sprintf("dbpassword=%s", sqlPw)
	//}
	fmt.Println("sqlUser: ", sqlUser)
	fmt.Println("sqlDbName: ", sqlDbName)

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
	//http.HandleFunc("/scouting", scoutingHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	respondWithJSON(w, http.StatusOK, map[string]string{"index": "success"})
}

func test(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	respondWithJSON(w, http.StatusOK, map[string]string{"test": "success"})
}

func playerHandler(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	basePath, rest := ShiftPath(req.URL.Path) // to take out "players" from the path
	log.Println("basePath: ", basePath, rest)

	head, tail := ShiftPath(rest)
	log.Println("head: ", head)
	log.Println("tail: ", tail)

	var pID string
	if head != "" {
		log.Println("head != '': ", head)
		id, err := strconv.Atoi(head)
		if err == nil {
			log.Println("id: ", id)
			pID = strconv.Itoa(id)
			head = "id" // hacky solution ...
		}
		switch head {
		case "reset":
			num, err := store.Reset()
			if err != nil {
				msg := "Could not reset the players to be undrafted"
				respondWithError(w, http.StatusInternalServerError, msg)
				return
			}
			respondWithJSON(w, http.StatusOK, num)
			return
		case "id":
			switch req.Method {
			case "GET":
				fmt.Println("getting player with id: ", pID)
				p, err := store.Player(id)
				if err != nil {
					msg := "Could not retrieve player with id: " + pID
					respondWithError(w, http.StatusNotFound, msg)
					return
				}
				respondWithJSON(w, http.StatusOK, p)
				return
			case "POST":
				fmt.Println("trying to draft player: ", pID)
				//playerID := int64(id)// strconv.ParseInt(strID, 10, 64)
				playerID, err := store.DraftPlayer(id)

				switch {
				case err == AlreadyDraftedErr:
					msg := (AlreadyDraftedErr).Error()
					fmt.Printf("Already drafted err. %s", msg)
					respondWithError(w, http.StatusNotFound, msg)
					return
				case playerID == 0:
					fmt.Println("did not find player")
					msg := "Player with id: " + pID + " does not exist."
					respondWithJSON(w, http.StatusNotFound, msg)
					return
				case err != nil:
					msg := "Could not draft player with id: " + pID
					respondWithError(w, http.StatusNotFound, msg)
					return
				default:
					msg := "Congrats, you have successfully drafted player: " + pID
					respondWithJSON(w, http.StatusOK, msg)
					return
				}
			default:
				msg := "Method not allowed"
				respondWithError(w, http.StatusMethodNotAllowed, msg)
				return
			}
		default:
			respondWithError(w, http.StatusNotFound, "Not found")
			return
		}
	}

	switch req.Method {
	case "GET":
		// return all players since they are not looking for specific player
		//	// TODO: implement sorting
		//	// query := url.Query()
		//	// sort := query.Get("sort")
		//	// dir := query.Get("dir")
		//	// fmt.Printf("sort: %s , dir: %s \n", sort, dir)
		fmt.Println("returning all players")
		players, err := store.Players()
		if err != nil {
			msg := "Could not retrieve players"
			respondWithError(w, http.StatusNotFound, msg)
			return
		}
		respondWithJSON(w, http.StatusOK, players)
		return
	case "POST":
		// TODO - make route to create new players
		msg := "Desired method not yet implemented"
		respondWithError(w, http.StatusNotImplemented, msg)
		return
	default:
		msg := "Method not allowed"
		respondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func ShiftPath(p string) (head, tail string) {
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

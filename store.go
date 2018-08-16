package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var AlreadyDraftedErr = errors.New("the player was already drafted. pick again")


// Store interface contains all methods available for players
type Store interface {
	DraftPlayer(id int) (int64, error)
	Players() ([]*Player, error)
	Player(id int) (*Player, error)
	Scout() ([]Player, error)
}

type dbStore struct {
	db *sql.DB
}

// CreatePlayer is a function that creates a player and returns an error if there is one
func (store *dbStore) DraftPlayer(id int) (int64, error) {

	player, err := store.Player(id)
	if err != nil {
		return 0, err
	}

	fmt.Println("player is drafted: ", player.Drafted)
	if player.Drafted == true {
		err := AlreadyDraftedErr //errors.New("the player was already drafted. pick again")
		fmt.Println("player already drafted error: ", err)
		return 0, err
	}

	lastInsertID := 0

	sqlStatement := `
		UPDATE players 
		SET drafted = true
		WHERE id = $1
		RETURNING id`

	e := store.db.QueryRow(sqlStatement, id).Scan(&lastInsertID)

	switch {
	case e == sql.ErrNoRows:
		return 0, e
	case e != nil:
		fmt.Println(e)
		return 0, e
	default:
		fmt.Printf("Player is %d\n ", lastInsertID)
		return int64(lastInsertID), nil
	}
}

// Player is a function that returns player for an id (if player exists)
// https://golang.org/src/database/sql/example_test.go
func (store *dbStore) Player(id int) (*Player, error) {
	player := &Player{}

	err := store.db.QueryRow("SELECT * FROM players WHERE id = $1", id).Scan(&player.ID, &player.Name, &player.School, &player.Position, &player.Drafted)

	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("No player with that ID.")
		return nil, err
	case err != nil:
		//log.Fatal(err)
		fmt.Println(err)
		return nil, err
	default:
		fmt.Printf("Player is %s\n", player)
		return player, nil
	}
}

// Players is a function that returns all players
func (store *dbStore) Players() ([]*Player, error) {
	rows, err := store.db.Query("SELECT * FROM players")

	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil, err
	}
	defer rows.Close()

	players := []*Player{}

	for rows.Next() {
		player := &Player{}

		if err := rows.Scan(&player.ID, &player.Name, &player.School, &player.Position, &player.Drafted); err != nil {
			log.Fatal(err)
		}

		players = append(players, player)
	}

	return players, nil
}

// Scout is a function that returns all players that have not been drafted
func (store *dbStore) Scout() ([]Player, error) {
	rows, err := store.db.Query("SELECT * FROM players WHERE drafted = false")

	if err != nil {
		fmt.Println("Failed to run query", err)
		return nil, err
	}
	defer rows.Close()

	players := []Player{}

	for rows.Next() {
		player := &Player{}

		if err := rows.Scan(&player.ID, &player.Name, &player.School, &player.Position, &player.Drafted); err != nil {
			log.Fatal(err)
		}

		players = append(players, *player)
	}

	return players, nil
}

var store Store

// InitStore function to create new instance of the store
func InitStore(s Store) {
	store = s
}

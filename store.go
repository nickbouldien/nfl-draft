package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Store interface contains all methods available for players
type Store interface {
	DraftPlayer(id int) (bool, error)
	GetPlayers() ([]*Player, error)
	GetPlayer(id int) (*Player, error)
}

type dbStore struct {
	db *sql.DB
}

// CreatePlayer is a function that creates a player and returns an error if there is one
func (store *dbStore) DraftPlayer(id int) (bool, error) {
	_, err := store.db.Query("INSERT INTO players(full_name, position) VALUES ($1)", id)

	return false, err

}

// GetPlayer is a function that returns all players
// https://golang.org/src/database/sql/example_test.go
func (store *dbStore) GetPlayer(id int) (*Player, error) {
	player := &Player{}

	err := store.db.QueryRow("SELECT * FROM players WHERE id = $1", id).Scan(&player.ID, &player.Name, &player.School, &player.Position, &player.Drafted)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No player with that ID.")
		return nil, err
	case err != nil:
		log.Fatal(err)
		return nil, err
	default:
		fmt.Printf("Player is %s\n", player)
		return player, nil
	}
}

// GetPlayers is a function that returns all players
func (store *dbStore) GetPlayers() ([]*Player, error) {
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

var store Store

// InitStore function to create new instance of the store
func InitStore(s Store) {
	store = s
}

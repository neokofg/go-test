package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/neokofg/go-test/database"
	"log"
	"time"
)

type Anime struct {
	ID       int
	Title    string
	Rating   int32
	Episodes int32
}

func CreateAnime(db *sql.DB, anime *Anime) (int, error) {
	query := `INSERT INTO animes (title, rating, episodes) VALUES ($1, $2, $3) RETURNING id;`
	err := db.QueryRow(query, anime.Title, anime.Rating, anime.Episodes).Scan(&anime.ID)
	if err != nil {
		return 0, err
	}
	return anime.ID, nil
}

func GetAnime(db *sql.DB, id int) (*Anime, error) {
	anime := &Anime{}
	query := `SELECT * FROM animes WHERE id = $1;`
	err := db.QueryRow(query, id).Scan(&anime.ID, &anime.Title, &anime.Rating, &anime.Episodes)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

func GetAllAnimes(db *sql.DB) ([]Anime, error) {
	query := `SELECT * FROM animes;`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var animes []Anime
	for rows.Next() {
		var a Anime
		if err := rows.Scan(&a.ID, &a.Title, &a.Rating, &a.Episodes); err != nil {
			return nil, err
		}
		animes = append(animes, a)
	}
	return animes, nil
}

func UpdateAnime(db *sql.DB, anime *Anime) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
	}(tx)

	query := `UPDATE animes SET title = $1, rating = $2, episodes = $3 WHERE id = $4;`
	_, err = tx.Exec(query, anime.Title, anime.Rating, anime.Episodes, anime.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func main() {
	connStr := "host=localhost port=5432 user=postgres password=password dbname=sweetify sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to postgres")

	migrator, err := database.NewMigrator(db, "database/migrations")
	if err != nil {
		log.Fatal(err)
	}

	if err := migrator.Up(); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrated successfully")

	id, err := CreateAnime(db, &Anime{
		Title:    "New Anime",
		Rating:   1,
		Episodes: 12,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New Anime:", id)
}

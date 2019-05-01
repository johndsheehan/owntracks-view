package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8081"
)

type Location struct {
	Latitude   float32
	Longitude  float32
	TimeString string
}

type Locations struct {
	Location []Location
}

func mapRender(w http.ResponseWriter, r *http.Request) {
	log.Println("request from: ", r.RemoteAddr)
	locations := Locations{}

	err := dbQuery(&locations)
	if err != nil {
		log.Println("failed to query db: ", err)
		return
	}

	parsedTemplate, _ := template.ParseFiles("templates/mapTemplate.html")
	err = parsedTemplate.Execute(w, locations)
	if err != nil {
		log.Printf("failed to parse template: ", err)
		return
	}
}

func dbInit() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config(DBHost),
		Config(DBPort),
		Config(DBUser),
		Config(DBPass),
		Config(DBName))

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func dbQuery(locations *Locations) error {
	rows, err := db.Query(`
	SELECT latitude, longitude, timestamp
	FROM owntracks
	WHERE locationaccuarcy<5000
	ORDER BY timestamp DESC`)
	if err != nil {
		log.Println("query failed: ", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		location := Location{}
		err = rows.Scan(
			&location.Latitude,
			&location.Longitude,
			&location.TimeString,
		)
		if err != nil {
			log.Println("scan failed: ", err)
			return err
		}

		location.TimeString = strings.Replace(location.TimeString, "T", " ", -1)
		location.TimeString = strings.Replace(location.TimeString, "Z", "", -1)

		locations.Location = append(locations.Location, location)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := ConfigParse()
	if err != nil {
		log.Fatal("failed to parse config")
	}

	err = dbInit()
	if err != nil {
		log.Fatal("failed to init db: ", err)
	}
	defer db.Close()

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", mapRender)
	err = http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}
}

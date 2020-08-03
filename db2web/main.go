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
		log.Println("failed to parse template: ", err)
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
	WHERE locationaccuarcy<5000 AND latitude!=0 AND longitude!=0
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

func redirectToHTTPS(colonPort, colonTLSPort string) {
	go func(colonPort, colonTLSPort string) {
		redirect := func(w http.ResponseWriter, r *http.Request) {
			log.Printf("http request from: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

			host := r.Host
			if strings.HasSuffix(r.Host, colonPort) {
				host = strings.TrimSuffix(host, colonPort)
			}

			url := fmt.Sprintf("https://%s%s/%s", host, colonTLSPort, r.RequestURI)
			log.Printf("redirect to: %s", url)

			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}

		err := http.ListenAndServe(colonPort, http.HandlerFunc(redirect))
		if err != nil {
			log.Fatal(err)
		}
	}(colonPort, colonTLSPort)
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

	useTLS := false
	if Config(FullChain) != "" && Config(PrivateKey) != "" {
		useTLS = true
	}

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/", mapRender)

	if useTLS {
		colonPort := ":" + Config(HTTP)
		colonTLSPort := ":" + Config(HTTPS)
		go redirectToHTTPS(colonPort, colonTLSPort)

		log.Printf("serving https on %s\n", colonTLSPort)
		err := http.ListenAndServeTLS(colonTLSPort, Config(FullChain), Config(PrivateKey), nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		host := Config(Host) + ":" + Config(HTTP)

		log.Printf("serving http on %s", host)
		err = http.ListenAndServe(host, nil)
		if err != nil {
			log.Fatal("failed to start server: ", err)
		}
	}
}

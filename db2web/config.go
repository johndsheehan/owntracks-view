package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
)

const (
	DBHost = "DBHOST"
	DBPort = "DBPORT"
	DBUser = "DBUSER"
	DBPass = "DBPASS"
	DBName = "DBNAME"
)

var cfgMap map[string]string
var cfgMux = &sync.Mutex{}

func Config(name string) string {
	cfgMux.Lock()
	defer cfgMux.Unlock()

	return cfgMap[name]
}

func ConfigParse() error {
	cfgMap = make(map[string]string)

	parseCMDLine()
	err := parseEnv()

	return err
}

func parseCMDLine() {
	host := flag.String("host", "", "db hostname")
	port := flag.String("port", "", "db port")
	user := flag.String("user", "", "db user")
	pass := flag.String("pass", "", "db pass")
	name := flag.String("name", "", "db name")

	flag.Parse()

	cfgMap[DBHost] = *host
	cfgMap[DBPort] = *port
	cfgMap[DBUser] = *user
	cfgMap[DBPass] = *pass
	cfgMap[DBName] = *name
}

func argFetch(name, fallback string, failOnDefault, logValue bool) (string, error) {
	logMsg := ""
	value, ok := os.LookupEnv(name)
	if !ok {
		if cfgMap[name] != "" {
			value = cfgMap[name]
			logMsg = "setting " + name + " from cmdline: "
		} else {
			if failOnDefault {
				msg := "no default for: " + name
				return "", errors.New(msg)
			}

			value = fallback
			logMsg = "setting " + name + " to default: : "
		}
	} else {
		logMsg = "setting " + name + " from env: "
	}

	if logValue {
		logMsg += value
	}
	log.Println(logMsg)

	return value, nil
}

func parseEnv() error {
	host, err := argFetch(DBHost, "localhost", false, true)
	if err != nil {
		log.Println("failed to parse host config: ", err)
		return err
	}
	cfgMap[DBHost] = host

	port, err := argFetch(DBPort, "5432", false, true)
	if err != nil {
		log.Println("failed to parse port config: ", err)
		return err
	}
	cfgMap[DBPort] = port

	user, err := argFetch(DBUser, "", true, true)
	if err != nil {
		log.Println("failed to parse user: ", err)
		return err
	}
	cfgMap[DBUser] = user

	pass, err := argFetch(DBPass, "", true, false)
	if err != nil {
		log.Println("failed to parse pass: ", err)
		return err
	}
	cfgMap[DBPass] = pass

	name, err := argFetch(DBName, "", true, true)
	if err != nil {
		log.Println("failed to parse name: ", err)
		return err
	}
	cfgMap[DBName] = name

	return nil
}

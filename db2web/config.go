package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
)

const (
	DBHost     = "DBHOST"
	DBPort     = "DBPORT"
	DBUser     = "DBUSER"
	DBPass     = "DBPASS"
	DBName     = "DBNAME"
	Host       = "HOST"
	HTTP       = "HTTP"
	HTTPS      = "HTTPS"
	FullChain  = "FULLCHAIN"
	PrivateKey = "PRIVATEKEY"
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
	dbhost := flag.String("dbhost", "127.0.0.1", "db hostname (default 127.0.0.1")
	dbport := flag.String("dbport", "5432", "db port (default 5432)")
	dbuser := flag.String("dbuser", "", "db user")
	dbpass := flag.String("dbpass", "", "db pass")
	dbname := flag.String("dbname", "", "db name")

	host := flag.String("host", "127.0.0.1", "host (default: 127.0.0.1)")
	httpPort := flag.String("http", "5080", "http port (default 5080)")
	httpsPort := flag.String("https", "5443", "https port (default 5443)")
	fullChain := flag.String("fullchain", "", "fullchain.pem")
	privateKey := flag.String("privatekey", "", "privKey.pem")

	flag.Parse()

	cfgMap[DBHost] = *dbhost
	cfgMap[DBPort] = *dbport
	cfgMap[DBUser] = *dbuser
	cfgMap[DBPass] = *dbpass
	cfgMap[DBName] = *dbname

	cfgMap[Host] = *host
	cfgMap[HTTP] = *httpPort
	cfgMap[HTTPS] = *httpsPort
	cfgMap[FullChain] = *fullChain
	cfgMap[PrivateKey] = *privateKey
}

func argFetch(name string, failOnDefault, logValue bool) (string, error) {
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

			value = ""
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
	dbhost, err := argFetch(DBHost, false, true)
	if err != nil {
		log.Println("failed to parse host config: ", err)
		return err
	}
	cfgMap[DBHost] = dbhost

	dbport, err := argFetch(DBPort, false, true)
	if err != nil {
		log.Println("failed to parse port config: ", err)
		return err
	}
	cfgMap[DBPort] = dbport

	dbuser, err := argFetch(DBUser, true, true)
	if err != nil {
		log.Println("failed to parse user: ", err)
		return err
	}
	cfgMap[DBUser] = dbuser

	dbpass, err := argFetch(DBPass, true, false)
	if err != nil {
		log.Println("failed to parse pass: ", err)
		return err
	}
	cfgMap[DBPass] = dbpass

	dbname, err := argFetch(DBName, true, true)
	if err != nil {
		log.Println("failed to parse name: ", err)
		return err
	}
	cfgMap[DBName] = dbname

	host, err := argFetch(Host, false, true)
	if err != nil {
		log.Println("failed to parse host: ", err)
		return err
	}
	cfgMap[Host] = host

	httpPort, err := argFetch(HTTP, false, true)
	if err != nil {
		log.Println("failed to parse http: ", err)
		return err
	}
	cfgMap[HTTP] = httpPort

	httpsPort, err := argFetch(HTTPS, false, true)
	if err != nil {
		log.Println("failed to parse  : ", err)
		return err
	}
	cfgMap[HTTPS] = httpsPort

	fullChain, err := argFetch(FullChain, false, true)
	if err != nil {
		log.Println("failed to parse fullchain: ", err)
		return err
	}
	cfgMap[FullChain] = fullChain

	privateKey, err := argFetch(PrivateKey, false, true)
	if err != nil {
		log.Println("failed to parse privateKey: ", err)
		return err
	}
	cfgMap[PrivateKey] = privateKey

	return nil
}

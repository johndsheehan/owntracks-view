package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"sync"
)

const (
	MQTTBroker   = "MQTTBROKER"
	MQTTTopic    = "MQTTTOPIC"
	MQTTClientID = "MQTTCLIENTID"
	MQTTUser     = "MQTTUSER"
	MQTTPass     = "MQTTPASS"

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

func parseCMDLine() {
	mqttBroker := flag.String("mqtt-broker", "", "mqtt broker")
	mqttTopic := flag.String("mqtt-topic", "", "mqtt topic")
	mqttClientID := flag.String("mqtt-client-id", "", "mqtt client id")
	mqttUser := flag.String("mqtt-user", "", "mqtt username")
	mqttPass := flag.String("mqtt-pass", "", "mqtt password")

	dbHost := flag.String("db-host", "", "db hostname")
	dbPort := flag.String("db-port", "", "db port")
	dbUser := flag.String("db-user", "", "db user")
	dbPass := flag.String("db-pass", "", "db pass")
	dbName := flag.String("db-name", "", "db name")

	flag.Parse()

	cfgMap[MQTTBroker] = *mqttBroker
	cfgMap[MQTTTopic] = *mqttTopic
	cfgMap[MQTTClientID] = *mqttClientID
	cfgMap[MQTTUser] = *mqttUser
	cfgMap[MQTTPass] = *mqttPass

	cfgMap[DBHost] = *dbHost
	cfgMap[DBPort] = *dbPort
	cfgMap[DBUser] = *dbUser
	cfgMap[DBPass] = *dbPass
	cfgMap[DBName] = *dbName
}

func parseEnv() error {
	mqttBroker, err := argFetch(MQTTBroker, "localhost:1883", false, true)
	if err != nil {
		log.Println("failed to parse mqtt broker: ", err)
		return err
	}
	cfgMap[MQTTBroker] = mqttBroker

	mqttTopic, err := argFetch(MQTTTopic, "", true, true)
	if err != nil {
		log.Println("failed to parse mqtt topic: ", err)
		return err
	}
	cfgMap[MQTTTopic] = mqttTopic

	mqttClientID, err := argFetch(MQTTClientID, "", true, true)
	if err != nil {
		log.Println("failed to parse mqtt client id: ", err)
		return err
	}
	cfgMap[MQTTClientID] = mqttClientID

	mqttUser, err := argFetch(MQTTUser, "", false, true)
	if err != nil {
		log.Println("failed to parse mqtt user: ", err)
		return err
	}
	cfgMap[MQTTUser] = mqttUser

	mqttPass, err := argFetch(MQTTPass, "", false, true)
	if err != nil {
		log.Println("failed to parse mqtt pass: ", err)
		return err
	}
	cfgMap[MQTTPass] = mqttPass

	dbHost, err := argFetch(DBHost, "localhost", false, true)
	if err != nil {
		log.Println("failed to parse host config: ", err)
		return err
	}
	cfgMap[DBHost] = dbHost

	dbPort, err := argFetch(DBPort, "5432", false, true)
	if err != nil {
		log.Println("failed to parse port config: ", err)
		return err
	}
	cfgMap[DBPort] = dbPort

	dbUser, err := argFetch(DBUser, "", true, true)
	if err != nil {
		log.Println("failed to parse user: ", err)
		return err
	}
	cfgMap[DBUser] = dbUser

	dbPass, err := argFetch(DBPass, "", true, false)
	if err != nil {
		log.Println("failed to parse pass: ", err)
		return err
	}
	cfgMap[DBPass] = dbPass

	dbName, err := argFetch(DBName, "", true, true)
	if err != nil {
		log.Println("failed to parse name: ", err)
		return err
	}
	cfgMap[DBName] = dbName

	return nil
}

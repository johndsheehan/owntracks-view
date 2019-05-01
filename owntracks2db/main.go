/*
CREATE DATABASE test;
\c test

CREATE TABLE owntracks (
  uid                 serial      NOT NULL,
  altitude            int         NOT NULL,
  altitudeAccuarcy    int         NOT NULL,
  battery             int         NOT NULL,
  barometricPressure  float       NOT NULL,
  connectivityStatus  char        NOT NULL,
  connectionTrigger   char        NOT NULL,
  locationAccuarcy    int         NOT NULL,
  longitude           float       NOT NULL,
  latitude            float       NOT NULL,
  timestamp           timestamp   NOT NULL,
  tid                 varchar(2)  NOT NULL,
  CONSTRAINT userinfo_pkey PRIMARY KEY (uid)
)
WITH (OIDS=FALSE);

*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
)

var msgs = make(chan MQTT.Message)

func mqttMsgHandler(client MQTT.Client, msg MQTT.Message) {
	msgs <- msg
}

func owntracksStore(topic string, msg []byte) {
	type ots struct {
		Altitude           int     `json:"alt"`
		AltitudeAccuarcy   int     `json:"vac"`
		Battery            int     `json:"batt"`
		BarometricPressure float32 `json:"p"`
		ConnectivityStatus string  `json:"conn"`
		ConnectionTrigger  string  `json:"t"`
		LocationAccuarcy   int     `json:"acc"`
		Latitude           float32 `json:"lat"`
		Longitude          float32 `json:"lon"`
		Timestamp          int64   `json:"tst"`
		TrackerID          string  `json:"tid"`
	}

	var ot ots
	err := json.Unmarshal(msg, &ot)
	if err != nil {
		log.Print("failed json unmarshall: ", err)
		log.Print(string(msg))
		return
	}
	fmt.Printf("%+v\n", ot)

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config(DBHost),
		Config(DBPort),
		Config(DBUser),
		Config(DBPass),
		Config(DBName))
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Print("db error: ", err)
		return
	}
	defer db.Close()

	unixTimeUTC := time.Unix(ot.Timestamp, 0)
	timeRFC3339 := unixTimeUTC.Format(time.RFC3339)

	var lastInsertID int
	stmt := "INSERT INTO owntracks(altitude, altitudeAccuarcy, battery, barometricPressure, connectivityStatus, connectionTrigger, locationAccuarcy, latitude, longitude, timestamp, tid) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning uid;"
	err = db.QueryRow(stmt,
		ot.Altitude,
		ot.AltitudeAccuarcy,
		ot.Battery,
		ot.BarometricPressure,
		ot.ConnectivityStatus,
		ot.ConnectionTrigger,
		ot.LocationAccuarcy,
		ot.Latitude,
		ot.Longitude,
		timeRFC3339,
		ot.TrackerID,
	).Scan(&lastInsertID)
	if err != nil {
		log.Print("db error: ", err)
		return
	}
	log.Println("last insert id: ", lastInsertID)
}

func main() {
	err := ConfigParse()
	if err != nil {
		log.Fatal("failed to parse config: ", err)
	}

	mqttOpts := MQTT.NewClientOptions()
	mqttOpts.AddBroker(Config(MQTTBroker))
	mqttOpts.SetClientID(Config(MQTTClientID))
	mqttOpts.SetCleanSession(true)
	mqttOpts.SetUsername(Config(MQTTUser))
	mqttOpts.SetPassword(Config(MQTTPass))
	mqttOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(Config(MQTTTopic), 0, mqttMsgHandler); token.Wait() && token.Error() != nil {
			log.Fatal("subscribe failed: ", token.Error())
		}
	}

	client := MQTT.NewClient(mqttOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("connect failed: ", token.Error())
	}

	for {
		select {
		case msg := <-msgs:
			if !msg.Duplicate() {
				go owntracksStore(msg.Topic(), msg.Payload())
			}
		}
	}
}

[Owntracks](https://owntracks.org) is a mobile phone application for logging location information.
It can also publish this information to a `MQTT` broker.
`owntracks2db` subscribes to the broker and writes any location data recieved to a postgres db.
`db2web` is a simple web server which plots the information to a webpage.

The schema used for the database is,
```
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
```

## to build
```
cd owntracks2db
go get
go build
# CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .  # for a static build
```

```
cd db2web
go get
go build
# CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .  # for a static build
```

## to run

**Note:** default argument values are overwritten by command line argument values which are overwritten by enviroment variables. Arguments which do not have a default value, and receive no value from either the command line or enviroment will trigger an error.

`owntracks2db`

```
cd owntracks2db
./owntracks2db  --mqtt-broker 'REPLACE'  --mqtt-topic 'REPLACE'  --mqtt-client-id 'REPLACE'  --mqtt-user 'REPLACE'  --mqtt-pass 'REPLACE'  --db-user 'REPLACE'  --db-pass 'REPLACE'  --db-name 'REPLACE'
```
the corresponding enviroment variables are,
```
MQTTBROKER  # default is localhost:1883
MQTTTOPIC
MQTTCLIENTID
MQTTUSER
MQTTPASS

DBHOST  # default is localhost
DBPORT  # default is 5432
DBUSER
DBPASS
DBNAME
```

`db2web`

```
./db2web --host 'REPLACE'  --port REPLACE  --name 'REPLACE'  --user 'REPLACE'  --pass 'REPLACE'
```
the corresponding enviroment variables are
```
DBHOST  # default is localhost
DBPORT  # default is 5432
DBUSER
DBPASS
DBNAME
```
the map is served on `http://localhost:8081`
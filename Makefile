.DEFAULT_GOAL := build

build: clean
	cd db2web ;  go get -u ;  go build
	cd owntracks2db ;  go get -u ;  go build

clean:
	rm -f db2web/db2web
	rm -f owntracks2db/owntracks2db

install:
	cp  db2web/db2web  /usr/local/bin
	cp  owntracks2db/owntracks2db  /usr/local/bin

uninstall:
	rm  /usr/local/bin/db2web
	rm  /usr/local/bin/owntracks2db

rpi: clean
	cd db2web ;  go get -u ;  env GOOS=linux GOARCH=arm GOARM=5 go build
	cd owntracks2db ;  go get -u ;  env GOOS=linux GOARCH=arm GOARM=5 go build

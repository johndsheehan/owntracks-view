FROM golang:alpine AS build

RUN apk update && apk add  --no-cache  git 

WORKDIR /app
COPY go.mod go.sum *.go /app/

ENV GO111MODULE="on"
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .


FROM scratch
COPY --from=build /app/db2web /db2web

ENTRYPOINT ["/db2web"]

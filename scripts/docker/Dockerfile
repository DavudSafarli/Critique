FROM golang:1.15.3

#RUN apk add --no-cache curl gcc libc-dev

WORKDIR /app

COPY ../.. .

RUN go get -d -v ./...

#RUN go install -v ./...
#RUN migrate -source "file://./scripts/db/migrations" -database "postgres://myuser:mypass@test-db/mydatabase?sslmode=disable" up
ENV CGO_ENABLED 0
#CMD sh -c "CGO_ENABLED=0 go test ./external/repository/postgres_repos    -v"
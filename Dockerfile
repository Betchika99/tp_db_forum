FROM golang:latest AS builder

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build main.go

FROM ubuntu:18.10

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y postgresql-10

EXPOSE 5000

USER postgres

RUN service postgresql start &&\
    psql --command="CREATE USER forum WITH NOSUPERUSER PASSWORD 'forum';" &&\
    createdb --owner=forum mydb &&\
    psql --dbname=mydb --command="CREATE EXTENSION IF NOT EXISTS citext;" &&\
    service postgresql stop

WORKDIR app
COPY --from=builder /usr/src/app .

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/10/main/pg_hba.conf
RUN echo "listen_addresses = '*'\nsynchronous_commit = off\nfsync = off\nunix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/10/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD service postgresql start && ./main

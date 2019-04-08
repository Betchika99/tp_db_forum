package database

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
)

const psqlURI = "postgresql://forum:forum@localhost:5432/mydb"
const dbSchema = "./database/db_scheme.sql"
const maxConn = 8

var Connection *pgx.ConnPool

func OpenConnect() error {
	configFromURI, err := pgx.ParseURI(psqlURI)
	if err != nil {
		return err
	}

	connConfig := pgx.ConnPoolConfig{ // type pgx.ConnPoolConfig
		ConnConfig: configFromURI,
		MaxConnections: maxConn,
	}

	conn, err := pgx.NewConnPool(connConfig)
	if err != nil {
		return err
	}

	SetConnect(conn)
	log.Printf("DB Connection opened")

	//go listen()

	return nil
}

//func listen() {
//	log.Println(Connection)
//	conn, err := Connection.Acquire()
//	if err != nil {
//		log.Println("ERROR is", err.Error())
//		return
//	}
//	defer Connection.Release(conn)
//
//	conn.Listen("chat")
//
//	for {
//		notification, err := conn.WaitForNotification(context.Background())
//		if err != nil {
//			fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
//			os.Exit(1)
//		}
//
//		fmt.Println("PID:", notification.PID, "Channel:", notification.Channel, "Payload:", notification.Payload)
//	}
//
//}

func SetConnect(connNew *pgx.ConnPool) {
	Connection = connNew
}

func GetConnect() *pgx.ConnPool{
	return Connection
}

//func Connect(psqURI string) (*pgx.ConnPool, error) {
//	config, err := pgx.ParseURI(psqURI)
//	if err != nil {
//		return nil, err
//	}
//	conn, err := pgx.NewConnPool(
//		pgx.ConnPoolConfig{
//			ConnConfig: config,
//			MaxConnections: maxConn,
//		})
//	if err != nil {
//		return nil, err
//	}
//	return conn, nil
//}

//var Connection, _ = Connect(PsqlURI)

//func MakeTransaction() *pgx.Tx {
//	conn, _ := Connection.Begin()
//	return conn
//}

func CloseConnect() {
	Connection.Close()
}

func LoadSchema() error {
	transaction, err := Connection.Begin()
	if err != nil {
		log.Println("ERROR is", err.Error())
		return err
	}

	content, err := ioutil.ReadFile(dbSchema)
	if err != nil {
		log.Println("ERROR is", err.Error())
		return err
	}

	if _, err = transaction.Exec(string(content)); err != nil {
		log.Println(err)
		if err2 := transaction.Rollback(); err2 != nil {
			log.Println("ERROR is", err2.Error())
			return err2
		}
		return err
	}

	if err = transaction.Commit(); err != nil {
		log.Println("ERROR is", err.Error())
		return err
	}

	log.Printf("DB schema loaded")
	return nil
}

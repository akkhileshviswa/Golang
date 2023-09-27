package main

import (
	"Golang/newsletter/user-service/event"
	"database/sql"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
	"usermanagement/data"

	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/lib/pq"
)

const webPort = "8081"

type RPCServer struct{}

type Config struct {
	DB     *sql.DB
	Models data.Models
}

const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "root"
	dbname   = "users"
)

func main() {
	log.Println("Starting user service")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	go app.rpcListen()
	rabbitConn, err := rabbitConnect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"UpdateRecord"}, rabbitConn)
	if err != nil {
		log.Println(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// This function is used to listen to rpc calls.
func (app *Config) rpcListen() error {
	server := new(RPCServer)
	rpc.Register(server)

	log.Println("Starting RPC server on port 8083")
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

// This function is used to connect to postgres database.
func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Established a successful connection!")
	return db
}

// This function is used to connect to rabbit mq.
func rabbitConnect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

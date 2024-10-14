package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo/Database"
	"todo/Server"

	"github.com/sirupsen/logrus"
)

const shutdownTimeout = 20 * time.Second

func Home(w http.ResponseWriter, r *http.Request) {
	// handle error
	w.Write([]byte("Welcome to the app - Home Page of the todo "))
}

func main() {

	//channel to do the task
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Welcome to the app - ToDo")

	connect := fmt.Sprintf("host %s port %s ", os.Getenv("DB_HOST"), os.Getenv(""))
	fmt.Println(connect)

	//server instance
	serv := Server.SetupRoutes()
	if err := Database.ConnectAndMigrate(
		//dbHost, dbPort, dbName, dbUser, dbPassword,

		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		Database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initilize and migrate database with err %+v", err)
	}
	logrus.Printf("migration sucessfull")

	go func() {
		if err := serv.Run(":8000"); err != nil {
			logrus.Panicf("Failed to run server with err %+v", err)
		}
	}()

	logrus.Printf("Server started at : 8000")

	<-done

	logrus.Printf("Server Shutdown")

}

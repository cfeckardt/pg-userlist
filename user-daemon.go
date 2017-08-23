package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

func doWork(db *sql.DB, work int64) {

}

func getUsers(db *sql.DB) {
}

func updateUsers(db *sql.DB) {

}

func waitForNotification(l *pq.Listener) {
	select {
	case <-l.Notify:
		fmt.Println("Received notification")
	case <-time.After(90 * time.Second):
		go l.Ping()
	}
}

func main() {
	args := os.Args[1:]
	var conninfo = strings.Join(args, " ")

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen("users_changed")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening for new users")

	for {
		updateUsers(db)
		waitForNotification(listener)
	}
}

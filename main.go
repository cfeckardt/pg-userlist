package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

var delimeter *string
var listen *string

func getUsers(db *sql.DB) {
	var rows, err = db.Query(`SELECT rolname, rolpassword FROM pg_authid WHERE rolpassword LIKE 'md5%';`)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	names := make([]string, 0)

	for rows.Next() {
		var name string
		var pass string
		err = rows.Scan(&name, &pass)
		if err != nil {
			panic(err)
		}
		names = append(names, name + " " + pass)
	}

	fmt.Println(strings.Join(names, *delimeter))
}

func waitForNotification(l *pq.Listener, db *sql.DB) {
	select {
	case <-l.Notify:
		fmt.Println("Received notification")
		getUsers(db)
	case <-time.After(90 * time.Second):
		go l.Ping()
	}
}

func main() {
	delimeter = flag.String("d", "\n", "Delimiter character")
	listen = flag.String("l", "", "")
	flag.Parse()
	var conninfo = strings.Join(flag.Args(), " ")

	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		panic(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}

	getUsers(db)

	if len(*listen) > 0 {
		listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
		err = listener.Listen("users_changed")
		if err != nil {
			panic(err)
		}

		for {
			waitForNotification(listener, db)
		}
	}
}

package main

import (
  "database/sql"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "time"

  "github.com/lib/pq"
)

var delimeter *string
var separator *string
var quoteOutput *bool
var listen *string
var saveFileName *string

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
    if *quoteOutput {
      name = fmt.Sprintf("\"%s\"", name)
      pass = fmt.Sprintf("\"%s\"", pass)
    }
    names = append(names, name + *separator + pass)
  }

  out := strings.Join(names, *delimeter)
  if len(*saveFileName) > 0 {
    saveFile(out)
  } else {
    fmt.Println(out)
  }
}

func saveFile(foo string) {
  fmt.Println(foo)
  d1 := []byte(foo)

  err := ioutil.WriteFile(*saveFileName, d1, 0644)
  if err != nil {
    panic(err)
  }
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
  separator = flag.String("s", ":", "Field separator/delimeter")
  quoteOutput = flag.Bool("q", false, "Quote user and password strings in output")
  saveFileName = flag.String("f", "", "Filename to save output to")
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
    err = listener.Listen(*listen)
    if err != nil {
      panic(err)
    }

    for {
      waitForNotification(listener, db)
    }
  }
}

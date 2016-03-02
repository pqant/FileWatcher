package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	_ "github.com/denisenkom/go-mssqldb"
	_ "database/sql"
	"flag"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "AAABBBCCC", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "localhost", "the database server")
var user = flag.String("user", "sa", "the database user")

func main() {
	//fmt.Printf("Sql is testing!\n")

	//flag.Parse() // parse the command line args
	//
	//if *debug {
	//	fmt.Printf(" password:%s\n", *password)
	//	fmt.Printf(" port:%d\n", *port)
	//	fmt.Printf(" server:%s\n", *server)
	//	fmt.Printf(" user:%s\n", *user)
	//}
	//
	//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	//if *debug {
	//	fmt.Printf(" connString:%s\n", connString)
	//}
	//conn, err := sql.Open("mssql", connString)
	//if err != nil {
	//	log.Fatal(">>> Open connection failed:\n", err.Error())
	//}
	//defer conn.Close()
	//
	//stmt, err := conn.Prepare("select 1, 'abc'")
	//if err != nil {
	//	log.Fatal(">>> Prepare failed:\n", err.Error())
	//}
	//defer stmt.Close()
	//
	//row := stmt.QueryRow()
	//var somenumber int64
	//var somechars string
	//err = row.Scan(&somenumber, &somechars)
	//if err != nil {
	//	log.Fatal("Scan failed:", err.Error())
	//}
	//fmt.Printf("somenumber:%d\n", somenumber)
	//fmt.Printf("somechars:%s\n", somechars)
	//
	//fmt.Printf("bye\n")
	//
	//return

	fmt.Println("Files are listening... \n")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("testDir")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()

}

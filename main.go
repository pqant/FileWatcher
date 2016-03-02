package main

import (
	"fmt"
	"log"
	"flag"
	_"github.com/pqant/FileWatcher/Ftp"
	_ "github.com/denisenkom/go-mssqldb"
	_"time"
	"github.com/pqant/FileWatcher/SqlUtility"
	"github.com/howeyc/fsnotify"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "!1q2w3e!", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "172.16.56.129", "the database server")
var user = flag.String("user", "sa", "the database user")
var database = flag.String("database", "FILEWATCHER", "database name")

func main() {

	/* Sample Comment by Windows! Selam !!*/
	flag.Parse()

	SqlUtility.ShowConnectionInfo(&SqlUtility.DbFlagContainer{
		Debug:debug,
		Server:server,
		User:user,
		Password:password,
		Port:port,
		Database:database,
	})

	fmt.Println("File system watcher is listening... \n")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				{
					go func() {
						if (ev.IsCreate()) {
							conn := SqlUtility.OpenConnection(debug, server, user, password, port, database)
							defer conn.Close()
							SqlUtility.SendToSql(conn, ev.Name, "CREATE")
						} else if (ev.IsDelete()) {
							conn := SqlUtility.OpenConnection(debug, server, user, password, port, database)
							defer conn.Close()
							SqlUtility.SendToSql(conn, ev.Name, "DELETE")
						} else {
							log.Printf("ev -> %v %v %v %v \n",ev.IsCreate(),ev.IsDelete(),ev.IsModify(),ev.IsRename())
						}
					}()
				}
			case err := <-watcher.Error:
				{
					log.Println("Error:", err)
				}
			}
		}
	}()

	err = watcher.Watch("TestDirectory")
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()

	/*

	fmt.Printf("FTP Connection is testing!\n")
	result := Ftp.FtpCheck("Pragmalinq","....","ftp.yapikredi.com.tr",21)
	fmt.Printf("FTP CONNECTION RESULT : %v\n",result)
	return

	*/

}

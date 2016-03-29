package main

import (
	"fmt"
	"log"
	"flag"
	_"github.com/pqant/FileWatcher/Ftp"
	_ "github.com/denisenkom/go-mssqldb"
	_"time"
	"github.com/pqant/FileWatcher/SqlUtility"
	"github.com/fsnotify/fsnotify"
	"os"
)

var debug = flag.Bool("debug", false, "enable debugging")
var password = flag.String("password", "___", "the database password")
var port *int = flag.Int("port", 1433, "the database port")
var server = flag.String("server", "11.1.1.1", "the database server")
var user = flag.String("user", "sa", "the database user")
var database = flag.String("database", "FILEWATCHER", "database name")

func main() {

	/* Sample Comment by Windows! Selam ! 22 !*/
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
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				{
					go func() {

						if ev.Op & fsnotify.Create == fsnotify.Create {
							conn := SqlUtility.OpenConnection(debug, server, user, password, port, database)
							defer conn.Close()
							SqlUtility.SendToSql(conn, ev.Name, "CREATE")
						} else if ev.Op & fsnotify.Remove == fsnotify.Remove {
							conn := SqlUtility.OpenConnection(debug, server, user, password, port, database)
							defer conn.Close()
							SqlUtility.SendToSql(conn, ev.Name, "DELETE")
						}
						log.Printf("ev -> %v - %v \n", ev.Name, ev.String())
					}()
				}
			case err := <-watcher.Errors:
				{
					log.Println("Error:", err)
				}
			}
		}
	}()


	path := "TestDirectory"
	if _, err := Exists(path); err == nil {
		fmt.Printf("File or Directory exists!\n")
	} else {
		fmt.Printf("Error >>> %s", err.Error())
	}

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add("TestDirectory")
	if err != nil {
		log.Fatal(err)
	}

	<-done

	/*

	fmt.Printf("FTP Connection is testing!\n")
	result := Ftp.FtpCheck("Pragmalinq","....","ftp.yapikredi.com.tr",21)
	fmt.Printf("FTP CONNECTION RESULT : %v\n",result)
	return

	*/

}


func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

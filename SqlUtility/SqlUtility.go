package SqlUtility

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/pqant/FileWatcher/Utility"
)

type DbFlagContainer struct {
	Debug    *bool
	Server   *string
	User     *string
	Password *string
	Port     *int
	Database *string
}

func OpenConnection(debug *bool, server *string, user *string, password *string, port *int, database *string) *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", *server, *user, *password, *port, *database)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal(">>> Open connection failed:\n", err.Error())
		os.Exit(-1)
	} else {
		fmt.Print("Connection is successfull!\n")
	}
	return conn
}

func ShowConnectionInfo(container *DbFlagContainer) {
	if *container.Debug {
		fmt.Print(">> Sql Connection Informations !<< \n")
		fmt.Printf("Password:%s\n", *container.Password)
		fmt.Printf("Port:%d\n", *container.Port)
		fmt.Printf("Server:%s\n", *container.Server)
		fmt.Printf("User:%s\n", *container.User)
		fmt.Printf("Database:%s\n", *container.Database)
	}
}

func SimpleSelections(conn *sql.DB) {
	stmt, err := conn.Prepare("select 1, 'abc'")
	if err != nil {
		log.Fatal(">>> Prepare failed:\n", err.Error())
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	var someNumber int64
	var someChars string
	err = row.Scan(&someNumber, &someChars)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Printf("somenumber:%d\n", someNumber)
	fmt.Printf("somechars:%s\n", someChars)
}

func BulkInsertionSample(conn *sql.DB) {
	controlFlag := make(chan bool)
	go func() {
		for i := 0; i < 1000; i++ {
			insertText := fmt.Sprintf("insert into Transactions (SENDER_HOST,TIME_STAMP,FILE_NAME,FILE_EVENT) values ('%s','%s','%s','%s')", "mac", time.Now().String()[:19],
				"TEMP_FILE.txt", "CREATE")

			stmt, err := conn.Prepare(insertText)
			if err != nil {
				log.Fatal(">>> Prepare failed:\n", err.Error())
			}
			defer stmt.Close()

			if result, err := stmt.Exec(); err != nil {
				log.Fatal(">>> Execution Error :\n", err.Error())
			} else {
				affectedRows, _ := result.RowsAffected()
				insertId, _ := result.LastInsertId()
				fmt.Printf("Affected Row(s) : %d , Last Inserted Row Id : %d\n", affectedRows, insertId)
			}
		}
		fmt.Println("Done!\n")
		controlFlag <- true
	}()
	<-controlFlag
}

func SendToSql(conn *sql.DB,fileName,eventType string) bool{
	hostName := Utility.HostName()
	insertText := fmt.Sprintf("insert into Transactions (SENDER_HOST,TIME_STAMP,FILE_NAME,FILE_EVENT) values ('%s','%s','%s','%s')", hostName, time.Now().String()[:19],
		fileName, eventType)
	stmt, err := conn.Prepare(insertText)
	if err != nil {
		log.Fatal(">>> Prepare failed:\n", err.Error())
	}
	defer stmt.Close()

	if result, err := stmt.Exec(); err != nil {
		log.Fatal(">>> Execution Error :\n", err.Error())
	} else {
		affectedRows, _ := result.RowsAffected()
		insertId, _ := result.LastInsertId()
		fmt.Printf("Affected Row(s) : %d , Last Inserted Row Id : %d\n", affectedRows, insertId)
		return true
	}
	return false
}



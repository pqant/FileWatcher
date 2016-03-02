package Ftp

import (
	"github.com/smallfish/ftp"
	"fmt"
	"os"
)

func FtpCheck(userName,password,ftpHost string,portNo int) bool {
	retValue := false
	if userName=="" {
		return retValue
	}
	if password=="" {
		return retValue
	}
	if ftpHost=="" {
		return retValue
	}

	ftp := new(ftp.FTP)

	if (ftp!=nil) {
		fmt.Println("Burda 3!\n")

	}
	// debug default false
	ftp.Debug = true
	ftp.Connect(ftpHost, portNo)

	// login
	ftp.Login(userName,password)
	if ftp.Code == 530 {
		fmt.Println("error: login failure")
		retValue = false
		os.Exit(-1)
	}
	retValue = true

	// pwd
	ftp.Pwd()
	fmt.Println("code:", ftp.Code, ", message:", ftp.Message)

	/*
	// make dir
	ftp.Mkd("/path")
	ftp.Request("TYPE I")

	// stor file
	b, _ := ioutil.ReadFile("/path/a.txt")
	ftp.Stor("/path/a.txt", b)
	*/

	ftp.Quit()
	return retValue
}

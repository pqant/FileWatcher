package Utility

import "os"

func HostName() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
}
/*
A web based application to manage .ssh/config files. This is the backend. The frontend is a seperate application.
*/
package main

import (
	"fmt"
	"sshconfigmanager/configreader"
)

func main() {
	sshConfig := configreader.NewSshConfig()
	if err := sshConfig.ReadConfig(); err != nil {
		fmt.Printf("Error reading config. Error: %s\n", err.Error())
		return
	}

	sshConfig.PrintCurrentConfig()
}

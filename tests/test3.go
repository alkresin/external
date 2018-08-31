package main

import (
	"fmt"
	egui "github.com/alkresin/external"
)

func main() {

	s := "guiserver=\naddress=95.80.77.1\nport=2801"
	if egui.Init(s) {
		fmt.Println("Ok")
	} else {
		fmt.Println("No luck")
	}

	egui.Exit()
}

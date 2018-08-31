package main

import (
	egui "github.com/alkresin/external"
)

func main() {

	if !egui.Init("") {
		return
	}

	f := func (p []string)string {
		if p == nil {
		}
		egui.EvalProc("oLabel1:SetText(\"Hi, friends!\")")
		return ""
	}
	egui.RegFunc("fmenu1", f)

	egui.OpenMainForm("c:/papps/utils/tcpip/forms/example.xml")

	egui.Exit()

}


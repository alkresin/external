package main

import (
	egui "github.com/alkresin/external"
)

func main() {

	if !egui.Init("") {
		return
	}

	f := func (p []string)string {
		egui.EvalProc("oLabel1:SetText(\"Hi, friends! (" + p[0] + ")\")")
		return ""
	}
	egui.RegFunc("fmenu1", f)

	f = func (p []string)string {
		if p == nil {
		}
		return "Hi from Go!"
	}
	egui.RegFunc("fmenu2", f)

	egui.OpenMainForm("forms/example.xml")

	egui.Exit()

}


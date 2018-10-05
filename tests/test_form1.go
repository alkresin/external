// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a GNU general public
// license that can be found in the LICENSE file.

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

	f = func ([]string)string {	return "Hi from Go!" }
	egui.RegFunc("fmenu2", f)

	egui.OpenMainForm("forms/example.xml")

	egui.Exit()

}


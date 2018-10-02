package main

import egui "github.com/alkresin/external"

func main() {

	if !egui.Init("") {
		return
	}

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 400, H: 280, Title: "External"})
	egui.InitMainWindow(pWindow)

	egui.Menu("")
	egui.Menu( "File" )
	egui.AddMenuItem( "Open form", openf, "openf" )
	egui.AddMenuSeparator()
	egui.AddMenuItem( "Exit", nil, "hwg_EndWindow()" )
	egui.EndMenu()
	egui.Menu( "Help" )
	egui.AddMenuItem( "About", nil, "hwg_MsgInfo(hb_version()+chr(10)+chr(13)+hwg_version(),\"About\")" )
	egui.EndMenu()
	egui.EndMenu()

	pWindow.Activate()

	egui.Exit()
}

func openf([]string)string {

	egui.OpenForm("forms/testget2.xml")
	return ""
}

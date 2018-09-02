package main

import (
	"fmt"
	egui "github.com/alkresin/external"
)

var pLabel *egui.Widget
var pEdi1 *egui.Widget

func main() {

	if !egui.Init("") {
		return
	}

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 400, H: 220, Title: "External"})
	egui.InitMainWindow(pWindow)

	egui.Menu("")
	egui.Menu( "File" )
	egui.AddMenuItem( "New",
		func (p []string)string { pLabel.SetText(p[0]); return "" }, "fsett2", "Bye...1" )
	egui.AddMenuItem( "Open dialog", fsett3, "fsett3" )
	egui.AddMenuSeparator()
	egui.AddMenuItem( "Message box", fmbox1, "fmbox1" )
	egui.EndMenu()
	egui.Menu( "Help" )
	egui.AddMenuItem( "About", nil, "hwg_MsgInfo(\"Test\",\"About\")" )
	egui.EndMenu()
	egui.EndMenu()

	pLabel = pWindow.AddWidget(&(egui.Widget{Type: "label", Name: "l1",
		X: 20, Y: 20, W: 180, H: 24, Title: "Test of a label",
		AProps: map[string]string{"Transpa":"t"} }))

	pWindow.AddWidget(&(egui.Widget{Type: "label",
		X: 20, Y: 50, W: 180, H: 24, Title: "Second", TColor: 255,
		AProps: map[string]string{"Transpa":"t"} }))

	pWindow.AddWidget(&(egui.Widget{Type: "button", X: 200, Y: 16, W: 100, H: 32, Title: "Click"}))
	egui.PLastWidget.SetCallBackProc("bclick", nil, "private sss:=\"Done\"\r\nhwg_MsgInfo(sss)")

	pWindow.AddWidget(&(egui.Widget{Type: "button", X: 200, Y: 60, W: 100, H: 32, Title: "SetText"}))
	egui.PLastWidget.SetCallBackProc("bclick", fsett1, "fsett1", "first parameter")

	pWindow.Activate()

	egui.Exit()

}

func fsett1(p []string)string {

	pLabel.SetText( p[1] )
	//b := egui.EvalFunc( "Return GetWidgetByName(\"main.l1\"):GetText()")
	s := pLabel.GetText()
	fmt.Println( s )
	return ""
}


func fsett3(p []string)string {
	if p == nil {}

	pFont := egui.CreateFont( &(egui.Font{Name: "f1", Family: "Georgia", Height: 16}) )
	pDlg := &(egui.Widget{X: 300, Y: 200, W: 200, H: 370, Title: "Dialog Test", Font: pFont })
	egui.InitDialog(pDlg)

	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 20, W: 180, H: 24, Title: "Name:"}))
	pEdi1 = pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi1", X: 20, Y: 44, W: 160, H: 24 }))
	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 80, W: 180, H: 24, Title: "SurName:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi2", X: 20, Y: 104, W: 160, H: 24 }))
	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 140, W: 180, H: 24, Title: "Профессия:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi3", X: 20, Y: 164, W: 160, H: 24 }))

	pDlg.AddWidget(&(egui.Widget{Type: "button", X: 50, Y: 330, W: 100, H: 32, Title: "Ok"}))
	egui.PLastWidget.SetCallBackProc("bclick", fsett4, "fsett4")

	pDlg.Activate()
	return ""
}

func fsett4(p []string)string {
	if p == nil {}
	s := pEdi1.GetText()
	fmt.Println( s )
	egui.PLastWindow.Close()
	return ""
}

func fmbox1(p []string)string {
	if len(p) == 0 {
		egui.MsgInfo( "Test1", "MsgBox", "fmbox1", fmbox1, "mm1" )
	} else if p[0] == "mm1" {
		egui.MsgInfo( "Test2", "MsgBox", "", nil, "" )
	}
	return ""
}

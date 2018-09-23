package main

import (
	egui "github.com/alkresin/external"
)

const (
	CLR_LBLUE  = 16759929
	CLR_LBLUE0 = 12164479
	CLR_LBLUE2 = 16770002
	CLR_LBLUE3 = 16772062
	CLR_LBLUE4 = 16775920
	ES_MULTILINE = 4
)

func main() {

	if !egui.Init("port=3105\nlog") {
		return
	}

	egui.CreateStyle( &(egui.Style{Name: "st1", Orient: 1, Colors: []int32{CLR_LBLUE,CLR_LBLUE3}}) )
	egui.CreateStyle( &(egui.Style{Name: "st2", Colors: []int32{CLR_LBLUE}, BorderW: 3}) )
	egui.CreateStyle( &(egui.Style{Name: "st3", Colors: []int32{CLR_LBLUE},
		BorderW: 2, BorderClr: CLR_LBLUE0}) )
	egui.CreateStyle( &(egui.Style{Name: "st4", Colors: []int32{CLR_LBLUE2,CLR_LBLUE3},
		BorderW: 1, BorderClr: CLR_LBLUE}) )

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 400, H: 280, Title: "External"})
	egui.InitMainWindow(pWindow)

	pPanel := pWindow.AddWidget(&(egui.Widget{Type: "paneltop", H: 40,
		AProps: map[string]string{"HStyle":"st1"} }))

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 0, Y: 0, W: 56, H: 40, Title: "Date",
		AProps: map[string]string{"HStyles": egui.ArrStrings("st1","st2","st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_WriteStatus(HWindow():GetMain(),1,Dtoc(Date()),.T.)")

	pPanel = pWindow.AddWidget(&(egui.Widget{Type: "panel", X: 0, Y: 40, W: 200, H: 208 }))
	egui.PLastWidget.SetCallBackProc("onsize", nil, "{|o,x,y|o:Move(,,,y-72)}")

	pEdi := pWindow.AddWidget(&(egui.Widget{Type: "edit", Name: "edim", X: 204, Y: 40, W: 196, H: 180,
		Winstyle: ES_MULTILINE }))
	egui.PLastWidget.SetCallBackProc("onsize", nil, "{|o,x,y|o:Move(,,x-o:nLeft,y-72)}")

	pWindow.AddWidget(&(egui.Widget{Type: "splitter", X: 200, Y: 40, W: 4, H: 208,
		Anchor: egui.A_TOPABS+egui.A_BOTTOMABS,
		AProps: map[string]string{"ALeft": egui.ArrWidgs(pPanel), "ARight": egui.ArrWidgs(pEdi)} }))
	egui.PLastWidget.SetCallBackProc("onsize", nil, "{|o,x,y|o:Move(,,,y-72)}")

	pWindow.AddWidget(&(egui.Widget{Type: "panelbot", H: 32,
		AProps: map[string]string{"HStyle":"st4","AParts": egui.ArrInts(120,120,0)} }))

	pWindow.Activate()

	egui.Exit()

}


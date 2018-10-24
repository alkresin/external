// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	egui "github.com/alkresin/external"
	"strconv"
	"time"
)

const (
	CLR_LBLUE  = 16759929
	CLR_LBLUE0 = 12164479
	CLR_LBLUE2 = 16770002
	CLR_LBLUE3 = 16772062
	CLR_LBLUE4 = 16775920
)

func main() {

	if !egui.Init("port=3105\nlog") {
		return
	}

	egui.CreateStyle(&(egui.Style{Name: "st1", Orient: 1, Colors: []int32{CLR_LBLUE, CLR_LBLUE3}}))
	egui.CreateStyle(&(egui.Style{Name: "st2", Colors: []int32{CLR_LBLUE}, BorderW: 3}))
	egui.CreateStyle(&(egui.Style{Name: "st3", Colors: []int32{CLR_LBLUE},
		BorderW: 2, BorderClr: CLR_LBLUE0}))
	egui.CreateStyle(&(egui.Style{Name: "st4", Colors: []int32{CLR_LBLUE2, CLR_LBLUE3},
		BorderW: 1, BorderClr: CLR_LBLUE}))

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 400, H: 280, Title: "External"})
	egui.InitMainWindow(pWindow)

	egui.Menu("")
	{
		egui.Menu("File")
		{
			egui.AddMenuItem("Set text to label",
				func(p []string) string { egui.Widg("main.l1").SetText(p[0]); return "" }, "fsett2", "Bye...1")
			egui.AddMenuSeparator()
			egui.AddMenuItem("Printing", fprint, "fprint")
			egui.AddMenuSeparator()
			egui.AddMenuItem("Exit", nil, "hwg_EndWindow()")
		}
		egui.EndMenu()

		egui.Menu("Dialogs")
		{
			egui.AddMenuItem("Open dialog", fsett3, "fsett3")
			egui.AddMenuItem("Test Tab", ftab, "ftab")
			egui.AddMenuItem("Test browse", fbrowse, "fbrowse")
		}
		egui.EndMenu()

		egui.Menu("Standard dialogs")
		{
			egui.AddMenuItem("Message boxes", fmbox1, "fmbox1")
			egui.AddMenuItem("MsgGet box", fmbox2, "fmbox2")
			egui.AddMenuItem("Choice", fmbox3, "fmbox3")
			egui.AddMenuItem("Select color", fsele_color, "fsele_color")
			egui.AddMenuItem("Select font", fsele_font, "fsele_font")
			egui.AddMenuItem("Select file", fsele_file, "fsele_file")
		}
		egui.EndMenu()

		egui.Menu("Help")
		{
			egui.AddMenuItem("About", nil, "hwg_MsgInfo(hb_version()+chr(10)+chr(13)+hwg_version(),\"About\")")
		}
		egui.EndMenu()
	}
	egui.EndMenu()

	pPanel := pWindow.AddWidget(&(egui.Widget{Type: "paneltop", H: 40,
		AProps: map[string]string{"HStyle": "st1"}}))

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 0, Y: 0, W: 56, H: 40, Title: "Date",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_WriteStatus(HWindow():GetMain(),1,Dtoc(Date()),.T.)")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 56, Y: 0, W: 56, H: 40, Title: "Time",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_WriteStatus(HWindow():GetMain(),2,Time(),.T.)")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 112, Y: 0, W: 56, H: 40, Title: "Get",
		AProps: map[string]string{"HStyles": egui.ToString("st1", "st2", "st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett3, "fsett3")

	pWindow.AddWidget(&(egui.Widget{Type: "label", Name: "l1",
		X: 20, Y: 60, W: 180, H: 24, Title: "Test of a label",
		AProps: map[string]string{"Transpa": "t"}}))

	pWindow.AddWidget(&(egui.Widget{Type: "button", X: 200, Y: 56, W: 100, H: 32, Title: "SetText"}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett1, "fsett1", "first parameter")

	pWindow.AddWidget(&(egui.Widget{Type: "panelbot", H: 32,
		AProps: map[string]string{"HStyle": "st4", "AParts": egui.ToString(120, 120, 0)}}))

	pWindow.Activate()

	egui.Exit()

}

func fsett1(p []string) string {

	pLabel := egui.Widg("main.l1")
	fmt.Println(pLabel.GetText())
	pLabel.SetText(p[1])

	return ""
}

func fsett3([]string) string {

	egui.BeginPacket()
	egui.SetDateFormat("DD.MM.YYYY")
	pFont := egui.CreateFont(&(egui.Font{Name: "f1", Family: "Georgia", Height: 16}))
	pDlg := &(egui.Widget{Name: "dlg", X: 300, Y: 200, W: 400, H: 260, Title: "Dialog Test", Font: pFont})
	egui.InitDialog(pDlg)

	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 10, W: 160, H: 24, Title: "Identifier:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi1", X: 20, Y: 32, W: 160, H: 24,
		AProps: map[string]string{"Picture": "@!R /XXX:XXX/"}}))
	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 220, Y: 10, W: 160, H: 24, Title: "Date:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi2", X: 220, Y: 32, W: 120, H: 24,
		Title: time.Now().Format("20060102"), AProps: map[string]string{"Picture": "D@D"}}))

	pDlg.AddWidget(&(egui.Widget{Type: "combo", Name: "comb", X: 20, Y: 68, W: 160, H: 24,
		AProps: map[string]string{"AItems": egui.ToString("first", "second", "third")}}))

	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 220, Y: 68, W: 80, H: 24, Title: "Age:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "updown", Name: "upd1", X: 280, Y: 68, W: 60, H: 24}))

	pDlg.AddWidget(&(egui.Widget{Type: "group", X: 10, Y: 110, W: 180, H: 76, Title: "Check"}))
	pDlg.AddWidget(&(egui.Widget{Type: "check", Name: "chk1", X: 24, Y: 124, W: 150, H: 24, Title: "Married"}))
	pDlg.AddWidget(&(egui.Widget{Type: "check", Name: "chk2", X: 24, Y: 148, W: 150, H: 24, Title: "Has children"}))

	pGroup := pDlg.AddWidget(&(egui.Widget{Type: "radiogr", Name: "rg", X: 200, Y: 110, W: 180, H: 76, Title: "Radio"}))
	pDlg.AddWidget(&(egui.Widget{Type: "radio", X: 224, Y: 124, W: 150, H: 24, Title: "Male"}))
	pDlg.AddWidget(&(egui.Widget{Type: "radio", X: 224, Y: 148, W: 150, H: 24, Title: "Female"}))
	egui.RadioEnd(pGroup, 1)

	pDlg.AddWidget(&(egui.Widget{Type: "button", X: 150, Y: 220, W: 100, H: 32, Title: "Ok"}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett4, "fsett4")

	pDlg.Activate()
	egui.EndPacket()

	return ""
}

func fsett4([]string) string {
	arr := egui.GetValues(egui.Wnd("dlg"), []string{"edi1", "edi2", "comb", "chk1", "chk2", "rg", "upd1"})
	egui.MsgInfo("Id: "+arr[0]+"\r\n"+"Date: "+arr[1]+"\r\n"+"Combo: "+arr[2]+"\r\n"+
		"Married: "+arr[3]+"\r\n"+"Has children: "+arr[4]+"\r\n"+"Sex: "+arr[5]+"\r\n"+
		"Age: "+arr[6], "Result", "", nil, "")
	egui.PLastWindow.Close()
	return ""
}

func ftab([]string) string {

	egui.BeginPacket()
	pFont := egui.CreateFont(&(egui.Font{Name: "f1", Family: "Georgia", Height: 16}))
	pDlg := &(egui.Widget{Name: "dlg2", X: 300, Y: 200, W: 200, H: 340, Title: "Tab", Font: pFont})
	egui.InitDialog(pDlg)

	pTab := pDlg.AddWidget(&(egui.Widget{Type: "tab", X: 10, Y: 10, W: 180, H: 280}))

	egui.TabPage(pTab, "First")
	pTab.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 30, W: 140, H: 24, Title: "Name:"}))
	pTab.AddWidget(&(egui.Widget{Type: "edit", X: 20, Y: 52, W: 140, H: 24}))
	pTab.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 84, W: 140, H: 24, Title: "Surname:"}))
	pTab.AddWidget(&(egui.Widget{Type: "edit", X: 20, Y: 106, W: 140, H: 24}))
	egui.TabPageEnd(pTab)

	egui.TabPage(pTab, "Second")
	pTab.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 40, W: 140, H: 24, Title: "Age:"}))
	pTab.AddWidget(&(egui.Widget{Type: "edit", X: 20, Y: 62, W: 140, H: 24}))
	pTab.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 94, W: 140, H: 24, Title: "Profession:"}))
	pTab.AddWidget(&(egui.Widget{Type: "edit", X: 20, Y: 116, W: 140, H: 24}))
	egui.TabPageEnd(pTab)

	pDlg.AddWidget(&(egui.Widget{Type: "button", X: 60, Y: 300, W: 100, H: 32, Title: "Ok"}))
	egui.PLastWidget.SetCallBackProc("onclick", ftabclose, "ftabclose")

	pDlg.Activate()
	egui.EndPacket()

	return ""
}

func fbrowse([]string) string {

	var arr = [][]string{{"Alex", "17", "1200"}, {"Victor", "42", "1600"}, {"John", "31", "1000"}}

	egui.BeginPacket()
	pFont := egui.CreateFont(&(egui.Font{Name: "f1", Family: "Georgia", Height: 16}))
	pDlg := &(egui.Widget{Name: "dlg2", X: 300, Y: 200, W: 280, H: 250, Title: "browse", Font: pFont})
	egui.InitDialog(pDlg)

	pBrw := pDlg.AddWidget(&(egui.Widget{Type: "browse", X: 10, Y: 10, W: 260, H: 180}))
	pBrw.SetParam("oStyleHead", egui.GetStyle("st1"))
	egui.BrwSetArray(pBrw, arr)
	egui.BrwSetColumn(pBrw, 1, "Name", 1, 0, false)
	egui.BrwSetColumn(pBrw, 2, "Age", 1, 0, false)
	egui.BrwSetColumn(pBrw, 3, "Salary", 1, 0, false)

	pDlg.AddWidget(&(egui.Widget{Type: "button", X: 90, Y: 210, W: 100, H: 32, Title: "Ok"}))
	egui.PLastWidget.SetCallBackProc("onclick", ftabclose, "ftabclose")

	pDlg.Activate()
	egui.EndPacket()

	return ""
}

func ftabclose([]string) string {
	egui.PLastWindow.Close()
	return ""
}

func fmbox1(p []string) string {
	if len(p) == 0 {
		egui.MsgYesNo("Yes or No???", "MsgBox", "fmbox1", fmbox1, "mm1")
	} else if p[0] == "mm1" {
		if p[1] == "t" {
			egui.MsgInfo("Yes!", "Answer", "", nil, "")
		} else {
			egui.MsgInfo("No...", "Answer", "", nil, "")
		}
	}
	return ""
}

func fmbox2(p []string) string {
	if len(p) == 0 {
		egui.MsgGet("Input something:", "MsgGet", 0, "fmbox2", fmbox2, "mm1")
	} else if p[0] == "mm1" {
		egui.MsgInfo(p[1], "Answer", "", nil, "")
	}
	return ""
}

func fmbox3(p []string) string {
	if len(p) == 0 {
		arr := []string{"Alex Petrov", "Serg Lama", "Jimmy Hendrix", "Dorian Gray", "Victor Peti"}
		egui.Choice(arr, "Select from a list", "fmbox3", fmbox3, "mm1")
	} else if p[0] == "mm1" {
		egui.MsgInfo(p[1], "Answer", "", nil, "")
	}
	return ""
}

func fsele_color(p []string) string {
	if len(p) == 0 {
		egui.SelectColor(0, "fsele_color", fsele_color, "mm1")
	} else {
		iColor, _ := strconv.Atoi(p[1])
		egui.Widg("main.l1").SetColor(int32(iColor), -1)
	}
	return ""
}

func fsele_font(p []string) string {
	if len(p) == 0 {
		egui.SelectFont("fsele_font", fsele_font, "")
	} else {
		fmt.Println("font id: ", p[0])
		if pFont := egui.GetFont(p[0]); pFont != nil {
			if len(p) < 8 {
			} else {
				fmt.Println("font fam: ", p[1])
			}
		}
	}
	return ""
}

func fsele_file(p []string) string {
	if len(p) == 0 {
		egui.SelectFile("", "fsele_file", fsele_file, "mm1")
	} else {
		if p[1] == "" {
			egui.MsgInfo("Nothing selected", "Result", "", nil, "")
		} else {
			egui.MsgInfo(p[1], "File selected", "", nil, "")
		}
	}
	return ""
}

func fprint(p []string) string {
	if len(p) == 0 {
		egui.InitPrinter(&(egui.Printer{SPrinter: "...", BPreview: true}), "fprint", fprint, "mm1")
	} else {
		pPrinter := egui.PLastPrinter
		pFont := pPrinter.AddFont(&(egui.Font{Family: "Times New Roman", Height: 10}))
		pPrinter.StartPage()
		pPrinter.SetFont(pFont)
		pPrinter.Box(5, 5, 200, 282)
		pPrinter.Say(50, 10, 165, 26, "Printing first sample !", egui.DT_CENTER)
		pPrinter.Line(45, 30, 170, 30)
		pPrinter.Line(45, 5, 45, 30)
		pPrinter.Line(170, 5, 170, 30)
		pPrinter.Say(50, 120, 150, 132, "----------", egui.DT_CENTER)
		pPrinter.Box(50, 134, 160, 146)
		pPrinter.Say(50, 135, 160, 146, "End Of Report", egui.DT_CENTER)
		pPrinter.EndPage()
		pPrinter.End()

	}
	return ""
}

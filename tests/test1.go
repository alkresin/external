package main

import (
	"fmt"
	"strconv"
	egui "github.com/alkresin/external"
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

	egui.CreateStyle( &(egui.Style{Name: "st1", Orient: 1, Colors: []int32{CLR_LBLUE,CLR_LBLUE3}}) )
	egui.CreateStyle( &(egui.Style{Name: "st2", Colors: []int32{CLR_LBLUE}, BorderW: 3}) )
	egui.CreateStyle( &(egui.Style{Name: "st3", Colors: []int32{CLR_LBLUE},
		BorderW: 2, BorderClr: CLR_LBLUE0}) )
	egui.CreateStyle( &(egui.Style{Name: "st4", Colors: []int32{CLR_LBLUE2,CLR_LBLUE3},
		BorderW: 1, BorderClr: CLR_LBLUE}) )

	pWindow := &(egui.Widget{X: 100, Y: 100, W: 400, H: 280, Title: "External"})
	egui.InitMainWindow(pWindow)

	egui.Menu("")
	{
		egui.Menu( "File" )
		{
			egui.AddMenuItem( "Set text",
				func (p []string)string { egui.GetWidg("main.l1").SetText(p[0]); return "" }, "fsett2", "Bye...1" )
			egui.AddMenuItem( "Open dialog", fsett3, "fsett3" )
			egui.AddMenuSeparator()
			egui.AddMenuItem( "Message boxes", fmbox1, "fmbox1" )
			egui.AddMenuItem( "MsgGet box", fmbox2, "fmbox2" )
			egui.AddMenuItem( "Choice", fmbox3, "fmbox3" )
			egui.AddMenuItem( "Select color", fsele_color, "fsele_color" )
			egui.AddMenuItem( "Select font", fsele_font, "fsele_font" )
			egui.AddMenuItem( "Select file", fsele_file, "fsele_file" )
			egui.AddMenuSeparator()
			egui.AddMenuItem( "Exit", nil, "hwg_EndWindow()" )
		}
		egui.EndMenu()
		egui.Menu( "Help" )
		{
			egui.AddMenuItem( "About", nil, "hwg_MsgInfo(hb_version()+chr(10)+chr(13)+hwg_version(),\"About\")" )
		}
		egui.EndMenu()
	}
	egui.EndMenu()

	pPanel := pWindow.AddWidget(&(egui.Widget{Type: "paneltop", H: 40,
		AProps: map[string]string{"HStyle":"st1"} }))

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 0, Y: 0, W: 56, H: 40, Title: "Date",
		AProps: map[string]string{"HStyles": egui.ArrStrings("st1","st2","st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_WriteStatus(HWindow():GetMain(),1,Dtoc(Date()),.T.)")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 56, Y: 0, W: 56, H: 40, Title: "Time",
		AProps: map[string]string{"HStyles": egui.ArrStrings("st1","st2","st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", nil, "hwg_WriteStatus(HWindow():GetMain(),2,Time(),.T.)")

	pPanel.AddWidget(&(egui.Widget{Type: "ownbtn", X: 112, Y: 0, W: 56, H: 40, Title: "Get",
		AProps: map[string]string{"HStyles": egui.ArrStrings("st1","st2","st3")}}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett3, "fsett3")

	pWindow.AddWidget(&(egui.Widget{Type: "label", Name: "l1",
		X: 20, Y: 60, W: 180, H: 24, Title: "Test of a label",
		AProps: map[string]string{"Transpa":"t"} }))

	pWindow.AddWidget(&(egui.Widget{Type: "button", X: 200, Y: 56, W: 100, H: 32, Title: "SetText"}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett1, "fsett1", "first parameter")

	pWindow.AddWidget(&(egui.Widget{Type: "panelbot", H: 32,
		AProps: map[string]string{"HStyle":"st4","AParts": egui.ArrInts(120,120,0)} }))

	pWindow.Activate()

	egui.Exit()

}

func fsett1(p []string)string {

	pLabel := egui.GetWidg("main.l1")
	fmt.Println( pLabel.GetText() )
	pLabel.SetText( p[1] )

	return ""
}


func fsett3(p []string)string {
	if p == nil {}

	egui.StartPacket()
	pFont := egui.CreateFont( &(egui.Font{Name: "f1", Family: "Georgia", Height: 16}) )
	pDlg := &(egui.Widget{Name: "dlg", X: 300, Y: 200, W: 200, H: 440, Title: "Dialog Test", Font: pFont })
	egui.InitDialog(pDlg)

	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 10, W: 180, H: 24, Title: "Name:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi1", X: 20, Y: 32, W: 160, H: 24 }))
	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 64, W: 180, H: 24, Title: "Surname:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "edit", Name: "edi2", X: 20, Y: 86, W: 160, H: 24 }))

	pDlg.AddWidget(&(egui.Widget{Type: "combo", Name: "comb", X: 20, Y: 116, W: 160, H: 24,
	      AProps: map[string]string{"AItems": egui.ArrStrings("first","second","third")} }))

	pDlg.AddWidget(&(egui.Widget{Type: "group", X: 10, Y: 150, W: 180, H: 76, Title: "Check"}))
	pDlg.AddWidget(&(egui.Widget{Type: "check", Name: "chk1", X: 24, Y: 174, W: 150, H: 24, Title: "Married"}))
	pDlg.AddWidget(&(egui.Widget{Type: "check", Name: "chk2", X: 24, Y: 198, W: 150, H: 24, Title: "Has children"}))

	pDlg.AddWidget(&(egui.Widget{Type: "radiogr", Name: "rg", X: 10, Y: 240, W: 180, H: 76, Title: "Radio"}))
	pDlg.AddWidget(&(egui.Widget{Type: "radio", X: 24, Y: 264, W: 150, H: 24, Title: "Male"}))
	pDlg.AddWidget(&(egui.Widget{Type: "radio", X: 24, Y: 288, W: 150, H: 24, Title: "Female"}))

	pDlg.AddWidget(&(egui.Widget{Type: "label", X: 20, Y: 332, W: 80, H: 24, Title: "Age:"}))
	pDlg.AddWidget(&(egui.Widget{Type: "updown", Name: "upd1", X: 100, Y: 330, W: 60, H: 24 }))

	pDlg.AddWidget(&(egui.Widget{Type: "button", X: 50, Y: 390, W: 100, H: 32, Title: "Ok"}))
	egui.PLastWidget.SetCallBackProc("onclick", fsett4, "fsett4")

	pDlg.Activate()
	egui.EndPacket()

	return ""
}

func fsett4(p []string)string {
	if p == nil {}
	arr := egui.GetValues( egui.GetWnd("dlg"), []string{ "edi1","edi2","comb","chk1","chk2","rg","upd1" } )
	egui.MsgInfo( "Name: "+arr[0]+"\r\n" + "Surname: "+arr[1]+"\r\n" + "Combo: "+arr[2]+"\r\n" +
		"Married: "+arr[3]+"\r\n" + "Has children: "+arr[4]+"\r\n" + "Sex: "+arr[5]+"\r\n" +
		"Age: "+arr[6], "Result", "", nil, "" )
	egui.PLastWindow.Close()
	return ""
}

func fmbox1(p []string)string {
	if len(p) == 0 {
		egui.MsgYesNo( "Yes or No???", "MsgBox", "fmbox1", fmbox1, "mm1" )
	} else if p[0] == "mm1" {
		if p[1] == "t" {
			egui.MsgInfo( "Yes!", "Answer", "", nil, "" )
		} else {
			egui.MsgInfo( "No...", "Answer", "", nil, "" )
		}
	}
	return ""
}

func fmbox2(p []string)string {
	if len(p) == 0 {
		egui.MsgGet( "Input something:", "MsgGet", 0, "fmbox2", fmbox2, "mm1" )
	} else if p[0] == "mm1" {
		egui.MsgInfo( p[1], "Answer", "", nil, "" )
	}
	return ""
}

func fmbox3(p []string)string {
	if len(p) == 0 {
		arr := []string{ "Alex Petrov", "Serg Lama", "Jimmy Hendrix", "Dorian Gray", "Victor Peti" }
		egui.Choice( arr, "Select from a list", "fmbox3", fmbox3, "mm1" )
	} else if p[0] == "mm1" {
		egui.MsgInfo( p[1], "Answer", "", nil, "" )
	}
	return ""
}

func fsele_color(p []string)string {
	if len(p) == 0 {
		egui.SelectColor( 0, "fsele_color", fsele_color, "mm1" )
	} else {
		iColor,_ := strconv.Atoi(p[1])
		egui.GetWidg("main.l1").SetColor( int32(iColor),-1 )
	}
	return ""
}

func fsele_font(p []string)string {
	if len(p) == 0 {
		egui.SelectFont( "fsele_font", fsele_font, "" )
	} else {
		fmt.Println( "font id: ", p[0] )
		if pFont := egui.GetFont( p[0] ); pFont != nil {
			if len(p) < 8 {
			} else {
				fmt.Println( "font fam: ", p[1] )
			}
		}
	}
	return ""
}

func fsele_file(p []string)string {
	if len(p) == 0 {
		egui.SelectFile( "", "fsele_file", fsele_file, "mm1" )
	} else {
		if p[1] == "" {
			egui.MsgInfo( "Nothing selected", "Result", "", nil, "" )
		} else {
			egui.MsgInfo( p[1], "File selected", "", nil, "" )
		}
	}
	return ""
}

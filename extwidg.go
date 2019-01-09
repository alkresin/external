// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package external

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// A set of constants of the anchor values
const (
	A_TOPLEFT   = -1  // Anchors control to the top and left borders of the container and does not change the distance between the top and left borders. (Default)
	A_TOPABS    = 1   // Anchors control to top border of container and does not change the distance between the top border.
	A_LEFTABS   = 2   // Anchors control to left border of container and does not change the distance between the left border.
	A_BOTTOMABS = 4   // Anchors control to bottom border of container and does not change the distance between the bottom border.
	A_RIGHTABS  = 8   // Anchors control to right border of container and does not change the distance between the right border.
	A_TOPREL    = 16  // Anchors control to top border of container and maintains relative distance between the top border.
	A_LEFTREL   = 32  // Anchors control to left border of container and maintains relative distance between the left border.
	A_BOTTOMREL = 64  // Anchors control to bottom border of container and maintains relative distance between the bottom border.
	A_RIGHTREL  = 128 // Anchors control to right border of container and maintains relative distance between the right border.
	A_HORFIX    = 256 // Anchors center of control relative to left and right borders but remains fixed in size.
	A_VERTFIX   = 512 // Anchors center of control relative to top and bottom borders but remains fixed in size.
)

// A set of constants of Winstyle values
const (
	DT_LEFT   = 0
	DT_CENTER = 1
	DT_RIGHT  = 2

	ES_PASSWORD  = 32
	ES_MULTILINE = 4
	ES_READONLY  = 2048

	WS_HSCROLL = 2097152
	WS_VSCROLL = 1048576

	WND_NOTITLE   = -1
	WND_NOSYSMENU = -2
	WND_NOSIZEBOX = -4
)

// A set of constants of the printer paper types
const (
	DMPAPER_A3 = 8  // A3 297 x 420 mm
	DMPAPER_A4 = 9  // A4 210 x 297 mm
	DMPAPER_A5 = 11 // A5 148 x 210 mm
	DMPAPER_A6 = 70 // A6 105 x 148 mm
)

// A set of constants of the code editor highliter
const (
	HILI_KEYW  = 1
	HILI_FUNC  = 2
	HILI_QUOTE = 3
	HILI_COMM  = 4
)

// The CodeBlock type for Harbour scripts, which are set via SetParam() method and BrwSetColumnEx()
type CodeBlock string

// The Font structure prepares data to create a new font
type Font struct {
	Family    string
	Name      string
	Height    int
	Bold      bool
	Italic    bool
	Underline bool
	Strikeout bool
	Charset   int16
}

// The Style structure prepares data to create a new style
type Style struct {
	Name      string
	Orient    int16
	Colors    []int32
	Corners   []int32
	BorderW   int8
	BorderClr int32
	Bitmap    string
}

// The Highlight structure serves to create a highlight rules for a code editor
type Highlight struct {
	Name string
}

// The Printer structure prepares data to initialize a printer
type Printer struct {
	Name       string
	SPrinter   string
	BPreview   bool
	IFormType  int
	BLandscape bool
}

// The Widget structure prepares data to create a new widget or window
type Widget struct {
	Parent   *Widget
	Type     string
	Name     string
	X        int
	Y        int
	W        int
	H        int
	Title    string
	Winstyle int32
	TColor   int32
	BColor   int32
	Tooltip  string
	Anchor   int32
	Font     *Font
	AProps   map[string]string
	aWidgets []*Widget
}

var mfu map[string]func([]string) string
var pMainWindow *Widget
var aDialogs []*Widget
var aFonts []*Font
var aStyles []*Style
var iIdCount int32

// PLastWindow is a pointer to a last used window structure (*Widget)
var PLastWindow *Widget

// PLastWidget is a pointer to a last used widget structure (*Widget)
var PLastWidget *Widget

// PLastPrinter is a pointer to a last used printer structure (*Printer)
var PLastPrinter *Printer

// Var mWidgs includes all possible widgets types with
// its properties, which may be installed, using AProps member of a Widget structure.
var mWidgs = map[string]map[string]string{
	"main":      {"Icon": "C"},
	"dialog":    nil,
	"label":     {"Transpa": "L"},
	"edit":      {"Picture": "C"},
	"button":    nil,
	"check":     {"Transpa": "L"},
	"radio":     {"Transpa": "L"},
	"radiogr":   nil,
	"group":     nil,
	"combo":     {"AItems": "AC"},
	"bitmap":    {"Transpa": "L", "TrColor": "N", "Image": "C"},
	"line":      {"Vertical": "L"},
	"panel":     {"HStyle": "C"},
	"paneltop":  {"HStyle": "C"},
	"panelbot":  {"HStyle": "C", "AParts": "AC"},
	"panelhead": {"HStyle": "C", "Xt": "N", "Yt": "N", "BtnClose": "L", "BtnMax": "L", "BtnMin": "L"},
	"ownbtn":    {"Transpa": "L", "TrColor": "N", "Image": "C", "HStyles": "AC", "Xt": "N", "Yt": "N"},
	"splitter":  {"Vertical": "L", "From": "N", "TO": "N", "ALeft": "AC", "ARight": "AC", "HStyle": "C"},
	"updown":    {"From": "N", "TO": "N"},
	"tree":      {"AImages": "AC", "EditLabel": "L"},
	"progress":  {"Maxpos": "N"},
	"tab":       nil,
	"browse":    {"Append": "L", "Autoedit": "L", "NoVScroll": "L", "NoBorder": "L"},
	"cedit":     {"NoVScroll": "L", "NoBorder": "L"},
	"link":      {"Link": "C", "ClrVisited": "N", "ClrLink": "N", "ClrOver": "N"},
	"monthcal":  {"NoToday": "L", "NoTodayCirc": "L", "WeekNumb": "L"}}

func widgFullName(pWidg *Widget) string {
	sName := pWidg.Name

	for pWidg.Parent != nil {
		pWidg = pWidg.Parent
		sName = pWidg.Name + "." + sName
	}
	return sName
}

// Returns a pointer to a Font structure with a Name member equal to sName argument.
func GetFont(sName string) *Font {
	if aFonts != nil {
		for _, o := range aFonts {
			if o.Name == sName {
				return o
			}
		}
	}
	return nil
}

// Returns a pointer to a Style structure with a Name member equal to sName argument.
func GetStyle(sName string) *Style {
	if aStyles != nil {
		for _, o := range aStyles {
			if o.Name == sName {
				return o
			}
		}
	}
	return nil
}

// Wnd returns a pointer to a Widget structure (a window or a dialog) with a Name member equal to sName argument.
func Wnd(sName string) *Widget {
	if sName == "main" {
		return pMainWindow
	} else if aDialogs != nil {
		for _, o := range aDialogs {
			if o.Name == sName {
				return o
			}
		}
	}
	return nil
}

// Widg returns a pointer to a Widget structure (a widget) with a Name member corresponding to sName argument.
// The sName must be compound name, containing a names of all parent widgets and windows, defined by dots.
func Widg(sName string) *Widget {
	npos := strings.Index(sName, ".")
	if npos == -1 {
		return Wnd(sName)
	}
	sWnd := sName[:npos]
	sName = sName[npos+1:]
	if oWnd := Wnd(sWnd); oWnd != nil {
		for npos = strings.Index(sName, "."); npos > -1; npos = strings.Index(sName, ".") {
			sWnd := sName[:npos]
			sName = sName[npos+1:]
			for _, o := range oWnd.aWidgets {
				if o.Name == sWnd {
					oWnd = o
					break
				}
			}
			if oWnd == nil {
				return nil
			}
		}
		for _, o := range oWnd.aWidgets {
			if o.Name == sName {
				return o
			}
		}
	}
	return nil
}

func setprops(pWidg *Widget, mwidg map[string]string) string {

	sPar := ""
	if pWidg.Winstyle != 0 {
		sPar += fmt.Sprintf(",\"Winstyle\": %d", pWidg.Winstyle)
	}
	if pWidg.TColor != 0 {
		sPar += fmt.Sprintf(",\"TColor\": %d", pWidg.TColor)
	}
	if pWidg.BColor != 0 {
		sPar += fmt.Sprintf(",\"BColor\": %d", pWidg.BColor)
	}
	if pWidg.Tooltip != "" {
		sPar += fmt.Sprintf(",\"Tooltip\": \"%s\"", pWidg.Tooltip)
	}
	if pWidg.Font != nil {
		sPar += fmt.Sprintf(",\"Font\": \"%s\"", pWidg.Font.Name)
	}
	if pWidg.Anchor != 0 {
		if pWidg.Anchor == A_TOPLEFT {
			pWidg.Anchor = 0
		}
		sPar += fmt.Sprintf(",\"Anchor\": %d", pWidg.Anchor)
	}
	if pWidg.AProps != nil {
		for name, val := range pWidg.AProps {
			cType, bOk := mwidg[name]
			if bOk {
				if cType == "C" {
					sPar += fmt.Sprintf(",\"%s\": \"%s\"", name, val)
				} else if cType == "L" {
					sPar += fmt.Sprintf(",\"%s\": \"%s\"", name, val)
				} else if cType == "N" {
					sPar += fmt.Sprintf(",\"%s\": %s", name, val)
				} else if cType == "AC" {
					sPar += fmt.Sprintf(",\"%s\": %s", name, val)
				}
			} else {
				WriteLog(fmt.Sprintf("Error! \"%s\" does not defined for \"%s\"\r\n", name, pWidg.Type))
			}
		}
	}
	if sPar != "" {
		sPar = ",{" + sPar[1:] + "}"
	}
	return sPar
}

// ToString converts function arguments to a json string
func ToString(xParam ...interface{}) string {

	for i, x := range xParam {
		switch v := x.(type) {
		case *Font:
			xParam[i] = v.Name
		case *Style:
			xParam[i] = v.Name
		case *Widget:
			xParam[i] = v.Name
		}
	}

	b, _ := json.Marshal(xParam)
	return string(b)
}

// OpenMainForm reads a main window description from a xml file, prepared by HwGUI's Designer,
// initialises and activates this window with all its widgets
func OpenMainForm(sForm string) bool {
	var bres bool
	b, err := json.Marshal(sForm)
	if err != nil {
		WriteLog( fmt.Sprintln(err) )
		return false
	}
	bres = sendout("[\"openformmain\"," + string(b) + "]")
	wait()
	return bres
}

// OpenForm reads a dialog window description from a xml file, prepared by HwGUI's Designer,
// initialises and activates this dialog with all its widgets
func OpenForm(sForm string) bool {
	var bres bool
	b, err := json.Marshal(sForm)
	if err != nil {
		WriteLog( fmt.Sprintln(err) )
		return false
	}
	bres = sendout("[\"openform\"," + string(b) + "]")
	return bres
}

// OpenReport reads a report description from a xml file, prepared by HwGUI's Designer
// and prints this report
func OpenReport(sForm string) bool {
	var b bool
	b = sendout("[\"openreport\",\"" + sForm + "\"]")
	return b
}

// CreateFont creates a font with parameters, defined in a structure, pointed by pFont argument.
func CreateFont(pFont *Font) *Font {

	pFont.new()
	sParams := fmt.Sprintf("[\"crfont\",\"%s\",\"%s\",%d,%t,%t,%t,%t,%d]", pFont.Name, pFont.Family, pFont.Height,
		pFont.Bold, pFont.Italic, pFont.Underline, pFont.Strikeout, pFont.Charset)
	sendout(sParams)
	return pFont
}

// CreateStyle creates a style with parameters, defined in a structure, pointed by pStyle argument.
func CreateStyle(pStyle *Style) *Style {

	if pStyle.Name == "" {
		pStyle.Name = fmt.Sprintf("s%d", iIdCount)
		iIdCount++
	}
	if aStyles == nil {
		aStyles = make([]*Style, 0, 16)
	}
	aStyles = append(aStyles, pStyle)
	b1, _ := json.Marshal(pStyle.Colors)
	b2, _ := json.Marshal(pStyle.Corners)
	sParams := fmt.Sprintf("[\"crstyle\",\"%s\",%s,%d,%s,%d,%d,\"%s\"]", pStyle.Name,
		string(b1), pStyle.Orient, string(b2),
		pStyle.BorderW, pStyle.BorderClr, pStyle.Bitmap)
	sendout(sParams)
	return pStyle
}

// CreateHighliter creates a highlight rules for a code editor ("cedit" widget)
func CreateHighliter(sName string, sCommands string, sFuncs string,
	sSingleLineComm string, sMultiLineComm string, bCase bool) *Highlight {

	sParams := fmt.Sprintf("[\"highl\",\"%s\",\"%s\",\"%s\",\"%s\",\"%s\",%t]", sName,
		sCommands, sFuncs, sSingleLineComm, sMultiLineComm, bCase)
	sendout(sParams)
	return &(Highlight{Name: sName})
}

// SetHighliter sets or unsets (if p == nil) a given Highliter to a "cedit" widget.
func SetHighliter(pEdit *Widget, p *Highlight) {
	var sHiliName string
	if p == nil {
		sHiliName = ""
	} else {
		sHiliName = p.Name
	}
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"hili\",\"%s\"]",
		widgFullName(pEdit), sHiliName)
	sendout(sParams)
}

// SetHili defines highlighting options for a code editor ("cedit" widget): a font, text color and background color
func SetHiliOpt(pEdit *Widget, iGroup int, pFont *Font, tColor int32, bColor int32) {
	var sFontName string
	if pFont == nil {
		sFontName = ""
	} else {
		sFontName = pFont.Name
	}
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"hiliopt\",[%d,\"%s\",%d,%d]]",
		widgFullName(pEdit), iGroup, sFontName, tColor, bColor)
	sendout(sParams)
}

// InitPrinter initializes a printer, the name of a printer is passed in SPrinter member of
// a pPrinter structure. If it is an empty string, the default printer will be used, if it is
// defined as "...", printer setup dialog will be opened.
func InitPrinter(pPrinter *Printer, sFunc string, fu func([]string) string, sMark string) *Printer {

	if pPrinter.Name == "" {
		pPrinter.Name = fmt.Sprintf("p%d", iIdCount)
		iIdCount++
	}
	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sMark = ""
	}
	sParams := fmt.Sprintf("[\"prninit\",\"%s\",[\"%s\",%t,%d,%t],\"%s\",\"%s\"]", pPrinter.Name,
		pPrinter.SPrinter, pPrinter.BPreview, pPrinter.IFormType, pPrinter.BLandscape, sFunc, sMark)
	sendout(sParams)
	PLastPrinter = pPrinter
	return pPrinter
}

// AddFont method adds a font, described in Font structure, to the printer.
func (p *Printer) AddFont(pFont *Font) *Font {
	pFont.new()
	sParams := fmt.Sprintf("[\"print\",\"fontadd\",\"%s\",[\"%s\",\"%s\",%d,%t,%t,%t,%d]]", p.Name,
		pFont.Name, pFont.Family, pFont.Height,
		pFont.Bold, pFont.Italic, pFont.Underline, pFont.Charset)
	sendout(sParams)
	return pFont
}

// SetFont method sets a font, previously added with AddFont, as current while printing
func (p *Printer) SetFont(pFont *Font) {
	sParams := fmt.Sprintf("[\"print\",\"fontset\",\"%s\",[\"%s\"]]", p.Name, pFont.Name)
	sendout(sParams)
}

// Say method prints s text string sText in a rectangle with iTop, iLeft, iRight, iBottom
// coordinates, iOpt defines an alignment.
func (p *Printer) Say(iTop, iLeft, iRight, iBottom int32, sText string, iOpt int32) {

	sParams := fmt.Sprintf("[\"print\",\"text\",\"%s\",[\"%s\",%d,%d,%d,%d,%d]]",
		p.Name, sText, iTop, iLeft, iRight, iBottom, iOpt)
	sendout(sParams)
}

// Line methods prints a line from iTop, iLeft to iRight, iBottom
func (p *Printer) Line(iTop, iLeft, iRight, iBottom int32) {

	sParams := fmt.Sprintf("[\"print\",\"line\",\"%s\",[%d,%d,%d,%d]]", p.Name, iTop, iLeft, iRight, iBottom)
	sendout(sParams)
}

// Box method prints a rectangle with iTop, iLeft, iRight, iBottom coordinates
func (p *Printer) Box(iTop, iLeft, iRight, iBottom int32) {

	sParams := fmt.Sprintf("[\"print\",\"box\",\"%s\",[%d,%d,%d,%d]]", p.Name, iTop, iLeft, iRight, iBottom)
	sendout(sParams)
}

// StartPage method begins a new page printed
func (p *Printer) StartPage() {

	sParams := fmt.Sprintf("[\"print\",\"startpage\",\"%s\",[]]", p.Name)
	sendout(sParams)
}

// StartPage method ends a page printed
func (p *Printer) EndPage() {

	sParams := fmt.Sprintf("[\"print\",\"endpage\",\"%s\",[]]", p.Name)
	sendout(sParams)
}

// End method closes a printer
func (p *Printer) End() {

	sParams := fmt.Sprintf("[\"print\",\"end\",\"%s\",[]]", p.Name)
	sendout(sParams)
}

// Initialises a main window with parameters, defined in a structure, pointed by pWnd argument.
// To show this window on a screen it is necessary to use Activate() method.
func InitMainWindow(pWnd *Widget) bool {
	pMainWindow = pWnd
	PLastWindow = pWnd
	pWnd.Type = "main"
	pWnd.Name = "main"
	sPar2 := setprops(pWnd, mWidgs["main"])
	sParams := fmt.Sprintf("[\"crmainwnd\",[%d,%d,%d,%d,\"%s\"]%s]", pWnd.X, pWnd.Y, pWnd.W,
		pWnd.H, pWnd.Title, sPar2)
	return sendout(sParams)
}

// Initialises a dialog window with parameters, defined in a structure, pointed by pWnd argument.
// To show this window on a screen it is necessary to use Activate() method.
func InitDialog(pWnd *Widget) bool {
	PLastWindow = pWnd
	pWnd.Type = "dialog"
	if pWnd.Name == "" {
		pWnd.Name = fmt.Sprintf("w%d", iIdCount)
		iIdCount++
	}
	if aDialogs == nil {
		aDialogs = make([]*Widget, 0, 8)
	}
	aDialogs = append(aDialogs, pWnd)

	sPar2 := setprops(pWnd, mWidgs["dialog"])
	sParams := fmt.Sprintf("[\"crdialog\",\"%s\",[%d,%d,%d,%d,\"%s\"]%s]", pWnd.Name, pWnd.X, pWnd.Y, pWnd.W,
		pWnd.H, pWnd.Title, sPar2)
	return sendout(sParams)
}

// EvalProc sends a code fragment, written on Harbour to a GuiServer to execute
// and does not return a result.
func EvalProc(s string) {

	b, _ := json.Marshal(s)
	sendout("[\"evalcode\"," + string(b) + "]")
}

// EvalFunc sends a code fragment, written on Harbour to a GuiServer to execute
// and returns a result.
func EvalFunc(s string) []byte {

	b, _ := json.Marshal(s)
	b = sendoutAndReturn("[\"evalcode\"," + string(b) + ",\"t\"]")
	if b[0] == byte('+') && b[1] == byte('"') {
		b = b[2 : len(b)-1]
	}
	return b
}

// GetValues returns list of values from widgets of a pWnd window (main or a dialog),
// listed by names in aNames slice.
func GetValues(pWnd *Widget, aNames []string) []string {
	sParams := "[\"getvalues\",\"" + pWnd.Name + "\",["
	for i, v := range aNames {
		if i > 0 {
			sParams += ","
		}
		sParams += "\"" + v + "\""
	}
	sParams += "]]"
	b := sendoutAndReturn(sParams)
	arr := make([]string, len(aNames))
	err := json.Unmarshal(b[1:], &arr)
	if err != nil {
		return nil
	} else {
		return arr
	}
}

// GetVersion returns the version string in different verbosity level, depending of i value
// i == 0  - GuiServer version only ("1.3", for example);
// i == 1  - GuiServer version with "GuiServer" word;
// i == 2  - GuiServer, Harbour and HwGUI versions.
func GetVersion(i int) string {

	var sRes string
	b := sendoutAndReturn("[\"getver\"," + strconv.Itoa(i) + "]")
	if b[0] == byte('+') {
		b = b[1:len(b)]
	}
	err := json.Unmarshal(b, &sRes)
	if err != nil {
		return ""
	}
	return sRes
}

// MsgInfo creates a standard nessagebox
// sTitle - box title, sMessage - text in a box
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func MsgInfo(sMessage string, sTitle string, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	b, _ := json.Marshal(sMessage)
	sParams := fmt.Sprintf("[\"common\",\"minfo\",\"%s\",\"%s\",%s,\"%s\"]", sFunc, sName, string(b), sTitle)
	sendout(sParams)
}

// MsgStop creates a standard nessagebox
// sTitle - box title, sMessage - text in a box
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func MsgStop(sMessage string, sTitle string, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	b, _ := json.Marshal(sMessage)
	sParams := fmt.Sprintf("[\"common\",\"mstop\",\"%s\",\"%s\",%s,\"%s\"]", sFunc, sName, string(b), sTitle)
	sendout(sParams)
}

// MsgYesNo creates a standard nessagebox
// sTitle - box title, sMessage - text in a box
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func MsgYesNo(sMessage string, sTitle string, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	b, _ := json.Marshal(sMessage)
	sParams := fmt.Sprintf("[\"common\",\"myesno\",\"%s\",\"%s\",%s,\"%s\"]", sFunc, sName, string(b), sTitle)
	sendout(sParams)
}

// MsgGet creates a messagebox, which allows to input a string
// sTitle - box title, sMessage - text in a box, iStyle - a Winstyle for an "edit" widget (ES_PASSWORD, for example).
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func MsgGet(sMessage string, sTitle string, iStyle int32, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	b, _ := json.Marshal(sMessage)
	sParams := fmt.Sprintf("[\"common\",\"mget\",\"%s\",\"%s\",%s,\"%s\",%d]", sFunc, sName, string(b), sTitle, iStyle)
	sendout(sParams)
}

// Choice creates a dialog with a "browse" inside, which allows to select one of items in
// a passed slice arr.
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func Choice(arr []string, sTitle string, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	b, _ := json.Marshal(arr)
	sParams := fmt.Sprintf("[\"common\",\"mchoi\",\"%s\",\"%s\",%s,\"%s\"]", sFunc, sName, string(b), sTitle)
	sendout(sParams)
}

// SelectFile creates a standard dialog to select file
// sPath - initial path;
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func SelectFile(sPath string, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"cfile\",\"%s\",\"%s\",\"%s\"]", sFunc, sName, sPath)
	sendout(sParams)
}

// SelectColor creates a standard dialog to select color
// iColor - base color;
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func SelectColor(iColor int32, fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"ccolor\",\"%s\",\"%s\",%d]", sFunc, sName, iColor)
	sendout(sParams)
}

// SelectFont creates a standard dialog to select font
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier;
// sName - a parameter, passed to a callback procedure.
func SelectFont(fu func([]string) string, sFunc string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
	}
	pFont := &(Font{Name: sName})
	pFont.new()
	sParams := fmt.Sprintf("[\"common\",\"cfont\",\"%s\",\"%s\"]", sFunc, pFont.Name)
	sendout(sParams)
}

// InsertNode inserts a node to a tree widget pTree.
// sNodeName - a name of a parent node, or empty, if it is root node;
// sNodeNew - a name of an inserted node;
// sTitle - a caption of an inserted node;
// sNodeNext - a name of a node, you want to insert the new before;
// aImages - path to images for the node ( unselected, selected );
// fu, sCode - a definition of a callback procedure; fu - function, sCode - identifier or Harbour script.
func InsertNode(pTree *Widget, sNodeName string, sNodeNew string, sTitle string,
	sNodeNext string, aImages []string, fu func([]string) string, sCode string) {

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"node\",[\"%s\",\"%s\",\"%s\",\"%s\",",
		widgFullName(pTree), sNodeName, sNodeNew, sTitle, sNodeNext)
	if sCode != "" {
		sName := widgFullName(pTree)
		if fu != nil {
			RegFunc(sCode, fu)
			sCode = "pgo(\"" + sCode + "\",{\"" + sName + "\",\"" + sNodeNew + "\"})"
		}
		b, _ := json.Marshal(sCode)
		sCode = string(b)
	} else {
		sCode = "null"
	}

	if aImages == nil {
		sParams += "null"
	} else {
		b, _ := json.Marshal(aImages)
		sParams += string(b)
	}
	sParams += "," + sCode + "]]"

	sendout(sParams)
}

// PBarStep does a next step for a pPBar progress bar widget
func PBarStep(pPBar *Widget) {

	var sName = widgFullName(pPBar)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"step\",1]", sName)
	sendout(sParams)
}

// PBarSet sets a progress bar position
func PBarSet(pPBar *Widget, iPos int) {

	var sName = widgFullName(pPBar)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"setval\",%d]", sName, iPos)
	sendout(sParams)
}

// InitTray forces the main window to be placed in a tray,
// sIcon - a path to icon file,
// sMenuName - a name of a context menu,
// sTooltip - a tooltip for an icon in tray.
func InitTray(sIcon string, sMenuName string, sTooltip string) {

	sParams := fmt.Sprintf("[\"tray\",\"init\",\"%s\",\"%s\",\"%s\"]", sIcon, sMenuName, sTooltip)
	sendout(sParams)
}

// ModifyTrayIcon changes a tray icon of a main window,
// sIcon - a path to icon file.
func ModifyTrayIcon(sIcon string) {

	sParams := fmt.Sprintf("[\"tray\",\"icon\",\"%s\"]", sIcon)
	sendout(sParams)
}

// RadioEnd completes a group of radio buttons, started with a "radiogr" widget
func RadioEnd(p *Widget, iSel int) {

	var sName = widgFullName(p)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"radioend\",%d]", sName, iSel)
	sendout(sParams)
}

// TabPage initialises a new page of a tab widget.
func TabPage(pTab *Widget, sCaption string) {

	var sName = widgFullName(pTab)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"pagestart\",\"%s\"]", sName, sCaption)
	sendout(sParams)
}

// TabPageEnd completes a description of a page of a tab widget.
func TabPageEnd(pTab *Widget) {

	var sName = widgFullName(pTab)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"pageend\",1]", sName)
	sendout(sParams)
}

// BrwSetArray sets a two-dimensional slice to be represented in a browse widget p.
func BrwSetArray(p *Widget, arr *[][]string) {

	var sName = widgFullName(p)
	b, _ := json.Marshal(*arr)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwarr\",%s]", sName, string(b))
	sendout(sParams)
}

// BrwGetArray returns a two-dimensional slice from a browse widget p.
func BrwGetArray(p *Widget) [][]string {

	var sName = widgFullName(p)
	var arr [][]string

	sParams := fmt.Sprintf("[\"get\",\"%s\",\"brwarr\"]", sName)
	b := sendoutAndReturn(sParams)

	err := json.Unmarshal(b[1:], &arr)
	if err != nil {
		return nil
	} else {
		return arr
	}
}

// BrwSetColumn defines options for a column with number ic of a browse widget p.
// The options are:
//   sHead string - a column title;
//   iAlignHead int - the alignment of a column title ( 0 - left, 1 - center, 2 - right );
//   iAlignData int - the alignment of a column data ( 0 - left, 1 - center, 2 - right );
//   bEditable bool - is the data in a column editable.
//   iLength - column width in characters;
func BrwSetColumn(p *Widget, ic int, sHead string, iAlignHead int, iAlignData int,
	bEditable bool, iLength int) {
	var sName = widgFullName(p)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwcol\",[%d,\"%s\",%d,%d,%t,%d]]",
		sName, ic, sHead, iAlignHead, iAlignData, bEditable, iLength)
	sendout(sParams)
}

// BrwSetColumnEx sets options for a column with number ic of a browse widget p -
// those, which can not be set via BrwSetColumn.
// sParam - option name, xParam - option value
func BrwSetColumnEx(p *Widget, ic int, sParam string, xParam interface{}) {
	var sName = widgFullName(p)
	var sParValue string
	var sObj = "d"

	switch v := xParam.(type) {
	case *Font:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case *Style:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case CodeBlock:
		b, _ := json.Marshal(xParam)
		sParValue = string(b)
		sObj = "b"
	default:
		b, _ := json.Marshal(xParam)
		sParValue = string(b)
	}

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwcolx\",[%d,\"%s\",%s,\"%s\"]]",
		sName, ic, sParam, sParValue, sObj)
	sendout(sParams)
}

// BrwDelColumn deletes a column with number ic of a browse widget p.
func BrwDelColumn(p *Widget, ic int) {
	var sName = widgFullName(p)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwcoldel\",%d]", sName, ic)
	sendout(sParams)
}

// SetVar sets a variable value
func SetVar(sVarName string, sValue string) {

	sParams := fmt.Sprintf("[\"setvar\",\"%s\",\"%s\"]", sVarName, sValue)
	sendout(sParams)
}

// GetVar gets a variable value
func GetVar(sVarName string) string {

	var sRes string
	sParams := fmt.Sprintf("[\"getvar\",\"%s\"]", sVarName)
	b := sendoutAndReturn(sParams)
	if b[0] == byte('+') {
		b = b[1:len(b)]
	}
	err := json.Unmarshal(b, &sRes)
	if err != nil {
		return ""
	}

	return sRes
}

// SetImagePath sets a directory where GuiServer should look for image files.
func SetImagePath(sValue string) {

	sParams := fmt.Sprintf("[\"setparam\",\"bmppath\",\"%s\"]", sValue)
	sendout(sParams)
}

// SetPath sets a directory where GuiServer should write files
// and look for files to read.
func SetPath(sValue string) {

	sParams := fmt.Sprintf("[\"setparam\",\"path\",\"%s\"]", sValue)
	sendout(sParams)
}

// SetDateFormat sets a date display format,
// for example, "DD.MM.YYYY"
func SetDateFormat(sValue string) {

	sParams := fmt.Sprintf("[\"setparam\",\"datef\",\"%s\"]", sValue)
	sendout(sParams)
}

func (p *Font) new() *Font {
	if p.Name == "" {
		p.Name = fmt.Sprintf("f%d", iIdCount)
		iIdCount++
	}
	if aFonts == nil {
		aFonts = make([]*Font, 0, 16)
	}
	aFonts = append(aFonts, p)
	return p
}

func (p *Font) FillFont(arr []string) {
	p.Family = arr[1]
	i, _ := strconv.Atoi(arr[2])
	p.Height = int(i)
	p.Bold = (arr[3] == "t")
	p.Italic = (arr[4] == "t")
	p.Underline = (arr[5] == "t")
	p.Strikeout = (arr[6] == "t")
	i, _ = strconv.Atoi(arr[7])
	p.Charset = int16(i)
}

// Method Activate shows on the screen a main window or a dialog
func (o *Widget) Activate() bool {
	var sParams string
	if o.Type == "main" {
		sParams = fmt.Sprintf("[\"actmainwnd\",[\"f\"]]")
	} else if o.Type == "dialog" {
		sParams = fmt.Sprintf("[\"actdialog\",\"%s\",\"f\",[\"f\"]]", o.Name)
	} else {
		return false
	}
	b := sendout("" + sParams)
	if o.Type == "main" {
		wait()
	}
	return b
}

// Method Close closes a main window or a dialog
func (o *Widget) Close() bool {
	if o.Type == "main" || o.Type == "dialog" {
		sParams := fmt.Sprintf("[\"close\",\"%s\"]", o.Name)
		b := sendout("" + sParams)
		if o.Type == "dialog" {
			o.delete()
		}
		return b
	}
	return false
}

func (o *Widget) delete() bool {
	if o.Type == "dialog" {
		for i, od := range aDialogs {
			if o.Name == od.Name {
				aDialogs = append(aDialogs[:i], aDialogs[i+1:]...)
				return true
			}
		}
	} else if o.Type != "main" {
	}
	return false
}

// Method AddWidget adds new child widget
// o - parent window or widget
// pWidg - a Widget structure with definition of a new widget
func (o *Widget) AddWidget(pWidg *Widget) *Widget {
	pWidg.Parent = o
	mwidg, bOk := mWidgs[pWidg.Type]
	if !bOk {
		WriteLog(fmt.Sprintf("Error! \"%s\" does not defined\r\n", pWidg.Type))
		return nil
	}
	if pWidg.Name == "" {
		pWidg.Name = fmt.Sprintf("w%d", iIdCount)
		iIdCount++
	}

	sPar2 := setprops(pWidg, mwidg)
	sParams := fmt.Sprintf("[\"addwidg\",\"%s\",\"%s\",[%d,%d,%d,%d,\"%s\"]%s]",
		pWidg.Type, widgFullName(pWidg), pWidg.X, pWidg.Y, pWidg.W,
		pWidg.H, pWidg.Title, sPar2)
	sendout(sParams)
	PLastWidget = pWidg
	if o.aWidgets == nil {
		o.aWidgets = make([]*Widget, 0, 16)
	}
	o.aWidgets = append(o.aWidgets, pWidg)
	return pWidg
}

// Method SetText sets a text aText to a widget, pointed by o.
func (o *Widget) SetText(sText string) {

	var sName = widgFullName(o)
	o.Title = sText
	b, _ := json.Marshal(sText)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"text\",%s]", sName, string(b))
	sendout(sParams)
}

// Method SetImage sets an image widget, pointed by o,
// sImage - a path to an image
func (o *Widget) SetImage(sImage string) {

	var sName = widgFullName(o)

	mwidg, bOk := mWidgs[o.Type]
	if !bOk {
		return
	}
	_, bOk = mwidg["Image"]
	if !bOk {
		return
	}

	if o.AProps == nil {
		o.AProps = make(map[string]string)
	}
	o.AProps["Image"] = sImage
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"image\",\"%s\"]", sName, sImage)
	sendout(sParams)
}

// Method SetParam sets a property to a widget, pointed by o,
// sParam - a name of a property,
// xParam - a value of a property.
func (o *Widget) SetParam(sParam string, xParam interface{}) {

	var sName = widgFullName(o)
	var sParValue string
	var sObj = "d"

	switch v := xParam.(type) {
	case *Font:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case *Style:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case *Widget:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case *Highlight:
		sParValue = "\"" + v.Name + "\""
		sObj = "o"
	case CodeBlock:
		b, _ := json.Marshal(xParam)
		sParValue = string(b)
		sObj = "b"
	default:
		b, _ := json.Marshal(xParam)
		sParValue = string(b)
	}
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"xparam\",[\"%s\",%s,\"%s\"]]", sName, sParam, sParValue, sObj)
	sendout(sParams)
}

// Method GetText gets the text from a widget, pointed by o.
func (o *Widget) GetText() string {
	var sRes string
	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"get\",\"%s\",\"text\"]", sName)
	b := sendoutAndReturn(sParams)
	if b[0] == byte('+') {
		b = b[1:len(b)]
	}
	err := json.Unmarshal(b, &sRes)
	if err != nil {
		return ""
	}

	return sRes
}

// Method SetColor sets a text color tColor and background color bColor to a widget, pointed by o.
func (o *Widget) SetColor(tColor int32, bColor int32) {

	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"color\",[%d,%d]]", sName, tColor, bColor)
	sendout(sParams)
}

// Method SetFont sets a font pFont to a widget, pointed by o.
func (o *Widget) SetFont(pFont *Font) {

	var sName = widgFullName(o)
	o.Font = pFont
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"font\",\"%s\"]", sName, pFont.Name)
	sendout(sParams)
}

func (o *Widget) SetCallBackProc(sbName string, fu func([]string) string, sCode string, params ...string) {

	var sName = widgFullName(o)
	var sc1, sc2 string

	if fu != nil {
		RegFunc(sCode, fu)
		if sbName == "onposchanged" {
			sc1 = "o,n"
			sc2 = ",n"
		} else if sbName == "onrclick" || sbName == "onenter" {
			sc1 = "o,nc,nr"
			sc2 = ",nc,nr"
		} else {
			sc1 = ""
			sc2 = ""
		}
		sCode = "{|" + sc1 + "|pgo(\"" + sCode + "\",{\"" + sName + "\"" + sc2
		for _, v := range params {
			sCode += ",\"" + v + "\""
		}
		sCode += "})}"
	}
	b, _ := json.Marshal(sCode)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"cb.%s\",%s]", sName, sbName, string(b))
	sendout(sParams)
}

func (o *Widget) SetCallBackFunc(sbName string, fu func([]string) string, sCode string, params ...string) {

	var sName = widgFullName(o)

	if fu != nil {
		RegFunc(sCode, fu)
		sCode = "fgo(\"" + sCode + "\",{\"" + sName + "\""
		for _, v := range params {
			sCode += ",\"" + v + "\""
		}
		sCode += "})"
	}
	b, _ := json.Marshal(sCode)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"cb.%s\",%s]", sName, sbName, string(b))
	sendout(sParams)
}

func (o *Widget) Move(iLeft, iTop, iWidth, iHeight int32) {

	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"move\",[%d,%d,%d,%d]]", sName, iLeft, iTop, iWidth, iHeight)
	sendout(sParams)
}

func (o *Widget) Enable(bEnable bool) {

	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"enable\",%t]", sName, bEnable)
	sendout(sParams)
}

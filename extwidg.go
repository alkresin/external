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

const (
	DT_LEFT     = 0
	DT_CENTER   = 1
	DT_RIGHT    = 2
)

// A set of constants of the printer paper types
const (
	DMPAPER_A3  =  8  // A3 297 x 420 mm
	DMPAPER_A4  =  9  // A4 210 x 297 mm
	DMPAPER_A5  = 11  // A5 148 x 210 mm
	DMPAPER_A6  = 70  // A6 105 x 148 mm
)

// The Font structure prepares data to create a new font
type Font struct {
	Family    string
	Name      string
	Height    int16
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

// The Printer structure prepares data to initialize a printer
type Printer struct {
	Name     string
	SPrinter string
	BPreview   bool
	IFormType   int
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

var PLastWindow *Widget
var PLastWidget *Widget
var PLastPrinter *Printer

var mWidgs = make(map[string]map[string]string)

func init() {
	mWidgs["main"] = nil
	mWidgs["dialog"] = nil
	mWidgs["label"] = map[string]string{"Transpa": "L"}
	mWidgs["edit"] = map[string]string{"Picture": "C"}
	mWidgs["button"] = nil
	mWidgs["check"] = map[string]string{"Transpa": "L"}
	mWidgs["radio"] = map[string]string{"Transpa": "L"}
	mWidgs["radiogr"] = nil
	mWidgs["group"] = nil
	mWidgs["combo"] = map[string]string{"AItems": "AC"}
	mWidgs["bitmap"] = map[string]string{"Transpa": "L", "TrColor": "N", "Image": "C"}
	mWidgs["line"] = map[string]string{"Vertical": "L"}
	mWidgs["panel"] = map[string]string{"HStyle": "C"}
	mWidgs["paneltop"] = map[string]string{"HStyle": "C"}
	mWidgs["panelbot"] = map[string]string{"HStyle": "C", "AParts": "AC"}
	mWidgs["ownbtn"] = map[string]string{"Transpa": "L", "TrColor": "N", "Image": "C", "HStyles": "AC"}
	mWidgs["splitter"] = map[string]string{"Vertical": "L", "From": "N", "TO": "N", "ALeft": "AC", "ARight": "AC"}
	mWidgs["updown"] = map[string]string{"From": "N", "TO": "N"}
	mWidgs["tree"] = map[string]string{"AImages": "AC", "EditLabel": "L"}
	mWidgs["progress"] = map[string]string{"Maxpos": "N"}
	mWidgs["tab"] = nil
	mWidgs["browse"] = map[string]string{"Append": "L", "Autoedit": "L"}
}

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
					sPar += fmt.Sprintf(",\"%s\": %d", name, val)
				} else if cType == "AC" {
					sPar += fmt.Sprintf(",\"%s\": %s", name, val)
				}
			} else {
				WriteLog(fmt.Sprintf("Error! \"%s\" does not defined for \"%s\"", name, pWidg.Type))
				return ""
			}
		}
	}
	if sPar != "" {
		sPar = ",{" + sPar[1:] + "}"
	}
	return sPar
}

// Converts function arguments to a json string
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

// Reads a main window description from a xml file, prepared by HwGUI's Designer,
// initialises and activates this window with all its widgets
func OpenMainForm(sForm string) bool {
	var b bool
	b = sendout("[\"openformmain\",\"" + sForm + "\"]")
	wait()
	return b
}

// Reads a dialog window description from a xml file, prepared by HwGUI's Designer,
// initialises and activates this dialog with all its widgets
func OpenForm(sForm string) bool {
	var b bool
	b = sendout("[\"openform\",\"" + sForm + "\"]")
	return b
}

// Reads a report description from a xml file, prepared by HwGUI's Designer
// and prints this report
func OpenReport(sForm string) bool {
	var b bool
	b = sendout("[\"openreport\",\"" + sForm + "\"]")
	return b
}

// Creates a font with parameters, defined in a structure, pointed by pFont argument.
func CreateFont(pFont *Font) *Font {

	pFont.new()
	sParams := fmt.Sprintf("[\"crfont\",\"%s\",\"%s\",%d,%t,%t,%t,%t,%d]", pFont.Name, pFont.Family, pFont.Height,
		pFont.Bold, pFont.Italic, pFont.Underline, pFont.Strikeout, pFont.Charset)
	sendout(sParams)
	return pFont
}

// Creates a style with parameters, defined in a structure, pointed by pStyle argument.
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

// Creates a style with parameters, defined in a structure, pointed by pStyle argument.
func InitPrinter(pPrinter *Printer, sFunc string, fu func([]string) string, sName string) *Printer {

	if pPrinter.Name == "" {
		pPrinter.Name = fmt.Sprintf("p%d", iIdCount)
		iIdCount++
	}
	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"prninit\",\"%s\",[\"%s\",%t,%d,%t],\"%s\",\"%s\"]", pPrinter.Name,
		pPrinter.SPrinter, pPrinter.BPreview, pPrinter.IFormType, pPrinter.BLandscape, sFunc, sName )
	sendout(sParams)
	PLastPrinter = pPrinter
	return pPrinter
}

func (p *Printer) AddFont(pFont *Font) *Font {
	pFont.new()
	sParams := fmt.Sprintf("[\"print\",\"fontadd\",\"%s\",[\"%s\",\"%s\",%d,%t,%t,%t,%d]]", p.Name,
		pFont.Name, pFont.Family, pFont.Height,
		pFont.Bold, pFont.Italic, pFont.Underline, pFont.Charset)
	sendout(sParams)
	return pFont
}

func (p *Printer) SetFont(pFont *Font) {
	sParams := fmt.Sprintf("[\"print\",\"fontset\",\"%s\",[\"%s\"]]", p.Name, pFont.Name)
	sendout(sParams)
}

func (p *Printer) Say(iTop, iLeft, iRight, iBottom int32, sText string, iOpt int32) {

	sParams := fmt.Sprintf("[\"print\",\"text\",\"%s\",[\"%s\",%d,%d,%d,%d,%d]]",
		p.Name, sText, iTop, iLeft, iRight, iBottom, iOpt)
	sendout(sParams)
}

func (p *Printer) Line(iTop, iLeft, iRight, iBottom int32) {

	sParams := fmt.Sprintf("[\"print\",\"line\",\"%s\",[%d,%d,%d,%d]]", p.Name, iTop, iLeft, iRight, iBottom)
	sendout(sParams)
}

func (p *Printer) Box(iTop, iLeft, iRight, iBottom int32) {

	sParams := fmt.Sprintf("[\"print\",\"box\",\"%s\",[%d,%d,%d,%d]]", p.Name, iTop, iLeft, iRight, iBottom)
	sendout(sParams)
}

func (p *Printer) StartPage() {

	sParams := fmt.Sprintf("[\"print\",\"startpage\",\"%s\",[]]", p.Name)
	sendout(sParams)
}

func (p *Printer) EndPage() {

	sParams := fmt.Sprintf("[\"print\",\"endpage\",\"%s\",[]]", p.Name)
	sendout(sParams)
}

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

func EvalProc(s string) {

	b, _ := json.Marshal(s)
	sendout("[\"evalcode\"," + string(b) + "]")
}

func EvalFunc(s string) []byte {

	b, _ := json.Marshal(s)
	b = sendoutAndReturn("[\"evalcode\"," + string(b) + ",\"t\"]")
	if b[0] == byte('+') && b[1] == byte('"') {
		b = b[2 : len(b)-1]
	}
	return b
}

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

func MsgInfo(sMessage string, sTitle string, sFunc string, fu func([]string) string, sName string) {

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

func MsgStop(sMessage string, sTitle string, sFunc string, fu func([]string) string, sName string) {

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

func MsgYesNo(sMessage string, sTitle string, sFunc string, fu func([]string) string, sName string) {

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

func MsgGet(sMessage string, sTitle string, iStyle int32, sFunc string, fu func([]string) string, sName string) {

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

func Choice(arr []string, sTitle string, sFunc string, fu func([]string) string, sName string) {

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

func SelectFile(sPath string, sFunc string, fu func([]string) string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"cfile\",\"%s\",\"%s\",\"%s\"]", sFunc, sName, sPath)
	sendout(sParams)
}

func SelectColor(iColor int32, sFunc string, fu func([]string) string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"ccolor\",\"%s\",\"%s\",%d]", sFunc, sName, iColor)
	sendout(sParams)
}

func SelectFont(sFunc string, fu func([]string) string, sName string) {

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
	sParams += "," + sCode + "]"

	sendout(sParams)
}

func PBarStep(pPBar *Widget) {

	var sName = widgFullName(pPBar)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"step\",1]", sName)
	sendout(sParams)
}

func PBarSet(pPBar *Widget, iPos int) {

	var sName = widgFullName(pPBar)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"setval\",%d]", sName, iPos)
	sendout(sParams)
}

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
func BrwSetArray(p *Widget, arr [][]string) {

	var sName = widgFullName(p)
	b, _ := json.Marshal(arr)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwarr\",%s]", sName, string(b))
	sendout(sParams)
}

// BrwSetColumn defines options for a column with number ic of a browse widget p.
// The options are:
//   sHead string - a column title;
//   iAlignHead int - the alignment of a column title ( 0 - left, 1 - center, 2 - right );
//   iAlignData int - the alignment of a column data ( 0 - left, 1 - center, 2 - right );
//   bEditable bool - is the data in a column editable.
func BrwSetColumn(p *Widget, ic int, sHead string, iAlignHead int, iAlignData int, bEditable bool) {
	var sName = widgFullName(p)
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"brwcol\",[%d,\"%s\",%d,%d,%t]]", sName, ic, sHead, iAlignHead, iAlignData, bEditable)
	sendout(sParams)
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
	p.Height = int16(i)
	p.Bold = (arr[3] == "t")
	p.Italic = (arr[4] == "t")
	p.Underline = (arr[5] == "t")
	p.Strikeout = (arr[6] == "t")
	i, _ = strconv.Atoi(arr[7])
	p.Charset = int16(i)
}

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

func (o *Widget) Close() bool {
	if o.Type == "main" || o.Type == "dialog" {
		sParams := fmt.Sprintf("[\"close\",\"%s\"]", o.Name)
		b := sendout("" + sParams)
		return b
	}
	return false
}

func (o *Widget) Delete() bool {
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

func (o *Widget) AddWidget(pWidg *Widget) *Widget {
	pWidg.Parent = o
	mwidg, bOk := mWidgs[pWidg.Type]
	if !bOk {
		WriteLog(fmt.Sprintf("Error! \"%s\" does not defined", pWidg.Type))
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

func (o *Widget) SetText(sText string) {

	var sName = widgFullName(o)
	o.Title = sText
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"text\",\"%s\"]", sName, sText)
	sendout(sParams)
}

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

func (o *Widget) SetParam(sParam string, xParam interface{}) {

	var sName = widgFullName(o)
	var sParValue string
	var bObj = true

	switch v := xParam.(type) {
	case *Font:
		sParValue = "\"" + v.Name + "\""
	case *Style:
		sParValue = "\"" + v.Name + "\""
	case *Widget:
		sParValue = "\"" + v.Name + "\""
	default:
		b, _ := json.Marshal(xParam)
		sParValue = string(b)
		bObj = false
	}
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"xparam\",[\"%s\",%s,%t]]", sName, sParam, sParValue, bObj)
	sendout(sParams)
}

func (o *Widget) GetText() string {
	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"get\",\"%s\",\"text\"]", sName)
	b := sendoutAndReturn(sParams)
	if b[0] == byte('+') && b[1] == byte('"') {
		b = b[2 : len(b)-1]
	}
	return string(b)
}

func (o *Widget) SetColor(tColor int32, bColor int32) {

	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"color\",[%d,%d]]", sName, tColor, bColor)
	sendout(sParams)
}

func (o *Widget) SetFont(pFont *Font) {

	var sName = widgFullName(o)
	o.Font = pFont
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"font\",\"%s\"]", sName, pFont.Name)
	sendout(sParams)
}

func (o *Widget) SetCallBackProc(sbName string, fu func([]string) string, sCode string, params ...string) {

	var sName = widgFullName(o)

	if fu != nil {
		RegFunc(sCode, fu)
		sCode = "pgo(\"" + sCode + "\",{\"" + sName + "\""
		for _, v := range params {
			sCode += ",\"" + v + "\""
		}
		sCode += "})"
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

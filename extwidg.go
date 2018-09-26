package external

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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

type Style struct {
	Name      string
	Orient    int16
	Colors    []int32
	Corners   []int32
	BorderW   int8
	BorderClr int32
	Bitmap    string
}

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
}

func widgFullName(pWidg *Widget) string {
	sName := pWidg.Name

	for pWidg.Parent != nil {
		pWidg = pWidg.Parent
		sName = pWidg.Name + "." + sName
	}
	return sName
}

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

func GetWnd(sName string) *Widget {
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

func GetWidg(sName string) *Widget {
	npos := strings.Index(sName, ".")
	if npos == -1 {
		return GetWnd(sName)
	}
	sWnd := sName[:npos]
	sName = sName[npos+1:]
	if oWnd := GetWnd(sWnd); oWnd != nil {
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

func ArrStrings(sParam ...string) string {
	s := ""
	for _, v := range sParam {
		s += ",\"" + v + "\""
	}
	return "[" + s[1:] + "]"
}

func ArrInts(iParam ...int) string {
	s := ""
	for _, v := range iParam {
		s += "," + strconv.Itoa(v)
	}
	return "[" + s[1:] + "]"
}

func ArrWidgs(wParam ...*Widget) string {
	s := ""
	for _, w := range wParam {
		s += ",\"" + w.Name + "\""
	}
	return "[" + s[1:] + "]"
}

func OpenMainForm(sForm string) bool {
	var b bool
	b = Sendout("[\"openformmain\",\"" + sForm + "\"]")
	Wait()
	return b
}

func OpenForm(sForm string) bool {
	var b bool
	b = Sendout("[\"openform\",\"" + sForm + "\"]")
	return b
}

func OpenReport(sForm string) bool {
	var b bool
	b = Sendout("[\"openreport\",\"" + sForm + "\"]")
	return b
}

func CreateFont(pFont *Font) *Font {

	pFont.newfont()
	sParams := fmt.Sprintf("[\"crfont\",\"%s\",\"%s\",%d,%t,%t,%t,%t,%d]", pFont.Name, pFont.Family, pFont.Height,
		pFont.Bold, pFont.Italic, pFont.Underline, pFont.Strikeout, pFont.Charset)
	Sendout(sParams)
	return pFont
}

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
	Sendout(sParams)
	return pStyle
}

func InitMainWindow(pWnd *Widget) bool {
	pMainWindow = pWnd
	PLastWindow = pWnd
	pWnd.Type = "main"
	pWnd.Name = "main"
	sPar2 := setprops(pWnd, mWidgs["main"])
	sParams := fmt.Sprintf("[\"crmainwnd\",[%d,%d,%d,%d,\"%s\"]%s]", pWnd.X, pWnd.Y, pWnd.W,
		pWnd.H, pWnd.Title, sPar2)
	return Sendout(sParams)
}

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
	return Sendout(sParams)
}

func EvalProc(s string) {

	b, _ := json.Marshal(s)
	Sendout("[\"evalcode\"," + string(b) + "]")
}

func EvalFunc(s string) []byte {

	b, _ := json.Marshal(s)
	b = SendoutAndReturn("[\"evalcode\"," + string(b) + ",\"t\"]")
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
	b := SendoutAndReturn(sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
}

func SelectFile(sPath string, sFunc string, fu func([]string) string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"cfile\",\"%s\",\"%s\",\"%s\"]", sFunc, sName, sPath)
	Sendout(sParams)
}

func SelectColor(iColor int32, sFunc string, fu func([]string) string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
		sName = ""
	}
	sParams := fmt.Sprintf("[\"common\",\"ccolor\",\"%s\",\"%s\",%d]", sFunc, sName, iColor)
	Sendout(sParams)
}

func SelectFont(sFunc string, fu func([]string) string, sName string) {

	if fu != nil && sFunc != "" {
		RegFunc(sFunc, fu)
	} else {
		sFunc = ""
	}
	pFont := &(Font{Name: sName})
	pFont.newfont()
	sParams := fmt.Sprintf("[\"common\",\"cfont\",\"%s\",\"%s\"]", sFunc, pFont.Name)
	Sendout(sParams)
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
		b,_ := json.Marshal(aImages)
		sParams +=  string(b)
	}
	sParams += "," + sCode + "]"

	Sendout(sParams)
}

func SetImagePath(sValue string) {

	sParams := fmt.Sprintf("[\"setparam\",\"bmppath\",\"%s\"]", sValue)
	Sendout(sParams)
}

func SetPath(sValue string) {

	sParams := fmt.Sprintf("[\"setparam\",\"path\",\"%s\"]", sValue)
	Sendout(sParams)
}

func (p *Font) newfont() *Font {
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
	b := Sendout("" + sParams)
	if o.Type == "main" {
		Wait()
	}
	return b
}

func (o *Widget) Close() bool {
	if o.Type == "main" || o.Type == "dialog" {
		sParams := fmt.Sprintf("[\"close\",\"%s\"]", o.Name)
		b := Sendout("" + sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
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

	o.AProps["Image"] = sImage
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"image\",\"%s\"]", sName, sImage)
	Sendout(sParams)
}

func (o *Widget) GetText() string {
	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"get\",\"%s\",\"text\"]", sName)
	b := SendoutAndReturn(sParams)
	if b[0] == byte('+') && b[1] == byte('"') {
		b = b[2 : len(b)-1]
	}
	return string(b)
}

func (o *Widget) SetColor(tColor int32, bColor int32) {

	var sName = widgFullName(o)

	sParams := fmt.Sprintf("[\"set\",\"%s\",\"color\",[%d,%d]]", sName, tColor, bColor)
	Sendout(sParams)
}

func (o *Widget) SetFont(pFont *Font) {

	var sName = widgFullName(o)
	o.Font = pFont
	sParams := fmt.Sprintf("[\"set\",\"%s\",\"font\",\"%s\"]", sName, pFont.Name)
	Sendout(sParams)
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
	Sendout(sParams)
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
	Sendout(sParams)
}

package external

import (
	"encoding/json"
)

var sMenu = ""
var iStackLen = 0

func Menu(sName string) {

	if sMenu == "" {
		sMenu = "[\"menu\",["
	} else {
		if sMenu[len(sMenu)-1] != '[' {
			sMenu += ","
		}
		sMenu += "[\"" + sName + "\",["
		iStackLen++
	}
}

func EndMenu() {
	sMenu += "]]"
	if iStackLen == 0 {
		Sendout(sMenu)
		sMenu = ""
	} else {
		iStackLen--
	}
}

func AddMenuItem(sName string, fu func([]string)string, sCode string, params ...string) {
	if fu != nil {
		RegFunc(sCode, fu)
		sCode = "pgo(\"" + sCode + "\",{"
		for i, v := range params {
			if i > 0 {
				sCode += ","
			}
			sCode += "\"" + v + "\""
		}
		sCode += "})"
	}
	b, _ := json.Marshal(sCode)
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\""+sName+"\"," + string(b) + "]"
}

func AddMenuSeparator() {
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\"-\"]"
}

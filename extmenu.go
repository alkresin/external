// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a GNU general public
// license that can be found in the LICENSE file.

package external

import (
	"encoding/json"
)

var sMenu = ""
var iStackLen = 0

// Menu starts a window's menu or submenu definition, sName is a menu title.
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
		sendout(sMenu)
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

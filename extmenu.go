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

// EndMenu completes a window's menu or submenu definition
func EndMenu() {
	sMenu += "]]"
	if iStackLen == 0 {
		sendout(sMenu)
		sMenu = ""
	} else {
		iStackLen--
	}
}

// AddMenuItem adds a new item to the Window's menu or submenu,
// sName argument is a title of the item,
// fu - a function in the program, which must be called, when this menu item is selected,
// sCode - the identifier (name) of this function.
// If the fu value is nil, sCode contains the Harbour's code, which must be executed by
// the GuiServer when this menu item is selected.
// params - arguments for the fu function.
func AddMenuItem(sName string, fu func([]string) string, sCode string, params ...string) {
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
	sMenu += "[\"" + sName + "\"," + string(b) + "]"
}

// AddMenuSeparator adds a separator to the Window's menu or submenu,
func AddMenuSeparator() {
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\"-\"]"
}

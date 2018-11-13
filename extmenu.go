// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package external

import (
	"strconv"
	"encoding/json"
)

var sMenu = ""
var iStackLen = 0

// Menu starts a window's menu or submenu definition, sTitle is a menu title.
func Menu(sTitle string) {

	if sMenu == "" {
		sMenu = "[\"menu\",["
	} else {
		if sMenu[len(sMenu)-1] != '[' {
			sMenu += ","
		}
		sMenu += "[\"" + sTitle + "\",["
		iStackLen++
	}
}

// MenuContext starts a context menu, sName is a menu identifier
func MenuContext(sName string) {

	if sMenu == "" {
		sMenu = "[\"menucontext\",\"create\",\"" + sName + "\",["
	}
}

// Show context menu on the screen
func ShowMenuContext(sName string, pWnd *Widget) {
	var sWndName string
	if pWnd == nil {
		sWndName = ""
	} else {
		sWndName = pWnd.Name
	}
	sendout("[\"menucontext\",\"show\",\"" + sName + "\",\"" + sWndName + "\"]")
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

func getscode(fu func([]string) string, sCode string, params ...string) string {
	if fu != nil {
		RegFunc(sCode, fu)
		sCode = "pgo(\"" + sCode + "\",{\"menu\""
		for _, v := range params {
			sCode += ",\"" + v + "\""
		}
		sCode += "})"
	}
	b, _ := json.Marshal(sCode)
	return string(b)
}

// AddMenuItem adds a new item to the Window's menu or submenu,
// sName argument is a title of the item,
// id - menu item identifier; if 0 - it is created automatically
// fu - a function in the program, which must be called, when this menu item is selected,
// sCode - the identifier (name) of this function.
// If the fu value is nil, sCode contains the Harbour's code, which must be executed by
// the GuiServer when this menu item is selected.
// params - arguments for the fu function.
func AddMenuItem(sName string, id int, fu func([]string) string, sCode string, params ...string) {
	
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\"" + sName + "\"," + getscode(fu, sCode, params...) + "," + strconv.Itoa(id) + "]"
}

func AddCheckMenuItem(sName string, id int, fu func([]string) string, sCode string, params ...string) {
	
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\"" + sName + "\"," + getscode(fu, sCode, params...) + "," + strconv.Itoa(id) + ",true]"
}

// AddMenuSeparator adds a separator to the Window's menu or submenu,
func AddMenuSeparator() {
	if sMenu[len(sMenu)-1] != '[' {
		sMenu += ","
	}
	sMenu += "[\"-\"]"
}

func MenuItemEnable(sWndName string, sMenuName string, iItem int, bValue bool) {

	sendout("[\"menu\",\"enable\",\"" + sWndName + "\",\"" + sMenuName + "\"," +
		strconv.Itoa(iItem) + "," + strconv.FormatBool(bValue) + "]")
}

func MenuItemCheck(sWndName string, sMenuName string, iItem int, bValue bool) {

	sendout("[\"menu\",\"check\",\"" + sWndName + "\",\"" + sMenuName + "\"," +
		strconv.Itoa(iItem) + "," + strconv.FormatBool(bValue) + "]")
}

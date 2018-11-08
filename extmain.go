// Copyright 2018 Alexander S.Kresin <alex@kresin.ru>, http://www.kresin.ru
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Package external is a GUI framework for Go language.
// External is a Go library to build GUI application, using a
// standalone GUI server application: https://github.com/alkresin/guiserver.
package external

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	VerProto = "1.1"
	Version  = "1.1"
)

var sLogName = "egui.log"
var bEndProg = false

var connOut, connIn net.Conn
var bConnExist = false
var bPacket = false
var sPacketBuf string

// Init runs, if needed, the Guiserver application, and connects to it.
// It returns 0, if the connection is successful, 1 - in other case,
// 2 -if a protocol version of a GuiServer isn't equal to a local protocol.
// The sOpt argument specifies connection details. It may contain following strings:
// guiserver=<full path to GuiServer executable>
// ip=<ip address of a computer, where GuiServer runs>
// port=<tcp/ip port number>
// log (switch on logging)
func Init(sOpt string) int {

	var err error

	iPort := 3101
	sServer := "guiserver"
	sIp := "127.0.0.1"
	sLog := ""
	if sOpt != "" {
		var arr []string
		sep := "\r\n"
		if !strings.Contains(sOpt, sep) {
			if sep = "\n"; !strings.Contains(sOpt, sep) {
				sep = ""
			}
		}
		if sep == "" {
			arr = make([]string, 1, 1)
			arr[0] = sOpt
		} else {
			arr = strings.Split(sOpt, sep)
		}
		for i := 0; i < len(arr); i++ {
			s := strings.ToLower(arr[i])
			if len(s) > 9 && s[:9] == "guiserver" {
				sServer = strings.TrimSpace(s[10:])
			} else if len(s) > 8 && s[:7] == "address" {
				sIp = strings.TrimSpace(s[8:])
			} else if len(s) > 5 && s[:4] == "port" {
				iPort, _ = strconv.Atoi(strings.TrimSpace(s[5:]))
			} else if len(s) > 2 && s[:3] == "log" {
				if s[4:5] == "1" {
					sLog = "-log1"
				} else if s[4:5] == "2" {
					sLog = "-log2"
				}
			}
		}
	}

	buf := make([]byte, 128)

	if sServer != "" {
		cmd := exec.Command(sServer, fmt.Sprintf("-p%d", iPort), sLog)
		cmd.Start()
	}
	time.Sleep(100 * time.Millisecond)

	connOut, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort))
	if err != nil {
		time.Sleep(1000 * time.Millisecond)
		connOut, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort))
		if err != nil {
			time.Sleep(3000 * time.Millisecond)
			connOut, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort))
			WriteLog(fmt.Sprintln(sServer, sIp, iPort))
			WriteLog(fmt.Sprintln(err))
			return 1
		}
	}
	iBufLen, err := connOut.Read(buf)
	if err != nil {
		WriteLog(fmt.Sprintln(err))
		connOut.Close()
		return 1
	}
	sVer := string(buf[:iBufLen-1])
	sVer = sVer[(strings.Index(sVer, "/") + 1):]

	connIn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort+1))
	if err != nil {
		time.Sleep(1000 * time.Millisecond)
		connOut, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort+1))
		if err != nil {
			WriteLog(fmt.Sprintln(sServer, sIp, iPort+1))
			WriteLog(fmt.Sprintln(err))
			return 1
		}
	}
	_, err = connIn.Read(buf)
	if err != nil {
		WriteLog(fmt.Sprintln(err))
		connIn.Close()
		return 1
	}

	if sVer != VerProto {
		WriteLog("\r\nProtocol version mismatched. Need " + VerProto + ", received " + sVer)
		connIn.Close()
		connOut.Close()
		return 2
	}

	bConnExist = true
	go listen(iPort + 1)
	time.Sleep(100 * time.Millisecond)

	return 0

}

// Exit closes the connection to Guiserver.
func Exit() {
	if bConnExist {
		connOut.Close()
	}
}

func listen(iPort int) {

	var bErr bool

	buffer := make([]byte, 1024)
	arr := make([]string, 5)
	for {

		bErr = false
		//length, err := connIn.Read(buffer)
		length, err := read(connIn, &buffer)

		if err != nil {
			//WriteLog("Read error\r\n")
			return
		}

		if buffer[0] != byte('+') || buffer[length-1] != byte('\n') {
			bErr = true
		}

		if !bErr {
			err = json.Unmarshal(buffer[1:length-1], &arr)
			if err != nil {
				WriteLog("Unmarshal error\r\n")
				bErr = true
			}
		}

		//fmt.Printf("Received command %d\t:%s %d\n", length, string(buffer[:length]), len(arr))

		if !bErr && len(arr) > 0 {
			switch arr[0] {
			case "runproc":
				sendResponse(connIn, "[\"Ok\"]")
				if len(arr) > 1 {
					if fnc, bExist := mfu[arr[1]]; bExist {
						var ap []string
						if len(arr) > 2 {
							ap = make([]string, 5)
							err = json.Unmarshal([]byte(arr[2]), &ap)
							if err != nil {
								WriteLog(fmt.Sprintf("runproc param Unmarshal error (%s)\r\n", arr[2]))
							}
						}
						//WriteLog(fmt.Sprintf("pgo> (%s) len:%d\r\n",arr[2],len(ap) ))
						fnc(ap)
					}
				} else {
					bErr = true
				}
			case "runfunc":
				if len(arr) > 1 {
					if fnc, bExist := mfu[arr[1]]; bExist {
						var ap []string
						if len(arr) > 2 {
							ap = make([]string, 5)
							err = json.Unmarshal([]byte(arr[2]), &ap)
							if err != nil {
								WriteLog(fmt.Sprintf("runproc param Unmarshal error (%s)\r\n", arr[2]))
							}
						}
						//WriteLog(fmt.Sprintf("pgo> (%s) len:%d\r\n",arr[2],len(ap) ))
						s := fnc(ap)
						b, _ := json.Marshal(s)
						sendResponse(connIn, "[\""+string(b)+"\"]")
					} else {
						sendResponse(connIn, "[\"Err\"]")
					}
				} else {
					bErr = true
					sendResponse(connIn, "[\"Err\"]")
				}
			case "exit":
				sendResponse(connIn, "[\"Ok\"]")
				if len(arr) > 1 {
					oW := Wnd(arr[1])
					if oW != nil {
						oW.Delete()
					}
				} else {
					bErr = true
				}
			case "endapp":
				sendResponse(connIn, "[\"Goodbye\"]")
				time.Sleep(100 * time.Millisecond)
				bEndProg = true
				connIn.Close()
				return
			default:
				sendResponse(connIn, "[\"Error\"]")
				bErr = true
			}
		}
		if bErr {
			WriteLog(fmt.Sprintf("Wrong message: %s]\r\n", string(buffer[:length])))
		}
	}
}

func read(conn net.Conn, pBuff *[]byte) (int, error) {
	*pBuff = (*pBuff)[:0]
	tmp := make([]byte, 256)
	for {
		length, err := conn.Read(tmp)
		if err != nil {
			WriteLog("Read error\r\n")
			return 0, err
		}
		*pBuff = append(*pBuff, tmp[:length]...)
		if tmp[length-1] == '\n' {
			break
		}
	}
	return len(*pBuff), nil
}

func sendResponse(conn net.Conn, s string) {
	conn.Write([]byte("+" + s + "\n"))
}

func sendout(s string) bool {

	var err error

	if bPacket {
		sPacketBuf += "," + s
	} else {
		if !bConnExist {
			WriteLog( "sendout: No connection established.\r\n" )
			return false
		}

		_, err = connOut.Write([]byte("+" + s + "\n"))
		if err != nil {
			fmt.Println(err)
			return false
		}

		buf := make([]byte, 128)
		_, err = connOut.Read(buf)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

func sendoutAndReturn(s string) []byte {

	var err error
	buf := make([]byte, 1024)

	if !bConnExist {
		WriteLog( "sendoutAndReturn: No connection established.\r\n" )
		return []byte("")
	}

	_, err = connOut.Write([]byte("+" + s + "\n"))
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	//length, err = connOut.Read(buf)
	length, err := read(connOut, &buf)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	return buf[:length-1]
}

// BeginPacket begins a sequence of functions, which creates or modifies GUI elements,
// for to join messages to Guiserver to one packet.
func BeginPacket() {
	bPacket = true
	sPacketBuf = "[\"packet\""
}

// EndPacket completes a sequence of functions, started by BeginPacket
func EndPacket() {
	bPacket = false
	sendout(sPacketBuf + "]")
	sPacketBuf = ""
}

// WriteLog writes the sText to a log file egui.log.
func WriteLog(sText string) {

	f, err := os.OpenFile(sLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(sText)

}

// RegFunc adds the fu func to a map of functions,
// sName argument is a function identifier - a key of this map.
// You may need to call this function in case of using HWGui's xml forms.
func RegFunc(sName string, fu func([]string) string) {

	if mfu == nil {
		mfu = make(map[string]func([]string) string)
	}
	mfu[sName] = fu
}

func wait() {
	for !bEndProg {
		time.Sleep(20 * time.Millisecond)
	}
}

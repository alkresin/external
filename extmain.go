/*
 *  External - GUI framework for Go language
 *  External is a Go library to build GUI application, using external
 *       standalone GUI server application: https://github.com/alkresin/guiserver
 *
 * Copyright 2018 Alexander S.Kresin <alex@kresin.ru>
 * www - http://www.kresin.ru
 */
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

var sLogName = "egui.log"
var bEndProg = false

var connOut, connIn net.Conn
var bConnExist = false

func Init(sOpt string) bool {

	var err error

	iPort := 3101
	sServer := "guiserver.exe"
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
				sLog = "-log+"
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
			WriteLog(sLogName, fmt.Sprintln(sServer, sIp, iPort))
			WriteLog(sLogName, fmt.Sprintln(err))
			return false
		}
	}
	_, err = connOut.Read(buf)
	if err != nil {
		WriteLog(sLogName, fmt.Sprintln(err))
		connOut.Close()
		return false
	}

	connIn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", sIp, iPort+1))
	if err != nil {
		WriteLog(sLogName, fmt.Sprintln(sServer, sIp, iPort+1))
		WriteLog(sLogName, fmt.Sprintln(err))
		return false
	}
	_, err = connIn.Read(buf)
	if err != nil {
		WriteLog(sLogName, fmt.Sprintln(err))
		connIn.Close()
		return false
	}

	bConnExist = true
	go listen(iPort + 1)
	time.Sleep(100 * time.Millisecond)

	return true

}

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
		length, err := connIn.Read(buffer)
		if err != nil {
			WriteLog(sLogName, "Read error\r\n")
			return
		}

		if buffer[0] != byte('+') || buffer[length-1] != byte('\n') {
			bErr = true
		}

		if !bErr {
			err = json.Unmarshal(buffer[1:length-1], &arr)
			if err != nil {
				WriteLog(sLogName, "Unmarshal error\r\n")
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
								WriteLog(sLogName, fmt.Sprintf("runproc param Unmarshal error (%s)\r\n", arr[2]))
							}
						}
						//WriteLog(sLogName, fmt.Sprintf("pgo> (%s) len:%d\r\n",arr[2],len(ap) ))
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
								WriteLog(sLogName, fmt.Sprintf("runproc param Unmarshal error (%s)\r\n", arr[2]))
							}
						}
						//WriteLog(sLogName, fmt.Sprintf("pgo> (%s) len:%d\r\n",arr[2],len(ap) ))
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
					oW := GetWnd(arr[1])
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
			WriteLog(sLogName, fmt.Sprintf("Wrong message: %s]\r\n", string(buffer[:length])))
		}
	}
}

func sendResponse(conn net.Conn, s string) {
	conn.Write([]byte("+" + s + "\n"))
}

func Sendout(s string) bool {

	var err error
	buf := make([]byte, 128)

	_, err = connOut.Write([]byte("+" + s + "\n"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = connOut.Read(buf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SendoutAndReturn(s string, length int) []byte {

	var err error
	buf := make([]byte, length)

	_, err = connOut.Write([]byte("+" + s + "\n"))
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	length, err = connOut.Read(buf)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	return buf[:length-1]
}

func WriteLog(sName string, sText string) {

	f, err := os.OpenFile(sName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(sText)

}

func RegFunc(sName string, fu func([]string) string) {

	if mfu == nil {
		mfu = make(map[string]func([]string) string)
	}
	mfu[sName] = fu
}

func Wait() {
	for !bEndProg {
		time.Sleep(20 * time.Millisecond)
	}
}

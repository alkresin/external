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
	"sync"
	"time"
)

const (
	VerProto = "1.1"
	Version  = "1.1"
	FileRoot = "gs"
)

type ConnEx struct {
	iType int8
	iPort int
	sIp   string
	sFileName string
	conn  net.Conn
	f    *os.File
}

var sLogName = "egui.log"
var bEndProg = false
var bWait = false

//var connOut, connIn net.Conn
var bConnExist = false

var pConnOut, pConnIn *ConnEx
var sDir = ""
var sFileRoot = FileRoot

var bPacket = false
var sPacketBuf string

var aRunProc [][]string
var muxRunProc sync.Mutex

var aRunFu []func()
var muxRunFu sync.Mutex

// Init runs, if needed, the Guiserver application, and connects to it.
// It returns 0, if the connection is successful, 1 - in other case,
// 2 -if a protocol version of a GuiServer isn't equal to a local protocol.
// The sOpt argument specifies connection details. It may contain following strings:
// guiserver=<full path to GuiServer executable>
// address=<ip address of a computer, where GuiServer runs>
// port=<tcp/ip port number>
// log=<0, 1 or 2> - logging level
func Init(sOpt string) int {

	var err error

	iConnType := 1
	iPort := 3101
	sServer := "guiserver"
	sIp := "127.0.0.1"
	sFileName := ""
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
			} else if len(s) > 5 && s[:4] == "type" {
				iConnType, _ = strconv.Atoi(strings.TrimSpace(s[5:]))
			} else if len(s) > 5 && s[:3] == "dir" {
				sDir = strings.TrimSpace(s[4:])
			} else if len(s) > 5 && s[:4] == "file" {
				sFileRoot = strings.TrimSpace(s[5:])
				if strings.Contains(sFileRoot, "*") {
					sz := strconv.Itoa(int((time.Now().UnixNano() / 1000000) % 100000))
					sFileRoot = strings.Replace(sFileRoot, "*", sz, -1)
				}
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

	if iConnType == 2 {
		if sDir == "" {
			sDir = os.TempDir()
		}
		sFileName = sDir + string(os.PathSeparator) + sFileRoot
		if sServer != "" {
			os.Remove(sFileName)
		}
	}

	if sServer != "" {
		var cmd *exec.Cmd
		if iConnType == 2 && (sDir != "" || sFileRoot != FileRoot) {
			cmd = exec.Command(sServer, fmt.Sprintf("-p%d", iPort), "-t2", sLog, "-d"+sDir, "-f"+sFileRoot )
		} else {
			cmd = exec.Command(sServer, fmt.Sprintf("-p%d", iPort), fmt.Sprintf("-t%d", iConnType), sLog)
		}
		cmd.Start()
	}
	time.Sleep(100 * time.Millisecond)

	pConnOut = &ConnEx{ iType: int8(iConnType), iPort: iPort, sIp: sIp, sFileName: sFileName+".gs1" }
	pConnIn = &ConnEx{ iType: int8(iConnType), iPort: iPort+1, sIp: sIp, sFileName: sFileName+".gs2" }

	if !pConnOut.Connect() {
		return 1
	}

	bConnExist = true
	
	iBufLen, err := pConnOut.Read(&buf)

	if err != nil {
		WriteLog(fmt.Sprintln(err))
		pConnOut.Close()
		bConnExist = false
		return 1
	}
	sVer := string(buf[:iBufLen-1])
	sVer = sVer[(strings.Index(sVer, "/") + 1):]

	if !pConnIn.Connect() {
		return 1
	}

	_, err = pConnIn.Read(&buf)

	if err != nil {
		WriteLog(fmt.Sprintln(err))
		Exit()
		return 1
	}

	if sVer != VerProto {
		WriteLog("\r\nProtocol version mismatched. Need " + VerProto + ", received " + sVer)
		Exit()
		return 2
	}

	if iConnType == 2 {
		pConnIn.Write( "[\"Ok\"]" )
	}

	go listen(iPort + 1)
	time.Sleep(100 * time.Millisecond)

	return 0
}

// Exit closes the connection to Guiserver.
func Exit() {
	if bConnExist {
		bConnExist = false
		pConnOut.Write("+[\"exit\"]\n")
		time.Sleep(10 * time.Millisecond)
		pConnOut.Close()
		pConnIn.Close()
	}
}

func listen(iPort int) {

	var bErr bool

	buffer := make([]byte, 1024)
	arr := make([]string, 5)
	for {

		bErr = false
		length, err := pConnIn.Read(&buffer)

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
				//sendResponse(connIn, "[\"Ok\"]")
				pConnIn.Write( "[\"Ok\"]" )
				if len(arr) > 1 {
					if bWait {
						tmp := make([]string, len(arr))
						muxRunProc.Lock()
						copy(tmp, arr)
						aRunProc = append(aRunProc, tmp)
						muxRunProc.Unlock()
					} else {
						runproc(arr)
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
						//sendResponse(connIn, string(b))
						pConnIn.Write("+" + string(b) + "\n")
					} else {
						//sendResponse(connIn, "[\"Err\"]")
						pConnIn.Write("+[\"Err\"]" + "\n")
					}
				} else {
					bErr = true
					//sendResponse(connIn, "[\"Err\"]")
					pConnIn.Write("[\"Err\"]")
				}
			case "exit":
				//sendResponse(connIn, "[\"Ok\"]")
				pConnIn.Write("[\"Ok\"]")
				if len(arr) > 1 {
					oW := Wnd(arr[1])
					if oW != nil {
						oW.delete()
					}
				} else {
					bErr = true
				}
			case "endapp":
				//sendResponse(connIn, "[\"Goodbye\"]")
				pConnIn.Write("[\"Goodbye\"]")
				time.Sleep(100 * time.Millisecond)
				bEndProg = true
				//connIn.Close()
				//pConnIn.Close()
				Exit()
				//WriteLog("The End")
				return
			default:
				//sendResponse(connIn, "[\"Error\"]")
				pConnIn.Write("[\"Error\"]")
				bErr = true
			}
		}
		if bErr {
			WriteLog(fmt.Sprintf("Wrong message: %s]\r\n", string(buffer[:length])))
		}
	}
}

/*
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
*/

func sendout(s string) bool {

	var err error

	if bPacket {
		sPacketBuf += "," + s
	} else {
		if !bConnExist {
			WriteLog("sendout: No connection established.\r\n")
			return false
		}

		//_, err = connOut.Write([]byte("+" + s + "\n"))
		_, err = pConnOut.Write("+" + s + "\n")
		if err != nil {
			fmt.Println(err)
			return false
		}

		buf := make([]byte, 128)
		//_, err = connOut.Read(buf)
		_, err = pConnOut.Read(&buf)
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
		WriteLog("sendoutAndReturn: No connection established.\r\n")
		return []byte("")
	}

	//_, err = connOut.Write([]byte("+" + s + "\n"))
	_, err = pConnOut.Write("+" + s + "\n")
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	//length, err = connOut.Read(buf)
	//length, err := read(connOut, &buf)
	length, err := pConnOut.Read(&buf)
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

// AddFuncToIdle adds a function to be executed while wait state.
func AddFuncToIdle(fu func()) {
	muxRunFu.Lock()
	aRunFu = append(aRunFu, fu)
	muxRunFu.Unlock()
}

func runproc(arr []string) {
	if fnc, bExist := mfu[arr[1]]; bExist {
		var ap []string
		if len(arr) > 2 {
			ap = make([]string, 5)
			err := json.Unmarshal([]byte(arr[2]), &ap)
			if err != nil {
				WriteLog(fmt.Sprintf("runproc param Unmarshal error (%s)\r\n", arr[2]))
			}
		}
		//WriteLog(fmt.Sprintf("pgo> (%s) len:%d\r\n",arr[2],len(ap) ))
		fnc(ap)
	}
}

func Wait() {
	bWait = true
	for !bEndProg {
		for {
			muxRunFu.Lock()
			if len(aRunFu) == 0 {
				muxRunFu.Unlock()
				break
			}
			fu := aRunFu[0]
			aRunFu = append(aRunFu[:0], aRunFu[1:]...)
			muxRunFu.Unlock()

			fu()
		}
		for {
			muxRunProc.Lock()
			if len(aRunProc) == 0 {
				muxRunProc.Unlock()
				break
			}
			arr := aRunProc[0]
			aRunProc = append(aRunProc[:0], aRunProc[1:]...)
			muxRunProc.Unlock()

			runproc(arr)
		}
		time.Sleep(20 * time.Millisecond)
	}
	bWait = false
}

func (p *ConnEx) Connect() bool {

	var err error

	if p.iType == 1 {
		p.conn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", p.sIp, p.iPort))
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			p.conn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", p.sIp, p.iPort))
			if err != nil {
				time.Sleep(3000 * time.Millisecond)
				p.conn, err = net.Dial("tcp4", fmt.Sprintf("%s:%d", p.sIp, p.iPort))
				if err != nil {
					WriteLog(fmt.Sprintln(p.sIp, p.iPort))
					WriteLog(fmt.Sprintln(err))
					return false
				}
			}
		}
	} else if p.iType == 2 {
		//os.Remove(p.sFileName)
		p.f,err = os.OpenFile(p.sFileName, os.O_RDWR, 0644)
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			p.f,err = os.OpenFile(p.sFileName, os.O_RDWR, 0644)
			if err != nil {
				time.Sleep(3000 * time.Millisecond)
				p.f,err = os.OpenFile(p.sFileName, os.O_RDWR, 0644)
				if err != nil {
					WriteLog(p.sFileName)
					WriteLog(fmt.Sprintln(err))
					return false
				}
			}
		}
	}

	return true
}

func (p *ConnEx) Close() {

	if p.iType == 1 {
		p.conn.Close()
	} else if p.iType == 2 {
		p.f.Close()
	}
}

func (p *ConnEx) Read(pBuff *[]byte) (int, error) {

	*pBuff = (*pBuff)[:0]
	tmp := make([]byte, 256)
	if p.iType == 1 {
		for {
			length, err := p.conn.Read(tmp)
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
	} else if p.iType == 2 {
		b1 := make([]byte, 1)
		for bConnExist {
			p.f.Seek(0,0)
			_, err := p.f.Read(b1)
			if err != nil {
				WriteLog("Read error\r\n")
				return 0, err
			}
			if b1[0] == 2 {
				for ;; {
					n, err := p.f.Read(tmp)
					if err != nil {
						WriteLog("Read error\r\n")
						return 0, err
					}
					iPos := -1
					for i := 0; i < n; i++ {
						if tmp[i] == 10 {
							iPos = i
							break
						}
					}
					if iPos >= 0 {
						*pBuff = append(*pBuff, tmp[:iPos+1]...)
						break
					} else if n < len(tmp) {
						*pBuff = append(*pBuff, tmp[:n+1]...)
						break
					} else {
						*pBuff = append(*pBuff, tmp...)
					}
				}
				break
			} else {
				time.Sleep(2 * time.Millisecond)
			}
		}
		//WriteLog(string(*pBuff))
		return len(*pBuff), nil
	}

   return 0, nil
}

func (p *ConnEx) Write( s string ) (int, error) {

	var err error

	if p.iType == 1 {
		_, err = p.conn.Write( []byte(s) )
	} else if p.iType == 2 {
		p.f.WriteAt( []byte(s), 1 )
		p.f.WriteAt( []byte{1}, 0 )
	}

	return 0, err
}

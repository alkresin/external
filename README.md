# external

[![GoDoc](https://godoc.org/github.com/alkresin/external?status.svg)](https://godoc.org/github.com/alkresin/external)

External is a GUI library for Go (Golang), based on connection to external GUI server application.
The connection can be esstablished via tcp/ip sockets or via regular files.
To use it you need to have the GuiServer executable, which may be compiled from sources, hosted in https://github.com/alkresin/guiserver, or downloaded from http://www.kresin.ru/en/guisrv.html.
Join the multilanguage group https://groups.google.com/d/forum/guiserver to discuss the GuiServer, External and related issues.


To get rid of a console window, *use -ldflags "-H windowsgui"* option in *go build* statement for your application.

--------------------
Alexander S.Kresin
http://www.kresin.ru/
mailto: alkresin@yahoo.com

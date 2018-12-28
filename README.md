# external

[![GoDoc](https://godoc.org/github.com/alkresin/external?status.svg)](https://godoc.org/github.com/alkresin/external)
[![Build Status](https://travis-ci.org/alkresin/external.svg?branch=master)](https://travis-ci.org/alkresin/external)
[![Codecov](https://codecov.io/gh/alkresin/external/branch/master/graph/badge.svg)](https://codecov.io/gh/alkresin/external)
[![Go Report
Card](https://goreportcard.com/badge/github.com/alkresin/external)](https://goreportcard.com/report/github.com/alkresin/external)

External is a GUI library for Go (Golang), based on tcp/ip connection to external GUI server application.
To use it you need to have the GuiServer executable, which may be compiled from sources, hosted in https://github.com/alkresin/guiserver, or downloaded from http://www.kresin.ru/en/guisrv.html.
Join the multilanguage group https://groups.google.com/d/forum/guiserver to discuss the GuiServer, External and related issues.


To get rid of a console window, *use -ldflags "-H windowsgui"* option in *go build* statement.

--------------------
Alexander S.Kresin
http://www.kresin.ru/
mailto: alkresin@yahoo.com

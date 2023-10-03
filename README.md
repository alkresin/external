# external

[![GoDoc](https://godoc.org/github.com/alkresin/external?status.svg)](https://godoc.org/github.com/alkresin/external)

<b> Attention! Since October 6, 2023 we are forced to use two-factor authentication to be able to
   update the repository. Because it's not suitable for me, I will probably use another place for projects.
   Maybe, https://gitflic.ru/, maybe, Sourceforge... Follow the news on my website, http://www.kresin.ru/

   Внимание! С 6 октября 2023 года нас вынуждают использовать двухфакторную идентификацию для того, чтобы
   продолжать работать над проектами. Поскольку для меня это крайне неудобно, я, возможно, переведу проекты
   на другое место. Это может быть https://gitflic.ru/, Sourceforge, или что-то еще. Следите за новостями
   на моем сайте http://www.kresin.ru/ </b>

External is a GUI library for Go (Golang), based on connection to external GUI server application.
The connection can be esstablished via tcp/ip sockets or via regular files.
To use it you need to have the GuiServer executable, which may be compiled from sources, hosted in https://github.com/alkresin/guiserver, or downloaded from http://www.kresin.ru/en/guisrv.html.
Join the multilanguage group https://groups.google.com/d/forum/guiserver to discuss the GuiServer, External and related issues.


To get rid of a console window, *use -ldflags "-H windowsgui"* option in *go build* statement for your application.

--------------------
Alexander S.Kresin
http://www.kresin.ru/
mailto: alkresin@yahoo.com

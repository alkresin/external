# external

[![GoDoc](https://godoc.org/github.com/alkresin/external?status.svg)](https://godoc.org/github.com/alkresin/external)

<b> Attention! Since October 6, 2023, we have been forced to use two-factor identification in order to 
   log in to github.com under your account. I can still do <i>git push</i> from the command line, but I can't
   use other services, for example, to answer questions. That's why I'm opening new projects on 
   https://gitflic.ru /, Sourceforge, or somewhere else. Follow the news on my website http://www.kresin.ru/

   Внимание! С 6 октября 2023 года нас вынуждили использовать двухфакторную идентификацию для того, чтобы 
   входить на github.com под своим аккаунтом. Я пока могу делать <i>git push<i> из командной строки, но не могу
   использовать другие сервисы, например, отвечать на вопросы. Поэтому новые проекты я открываю на 
   https://gitflic.ru/, Sourceforge, или где-то еще. Следите за новостями на моем сайте http://www.kresin.ru/ </b>

External is a GUI library for Go (Golang), based on connection to external GUI server application.
The connection can be esstablished via tcp/ip sockets or via regular files.
To use it you need to have the GuiServer executable, which may be compiled from sources, hosted in https://github.com/alkresin/guiserver, or downloaded from http://www.kresin.ru/en/guisrv.html.
Join the multilanguage group https://groups.google.com/d/forum/guiserver to discuss the GuiServer, External and related issues.


To get rid of a console window, *use -ldflags "-H windowsgui"* option in *go build* statement for your application.

--------------------
Alexander S.Kresin
http://www.kresin.ru/
mailto: alkresin@yahoo.com

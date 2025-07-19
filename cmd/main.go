package main

//version 1.1.0

import (
	"webtest/colors"
	"webtest/intro"
	"webtest/server"
)

func init() {
	colors.SetColor(colors.Text_Red)
	intro.IntroLog()
	colors.ResetColor()
}

func main() {
	server.StartServer()
}

package main

import (
	"runtime"
	"strconv"

	blt "bearlibterminal"
)

const (
	//for setting blt window
	WindowSizeX = 30
	WindowSizeY = 30
	GameTitle   = "unnamed game"
	FontName    = "UbuntuMono-R.ttf"
	FontSize    = 18
)

func constrainThreads() {
	/*Constraining processor and threads is necessary,
	  because BearLibTerminal often crashes otherwise*/
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func InitializeBLT() {
	/*Constraining threads and setting blt window*/
	constrainThreads()
	blt.Open()
	sizeX, sizeY := strconv.Itoa(WindowSizeX), strconv.Itoa(WindowSizeY+2)
	sizeFont := strconv.Itoa(FontSize)
	window := "window: size=" + sizeX + "x" + sizeY
	blt.Set(window + ", title=' " + GameTitle + "'; font: " + FontName + ", size=" + sizeFont)
	blt.Clear()
	blt.Refresh()
}

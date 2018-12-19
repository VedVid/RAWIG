/*
Copyright (c) 2018 Tomasz "VedVid" Nowakowski

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.
*/

package main

import (
	"runtime"
	"strconv"

	blt "bearlibterminal"
)

const (
	// Setting BearLibTerminal window.
	WindowSizeX = 50
	WindowSizeY = 25
	MapSizeX = 30
	MapSizeY = 20
	UIPosX = MapSizeX
	UIPosY = 0
	UISizeX = WindowSizeX - MapSizeX
	UISizeY = WindowSizeY
	LogSizeX = MapSizeX
	LogSizeY = WindowSizeY - MapSizeY
	LogPosX = 0
	LogPosY = MapSizeY
	GameTitle   = "unnamed game"
	FontName    = "UbuntuMono-R.ttf"
	FontSize    = 18
)

func constrainThreads() {
	/* Constraining processor and threads is necessary,
	   because BearLibTerminal often crashes otherwise. */
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func InitializeBLT() {
	/* Constraining threads and setting BearLibTerminal window. */
	constrainThreads()
	blt.Open()
	sizeX, sizeY := strconv.Itoa(WindowSizeX), strconv.Itoa(WindowSizeY)
	sizeFont := strconv.Itoa(FontSize)
	window := "window: size=" + sizeX + "x" + sizeY
	blt.Set(window + ", title=' " + GameTitle + "'; font: " + FontName + ", size=" + sizeFont)
	blt.Clear()
	blt.Refresh()
}

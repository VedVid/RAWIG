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
	blt "bearlibterminal"
	"errors"
	"fmt"
	"unicode/utf8"
)

const (
	MaxMessageBuffer = WindowSizeY - MapSizeY
)

func PrintMenu(x, y int, header string, options []string) {
	/* Function PrintMenu takes four arguments: two ints that are
	   top-left corner of menu, header, and slice of options.
	   If header is empty, text is moved one tile higher to
	   avoid wasting space.
	   During execution, it joins header and all of options in
	   one text, with additional formatting.
	   For example, header "MyMenu" and options ["first", "two"]
	   would produce that kind of output:
	       MyMenu
	       a) first
	       b) two
	    It refreshed terminal and waits for player input at the end. */
	blt.ClearArea(UIPosX, UIPosY, UISizeX, UISizeY)
	if header == "" {
		y--
	}
	txt := header
	for i, v := range options {
		txt = txt + "\n" + OrderToCharacter(i) + ") " + v
	}
	txt = txt + "\n[[ESC]] back"
	blt.Print(x, y, txt)
	blt.Refresh()
}

func PrintInventoryMenu(x, y int, header string, options Objects) {
	/* PrintInventoryMenu is helper function that takes Objects
	   as its main argument, and adds their names (currently
	   their symbol representation, due to some strange decisions
	   made by dev, objects doesn't have names yet) to the opts
	   slice of strings, then calls PrintMenu using that list.
	   Unfortunately that kind of "hack" is necessary, because
	   Go doesn't support generics and optional arguments,
	   and still doesn't provide sensible alternatives.
	   I'd like to just pass Objects to the PrintMenu func. */
	var opts = []string{}
	for _, v := range options {
		opts = append(opts, v.Name)
	}
	PrintMenu(x, y, header, opts)
}

func PrintEquipmentMenu(x, y int, header string, options Objects) {
	/* Similar to PrintInventoryMenu, but it sorts options
	   by their Slots initially, and slot in showed before
	   item name.
	   Note that it shows Creature's slots,
	   not all equippable objects in inventory.
	   Because of this, it is necessary to find "true" length
	   of options, skipping all nil pointers.
	   Unfortunately, it may crash in future, with
	   more slots involved. */
	var opts = []string{}
	for i := 0; i < len(options); i++ {
		txt := ""
		if options[i] != nil {
			txt = "[[" + SlotStrings[i] + "]] " + options[i].Name
		} else {
			txt = "[[" + SlotStrings[i] + "]] empty"
		}
		opts = append(opts, txt)
	}
	PrintMenu(x, y, header, opts)
}

func PrintEquippables(x, y int, header string, options Objects) {
	/* PrintEquippables is function that prints list of equippables. */
	var opts = []string{}
	for _, v := range options {
		opts = append(opts, v.Name)
	}
	PrintMenu(x, y, header, opts)
}

func PrintMessages(x, y int, header string) {
	/* PrintMessages works as PrintMenu, but it
	   will not format text in special way. */
	if header == "" {
		y--
	}
	txt := header
	for _, v := range MsgBuf {
		txt = txt + "\n" + v
	}
	blt.Print(x, y, txt)
}

func AddMessage(message string) {
	/* AddMessage is function that adds message
	   to the MessageBuffer. It removes the oldest
	   line to keep size set in MaxMessageBuffer.
	   But first, it checks if passed message is
	   not longer than whole message log.
	   This is mostly harmless, so AddMessage
	   does not returns error, but prints it
	   at its own. */
	var err error
	messageLen := utf8.RuneCountInString(message)
	if messageLen > LogSizeX {
		txt := MessageLengthError(message, messageLen, LogSizeX)
		err = errors.New("Message is too long to fit message log. " + txt)
		fmt.Println(err)
	}
	if len(MsgBuf) < MaxMessageBuffer {
		MsgBuf = append(MsgBuf, message)
	} else {
		MsgBuf = append(MsgBuf[1:], message)
	}
	PrintLog()
	blt.Refresh()
}

func RemoveFirstMessage() {
	MsgBuf = MsgBuf[1:]
	blt.Layer(UILayer)
	blt.ClearArea(LogPosX, LogPosY, LogPosX+LogSizeX, LogPosY+LogSizeY)
	PrintLog()
	blt.Refresh()
}

func RemoveLastMessage() {
	/* Function RemoveLastMessage is called when it is necessary to remove
	   last message from buffer, even if said buffer is not full.
	   It removes last message, clears its area, and reprints log. */
	MsgBuf = MsgBuf[:len(MsgBuf)-1]
	blt.Layer(UILayer)
	blt.ClearArea(LogPosX, LogPosY, LogPosX+LogSizeX, LogPosY+LogSizeY)
	PrintLog()
	blt.Refresh()
}

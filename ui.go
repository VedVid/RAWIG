/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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

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
)

const (
	MaxMessageBuffer = 10
)

func PrintMenu(x, y int, header string, options []string) {
	/* Function PrintMenu takes four arguments: two ints that are
	   top-left corner of menu, header, and slice of options.
	   During execution, it joins header and all of options in
	   one text, with additional formatting.
	   For example, header "MyMenu" and options ["first", "two"]
	   would produce that kind of output:
	       MyMenu
	       a) first
	       b) two
	    */
	txt := header
	for i, v := range options {
		txt = txt + "\n" + OrderToCharacter(i) + ") " + v
	}
	blt.Print(x, y, txt)
}

func PrintMessages(x, y int, header string) {
	txt := header
	for _, v := range MsgBuf {
		txt = txt + "\n" + v
	}
	blt.Print(x, y, txt)
}

func AddMessage(message string) {
	if len(MsgBuf) < MaxMessageBuffer {
		MsgBuf = append(MsgBuf, message)
	} else {
		MsgBuf = append(MsgBuf[1:], message)
	}
}

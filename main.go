package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const CTRL_KEY = 0x1f

type editorConfig struct {
	screenRows    int
	screenColumns int
}

var editor editorConfig

func initEditor() {
	getWindowsSize(&editor.screenRows, &editor.screenColumns)
}
func getWindowsSize(rows *int, columns *int) {
	tempColumns, tempRows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	*rows = tempRows
	*columns = tempColumns
}
func editorDrawRows() {
	for x := 0; x < editor.screenRows; x++ {
		os.Stdout.Write([]byte("~\r\n"))
	}
}
func die(s string) {
	clearScreenSequence := "\x1b[2J"
	repositionCursorSequence := "\x1b[H"
	_, _ = os.Stdout.Write([]byte(clearScreenSequence))
	_, _ = os.Stdout.Write([]byte(repositionCursorSequence))
	fmt.Printf("%s", s)
	os.Exit(1)
}

func editorReadKey() rune {
	var buf [1]byte
	file := os.Stdin

	for {
		nread, err := file.Read(buf[:])
		if err != nil {
			fmt.Println("Error reading from standard input\r")
			break
		}
		if nread == -1 {
			die("read")
		}

	}
	return rune(buf[0])
}

func editorProcessKeypress() {
	c := editorReadKey()

	switch c {
	case CTRL_KEY & 'q':
		clearScreenSequence := "\x1b[2J"
		repositionCursorSequence := "\x1b[H"
		_, _ = os.Stdout.Write([]byte(clearScreenSequence))
		_, _ = os.Stdout.Write([]byte(repositionCursorSequence))
		os.Exit(0)
		break
	}
}
func editorRefreshScreen() {
	clearScreenSequence := "\x1b[2J"
	repositionCursorSequence := "\x1b[H"
	_, _ = os.Stdout.Write([]byte(clearScreenSequence))
	_, _ = os.Stdout.Write([]byte(repositionCursorSequence))
	editorDrawRows()
	_, _ = os.Stdout.Write([]byte(repositionCursorSequence))
}
func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	defer term.Restore(int(os.Stdin.Fd()), oldState)
	initEditor()
	for {
		editorRefreshScreen()
		editorProcessKeypress()
	}
}

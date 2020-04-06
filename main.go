package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

func main() {

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	s, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.Clear()

	stack := []string{}

	runes := []rune{}

	for {
		w, h := s.Size()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
					runes = append([]rune{ev.Rune()}, runes...)

					for i, r := range runes {
						s.SetContent(w - 1 - i, h - 1, r, nil, tcell.StyleDefault)
					}

					s.Show()
				}

			case tcell.KeyBackspace2:
				if len(runes) > 0 {
					s.SetContent(w - len(runes), h - 1, 0, nil, tcell.StyleDefault)

					runes = runes[1:]

					for i, r := range runes {
						s.SetContent(w - 1 - i, h - 1, r, nil, tcell.StyleDefault)
					}

					s.Show()
				}

			case tcell.KeyEnter:
				if len(runes) > 0 {
					stack = append([]string{string(runes)}, stack...)

					runes = runes[:0]

					s.Clear()

					for l, str := range stack {
						for i, r := range []rune(str) {
							s.SetContent(w - 1 - i, h - 2 - l, r, nil, tcell.StyleDefault)
						}
					}

					s.Show()
				}

			case tcell.KeyCtrlD:
				if len(stack) > 0 {
					stack = stack[1:]

					s.Clear()

					for l, str := range stack {
						for i, r := range []rune(str) {
							s.SetContent(w - 1 - i, h - 2 - l, r, nil, tcell.StyleDefault)
						}
					}

					s.Show()
				}

			case tcell.KeyCtrlU:
				if len(stack) > 0 {
					stack = append([]string{stack[0]}, stack...)

					s.Clear()

					for l, str := range stack {
						for i, r := range []rune(str) {
							s.SetContent(w - 1 - i, h - 2 - l, r, nil, tcell.StyleDefault)
						}
					}

					s.Show()
				}

			case tcell.KeyCtrlL:
				s.Sync()

			case tcell.KeyEscape:
				s.Fini()

				os.Exit(0)
			}

		case *tcell.EventResize:
			s.Sync()
		}
	}
}

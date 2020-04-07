package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
)

var sc tcell.Screen

func reverse(s string) string {
	runes := []rune{}

	for _, r := range s {
		runes = append([]rune{r}, runes...)
	}

	return string(runes)
}

func setInputContent(input []rune, w, h int) {
	for i, r := range input {
		sc.SetContent(w-1-i, h-1, r, nil, tcell.StyleDefault)
	}
}

func setStackContent(stack []string, w, h int) {
	for l, str := range stack {
		for i, r := range []rune(str) {
			sc.SetContent(w-1-i, h-2-l, r, nil, tcell.StyleDefault)
		}
	}
}

func init() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)

	var err error

	sc, err = tcell.NewScreen()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err = sc.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	sc.Clear()
}

func main() {
	stack := []string{}

	input := []rune{}

	for {
		w, h := sc.Size()

		ev := sc.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
					input = append([]rune{ev.Rune()}, input...)

					setInputContent(input, w, h)

					sc.Show()

				case '+', '-', '*', '/':
					if len(input) > 0 && len(stack) > 0 {
						lhs, _ := strconv.ParseFloat(reverse(string(input)), 32)
						rhs, _ := strconv.ParseFloat(reverse(stack[0]), 32)

						stack[0] = reverse(fmt.Sprintf("%v", lhs+rhs))

						input = input[:0]

						sc.Clear()

						setStackContent(stack, w, h)

						sc.Show()
					}
				}

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					sc.SetContent(w-len(input), h-1, 0, nil, tcell.StyleDefault)

					input = input[1:]

					setInputContent(input, w, h)

					sc.Show()
				}

			case tcell.KeyEnter:
				if len(input) > 0 {
					stack = append([]string{string(input)}, stack...)

					input = input[:0]

					sc.Clear()

					setStackContent(stack, w, h)

					sc.Show()
				}

			case tcell.KeyCtrlD:
				if len(stack) > 0 {
					stack = stack[1:]

					sc.Clear()

					setStackContent(stack, w, h)
					setInputContent(input, w, h)

					sc.Show()
				}

			case tcell.KeyCtrlU:
				if len(stack) > 0 {
					stack = append([]string{stack[0]}, stack...)

					sc.Clear()

					setStackContent(stack, w, h)
					setInputContent(input, w, h)

					sc.Show()
				}

			case tcell.KeyCtrlS:
				if len(stack) > 1 {
					tmp := stack[0]
					stack[0] = stack[1]
					stack[1] = tmp

					sc.Clear()

					setStackContent(stack, w, h)
					setInputContent(input, w, h)

					sc.Show()
				}

			case tcell.KeyCtrlR:
				if len(stack) > 1 {
					tmp := stack[0]
					stack = append(stack[1:len(stack)], tmp)

					sc.Clear()

					setStackContent(stack, w, h)
					setInputContent(input, w, h)

					sc.Show()
				}

			case tcell.KeyCtrlL:
				sc.Sync()

			case tcell.KeyEscape:
				sc.Fini()

				os.Exit(0)
			}

		case *tcell.EventResize:
			sc.Sync()
		}
	}
}

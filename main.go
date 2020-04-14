package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
)

var (
	sc    tcell.Screen
	st    tcell.Style
	stack []float64
	input []rune
)

func screen() {
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

	st = tcell.StyleDefault
	stack = []float64{}
	input = []rune{}
}

func draw(help bool) {
	w, h := sc.Size()

	sc.Clear()

	sc.SetContent(0, h-1, ':', nil, st)

	for c, r := range input {
		sc.SetContent(1+c, h-1, r, nil, st)
	}

	sc.ShowCursor(1+len(input), h-1)

	for l, n := range stack {
		for c, r := range []rune(fmt.Sprintf("%v: %v", l+1, n)) {
			sc.SetContent(c, h-2-l, r, nil, st)
		}
	}

	if help {
		keys := []string{
			"ADD :[+]",
			"SUB :[-]",
			"MUL :[*]",
			"DIV :[/]",
			"+/- :[i]",
			"DUP :[u]",
			"ROT :[r]",
			"POW :[p]",
			"SQRT :[t]",
			"DROP :[d]",
			"SWAP :[s]",
			"HELP :[h]",
			"QUIT :[q]",
		}

		for l, str := range keys {
			for c, r := range str {
				sc.SetContent(w-1-len(str)+c, h-1-len(keys)+l, r, nil, st)
			}
		}
	}

	sc.Show()
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s\nPress [h] for help\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	screen()

	help := false

	for {
		ev := sc.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				switch ev.Rune() {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
					input = append(input, ev.Rune())

				case '+', '-', '*', '/', 'p':
					if len(input) > 0 && len(stack) > 0 {
						if n, err := strconv.ParseFloat(string(input), 32); err == nil {
							switch ev.Rune() {
							case '+':
								stack[0] += n

							case '-':
								stack[0] -= n

							case '*':
								stack[0] *= n

							case '/':
								stack[0] /= n

							case 'p':
								stack[0] = math.Pow(stack[0], n)
							}

							input = input[:0]
						}
					} else if len(stack) > 1 {
						rhs, lhs := stack[0], stack[1]
						stack = stack[2:]

						switch ev.Rune() {
						case '+':
							stack = append([]float64{lhs + rhs}, stack...)

						case '-':
							stack = append([]float64{lhs - rhs}, stack...)

						case '*':
							stack = append([]float64{lhs * rhs}, stack...)

						case '/':
							stack = append([]float64{lhs / rhs}, stack...)

						case 'p':
							stack = append([]float64{math.Pow(lhs, rhs)}, stack...)
						}
					}

				case 'i', 't':
					if len(input) > 0 {
						switch ev.Rune() {
						case 'i':
							if input[0] == '-' {
								input = input[1:]
							} else {
								input = append([]rune{'-'}, input...)
							}

						case 't':
							if n, err := strconv.ParseFloat(string(input), 32); err == nil {
								input = []rune(fmt.Sprintf("%v", math.Sqrt(n)))
							}
						}
					} else if len(stack) > 0 {
						switch ev.Rune() {
						case 'i':
							stack[0] *= -1

						case 't':
							stack[0] = math.Sqrt(stack[0])
						}
					}

				case 'd', 'u':
					if len(stack) > 0 {
						switch ev.Rune() {
						case 'd':
							stack = stack[1:]

						case 'u':
							stack = append([]float64{stack[0]}, stack...)
						}
					}

				case 's', 'r':
					if len(stack) > 1 {
						switch ev.Rune() {
						case 's':
							tmp := stack[0]
							stack[0] = stack[1]
							stack[1] = tmp

						case 'r':
							tmp := stack[0]
							stack = append(stack[1:len(stack)], tmp)
						}
					}

				case 'h':
					help = !help

				case 'q':
					sc.Fini()
					os.Exit(0)
				}

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
				}

			case tcell.KeyEnter:
				if len(input) > 0 {
					if n, err := strconv.ParseFloat(string(input), 32); err == nil {
						stack = append([]float64{n}, stack...)
						input = input[:0]
					}
				}

			case tcell.KeyCtrlL:
				sc.Sync()
			}

			draw(help)

		case *tcell.EventResize:
			sc.Sync()
			draw(help)
		}
	}
}

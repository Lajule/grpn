package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/gdamore/tcell"
)

var (
	sc tcell.Screen
	st tcell.Style
)

func reverse(s string) string {
	runes := []rune{}

	for _, r := range s {
		runes = append([]rune{r}, runes...)
	}

	return string(runes)
}

func setInputContent(input []rune, w, h int) {
	for c, r := range input {
		sc.SetContent(w-1-c, h-1, r, nil, st)
	}
}

func clearInput(l, w, h int) {
	for c := 0; c < l; c++ {
		sc.SetContent(w-1-c, h-1, 0, nil, st)
	}
}

func setStackContent(stack []float64, w, h int) {
	for l, n := range stack {
		for c, r := range []rune(reverse(fmt.Sprintf("%v", n))) {
			sc.SetContent(w-1-c, h-2-l, r, nil, st)
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

	st = tcell.StyleDefault

	sc.Clear()
}

func main() {
	stack := []float64{}

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

				case '+', '-', '*', '/', 'p':
					if len(input) > 0 && len(stack) > 0 {
						if n, err := strconv.ParseFloat(reverse(string(input)), 32); err == nil {
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

							sc.Clear()
							setStackContent(stack, w, h)
							sc.Show()
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

						sc.Clear()
						setStackContent(stack, w, h)
						sc.Show()
					}

				case 'i', 't':
					if len(input) > 0 {
						switch ev.Rune() {
						case 'i':
							if input[len(input)-1] == '-' {
								input = input[:len(input)-1]

								sc.SetContent(w-len(input)-1, h-1, 0, nil, st)
							} else {
								input = append(input, '-')

								sc.SetContent(w-len(input), h-1, '-', nil, st)
							}

						case 't':
							if n, err := strconv.ParseFloat(reverse(string(input)), 32); err == nil {
								clearInput(len(input), w, h)

								input = []rune(reverse(fmt.Sprintf("%v", math.Sqrt(n))))

								setInputContent(input, w, h)
							}
						}

						sc.Show()
					} else if len(stack) > 0 {
						switch ev.Rune() {
						case 'i':
							stack[0] *= -1

						case 't':
							stack[0] = math.Sqrt(stack[0])
						}

						sc.Clear()
						setStackContent(stack, w, h)
						setInputContent(input, w, h)
						sc.Show()
					}

				case 'd', 'u':
					if len(stack) > 0 {
						switch ev.Rune() {
						case 'd':
							stack = stack[1:]

						case 'u':
							stack = append([]float64{stack[0]}, stack...)
						}

						sc.Clear()
						setStackContent(stack, w, h)
						setInputContent(input, w, h)
						sc.Show()
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

						sc.Clear()
						setStackContent(stack, w, h)
						setInputContent(input, w, h)
						sc.Show()
					}

				case 'q':
					sc.Fini()

					os.Exit(0)
				}

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					sc.SetContent(w-len(input), h-1, 0, nil, st)

					input = input[1:]

					setInputContent(input, w, h)
					sc.Show()
				}

			case tcell.KeyEnter:
				if len(input) > 0 {
					if n, err := strconv.ParseFloat(reverse(string(input)), 32); err == nil {
						stack = append([]float64{n}, stack...)

						input = input[:0]

						sc.Clear()
						setStackContent(stack, w, h)
						sc.Show()
					}
				}

			case tcell.KeyCtrlL:
				sc.Sync()

			}

		case *tcell.EventResize:
			sc.Sync()
		}
	}
}

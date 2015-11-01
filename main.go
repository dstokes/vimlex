package main

import (
	"fmt"
	"os"
)

func meta(b []byte) bool {
	return len(b) == 1 && b[0] < 28
}

func chord(b []byte) string {
	return "^" + string(b[0]+64)
}

func contains(s []rune, c int) bool {
	for _, cpt := range s {
		if c == int(cpt) {
			return true
		}
	}
	return false
}

func main() {
	l := lex(os.Stdin)

	// print output
	go func() {
		//var mode string
		for {
			t := <-l.tokens
			var out string
			out = string(t.bytes)
			if meta(t.bytes) {
				out = chord(t.bytes)
			}
			// if mode != t.mode {
			// 	mode = t.mode
			// 	fmt.Fprintf(os.Stdout, " - %s\n", mode)
			// }
			fmt.Fprintf(os.Stdout, "%s", out)
		}
	}()

	l.run()
}

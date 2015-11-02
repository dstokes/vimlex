package main

import (
	"bytes"
	"flag"
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
	// parse flags
	flags := flag.NewFlagSet("cFlags", flag.ContinueOnError)

	all := flags.Bool("a", false, "Show keystrokes in all modes")

	// check flag values
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		return
	}

	l := lex(os.Stdin)
	go l.run()

	// print output
	var mode string
	for t := range l.tokens {
		// discard "replacement" characters
		if bytes.Equal(t.bytes, []byte{239, 191, 189}) {
			continue
		}
		var out string
		out = string(t.bytes)
		if meta(t.bytes) {
			out = chord(t.bytes)
		}
		if mode != t.mode {
			mode = t.mode
			if t.mode != "normal" || *all {
				fmt.Fprintln(os.Stdout, "")
			}
		}
		if t.mode == "normal" || *all {
			fmt.Fprintf(os.Stdout, "%s", out)
		}
	}
}

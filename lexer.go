package main

import (
	"bufio"
	"bytes"
	"io"
)

const (
	ESC   int = 27
	ENTER int = 13
	CTRLC int = 3
	SPACE int = 20
)

var searches = []rune{'/'}
var visuals = []rune{'v', 'V'}
var commands = []rune{':'}
var inserts = []rune{'A', 'C', 'I', 'O', 'R', 'S', 'a', 'c', 'i', 'o', 's', '?'}
var vimBackspace = []byte{194, 128}

type token struct {
	mode  string
	bytes []byte
}
type lexer struct {
	mode    string
	scanner *bufio.Scanner
	tokens  chan *token
}

func (l *lexer) run() {
	for state := lexNormal; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) next() (error, []byte) {
	if l.scanner.Scan() == false {
		return io.EOF, []byte{}
	}
	return nil, l.scanner.Bytes()
}

func lex(r io.Reader) *lexer {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	return &lexer{
		mode:    "normal",
		scanner: scanner,
		tokens:  make(chan *token, 100)}
}

type stateFn func(l *lexer) stateFn

func lexNormal(l *lexer) stateFn {
	for {
		err, b := l.next()
		if err != nil {
			break
		}

		// account for special vim backspace code (<80>kb)
		if len(b) > 1 {
			if bytes.Equal(b, vimBackspace) {
				b = []byte{0x8}
				for i := 0; i < 2; i++ {
					l.next()
				}
			}
			l.tokens <- &token{mode: "normal", bytes: b}
			continue
		}
		code := int(b[0])
		// discard invalid characters
		if code >= 126 {
			continue
		}
		l.tokens <- &token{mode: "normal", bytes: b}
		if contains(inserts, code) {
			return lexInsert(l)
		} else if contains(commands, code) {
			return lexCommand(l)
		} else if contains(searches, code) {
			return lexSearch(l)
		} else if contains(visuals, code) || code == 22 {
			return lexVisual(l)
		}
	}
	return nil
}

func lexInsert(l *lexer) stateFn {
	for {
		err, b := l.next()
		if err != nil {
			break
		}
		code := int(b[0])
		if code == ESC || code == CTRLC {
			l.tokens <- &token{mode: "insert", bytes: b}
			return lexNormal(l)
		}
	}
	return nil
}

func lexCommand(l *lexer) stateFn {
	for {
		err, b := l.next()
		if err != nil {
			break
		}
		code := int(b[0])
		l.tokens <- &token{mode: "command", bytes: b}
		if code == ENTER || code == ESC {
			return lexNormal(l)
		}
	}
	return nil
}

func lexSearch(l *lexer) stateFn {
	for {
		err, b := l.next()
		if err != nil {
			break
		}
		code := int(b[0])
		if code == ENTER || code == ESC {
			l.tokens <- &token{mode: "search", bytes: b}
			return lexNormal(l)
		}
	}
	return nil
}

func lexVisual(l *lexer) stateFn {
	for {
		err, b := l.next()
		if err != nil {
			break
		}

		code := int(b[0])
		l.tokens <- &token{mode: "visual", bytes: b}
		if contains(inserts, code) {
			return lexInsert(l)
		} else if contains(commands, code) {
			return lexCommand(l)
		} else if code == ENTER || code == ESC {
			return lexNormal(l)
		}
	}
	return nil
}

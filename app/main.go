package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf8"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	// Handle \d - digit character class
	if pattern == "\\d" {
		for _, b := range line {
			if b >= '0' && b <= '9' {
				return true, nil
			}
		}
		return false, nil
	}

	// Handle \w - word character class
	if pattern == "\\w" {
		for _, b := range line {
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' {
				return true, nil
			}
		}
		return false, nil
	}

	// Handle positive character groups [abc]
	if len(pattern) > 2 && pattern[0] == '[' && pattern[len(pattern)-1] == ']' {
		chars := pattern[1 : len(pattern)-1]
		
		// Check if it's a negative character group [^abc]
		if len(chars) > 0 && chars[0] == '^' {
			negativeChars := chars[1:]
			for _, b := range line {
				found := false
				for i := 0; i < len(negativeChars); i++ {
					if b == negativeChars[i] {
						found = true
						break
					}
				}
				if !found {
					return true, nil
				}
			}
			return false, nil
		}
		
		// Positive character group
		for _, b := range line {
			for i := 0; i < len(chars); i++ {
				if b == chars[i] {
					return true, nil
				}
			}
		}
		return false, nil
	}

	// Handle single literal character
	if utf8.RuneCountInString(pattern) != 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	return bytes.ContainsAny(line, pattern), nil
}

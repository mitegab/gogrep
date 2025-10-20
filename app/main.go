package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchPattern(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func matchPattern(text, pattern string) (bool, error) {
	// anchors
	startAnchor := false
	endAnchor := false
	if len(pattern) > 0 && pattern[0] == '^' {
		startAnchor = true
		pattern = pattern[1:]
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '$' {
		
		endAnchor = true
		pattern = pattern[:len(pattern)-1]
	}

	if startAnchor && endAnchor {
		pos, _, _ := matchHere(text, pattern, 0, "", false)
		return pos == len(text), nil
	}
	if startAnchor {
		pos, _, _ := matchHere(text, pattern, 0, "", false)
		return pos >= 0, nil
	}
	if endAnchor {
		for i := 0; i <= len(text); i++ {
			pos, _, _ := matchHere(text, pattern, i, "", false)
			if pos == len(text) {
				return true, nil
			}
		}
		return false, nil
	}
	for i := 0; i <= len(text); i++ {
		if pos, _, _ := matchHere(text, pattern, i, "", false); pos >= 0 {
			return true, nil
		}
	}
	return false, nil
}

func matchHere(text, pattern string, textPos int, cap1 string, haveCap bool) (int, string, bool) {
	patPos := 0
	for patPos < len(pattern) {
		// parentheses: alternation or capturing group
		if pattern[patPos] == '(' {
			end := findClosingParenAt(pattern, patPos)
			if end < 0 {
				return -1, cap1, haveCap
			}
			inside := pattern[patPos+1 : end]
			rest := pattern[end+1:]
			alts := splitAlternation(inside)
			if len(alts) > 1 {
				for _, alt := range alts {
					if res, nc, nh := matchHere(text, alt+rest, textPos, cap1, haveCap); res >= 0 {
						return res, nc, nh
					}
				}
				return -1, cap1, haveCap
			}
			// capturing group (single)
			if res, _, _ := matchHere(text, inside, textPos, cap1, haveCap); res >= 0 {
				captured := text[textPos:res]
				return matchHere(text, rest, res, captured, true)
			}
			return -1, cap1, haveCap
		}

		// determine element and possible quantifier
		elemLen := 0
		elem := ""
		if pattern[patPos] == '\\' && patPos+1 < len(pattern) {
			elemLen = 2
			elem = pattern[patPos : patPos+2]
		} else if pattern[patPos] == '[' {
			end := patPos + 1
			for end < len(pattern) && pattern[end] != ']' {
				end++
			}
			if end >= len(pattern) {
				return -1, cap1, haveCap
			}
			elemLen = end - patPos + 1
			elem = pattern[patPos : end+1]
		} else if pattern[patPos] == '.' {
			elemLen = 1
			elem = "."
		} else {
			elemLen = 1
			elem = pattern[patPos : patPos+1]
		}

		// '+' quantifier
		if patPos+elemLen < len(pattern) && pattern[patPos+elemLen] == '+' {
			rest := pattern[patPos+elemLen+1:]
			saved := textPos
			// consume greedily
			count := 0
			for textPos < len(text) && matchElement(text[textPos], elem) {
				textPos++
				count++
			}
			if count == 0 {
				return -1, cap1, haveCap
			}
			for tp := textPos; tp >= saved+1; tp-- {
				if res, nc, nh := matchHere(text, rest, tp, cap1, haveCap); res >= 0 {
					return res, nc, nh
				}
			}
			return -1, cap1, haveCap
		}
		// '?' quantifier
		if patPos+elemLen < len(pattern) && pattern[patPos+elemLen] == '?' {
			rest := pattern[patPos+elemLen+1:]
			// try with element
			if textPos < len(text) && matchElement(text[textPos], elem) {
				if res, nc, nh := matchHere(text, rest, textPos+1, cap1, haveCap); res >= 0 {
					return res, nc, nh
				}
			}
			// try without element
			return matchHere(text, rest, textPos, cap1, haveCap)
		}

		// consume single element
		if elem == "\\d" {
			if textPos >= len(text) || text[textPos] < '0' || text[textPos] > '9' {
				return -1, cap1, haveCap
			}
			textPos++
			patPos += elemLen
			continue
		}
		if elem == "\\w" {
			if textPos >= len(text) {
				return -1, cap1, haveCap
			}
			b := text[textPos]
			if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_') {
				return -1, cap1, haveCap
			}
			textPos++
			patPos += elemLen
			continue
		}
		if elem == "\\1" {
			if !haveCap {
				return -1, cap1, haveCap
			}
			l := len(cap1)
			if textPos+l > len(text) || text[textPos:textPos+l] != cap1 {
				return -1, cap1, haveCap
			}
			textPos += l
			patPos += elemLen
			continue
		}
		if len(elem) > 0 && elem[0] == '[' {
			group := elem[1 : len(elem)-1]
			if textPos >= len(text) || !matchCharGroup(text[textPos], group) {
				return -1, cap1, haveCap
			}
			textPos++
			patPos += elemLen
			continue
		}
		if elem == "." {
			if textPos >= len(text) {
				return -1, cap1, haveCap
			}
			textPos++
			patPos += elemLen
			continue
		}
		// literal
		if textPos >= len(text) || text[textPos] != elem[0] {
			return -1, cap1, haveCap
		}
		textPos++
		patPos += elemLen
	}
	return textPos, cap1, haveCap
}

func matchElement(char byte, elem string) bool {
	if elem == "." {
		return true
	}
	if elem == "\\d" {
		return char >= '0' && char <= '9'
	}
	if elem == "\\w" {
		return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_'
	}
	if len(elem) > 0 && elem[0] == '[' && elem[len(elem)-1] == ']' {
		return matchCharGroup(char, elem[1:len(elem)-1])
	}
	if len(elem) == 1 {
		return char == elem[0]
	}
	return false
}

func matchCharGroup(char byte, group string) bool {
	if len(group) > 0 && group[0] == '^' {
		negGroup := group[1:]
		for i := 0; i < len(negGroup); i++ {
			if char == negGroup[i] {
				return false
			}
		}
		return true
	}
	
	for i := 0; i < len(group); i++ {
		if char == group[i] {
			return true
		}
	}
	return false
}

func findClosingParenAt(s string, start int) int {
	if start >= len(s) || s[start] != '(' {
		return -1
	}
	depth := 0
	for i := start; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func splitAlternation(s string) []string {
	var result []string
	var current string
	depth := 0
	
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			depth++
			current += string(s[i])
		} else if s[i] == ')' {
			depth--
			current += string(s[i])
		} else if s[i] == '|' && depth == 0 {
			result = append(result, current)
			current = ""
		} else {
			current += string(s[i])
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}


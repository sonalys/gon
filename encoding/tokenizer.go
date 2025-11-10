package encoding

import (
	"bytes"
	"unicode"
)

type Token = []byte

func tokenize(input Token) []Token {
	var tokens []Token
	var curTokenStartIndex int
	var inString bool
	var inComment bool

	for i, r := range input {
		resetCursor := func() {
			curTokenStartIndex = i + 1
		}

		getCurrent := func(inclusive bool) Token {
			defer resetCursor()
			if inclusive {
				return input[curTokenStartIndex : i+1]
			}
			return input[curTokenStartIndex:i]
		}

		curLength := i - curTokenStartIndex

		switch {
		case r == '\n':
			if !inComment {
				if curLength > 0 {
					tokens = append(tokens, getCurrent(false))
				}
			}
			inComment = false
			resetCursor()
		case inComment:
		case len(input) > i+2 && bytes.Equal(input[i:i+2], Token("//")):
			inComment = true
		case r == '"' && !inString:
			inString = true
		case r == '"' && inString:
			tokens = append(tokens, bytes.TrimSpace(getCurrent(true)))
			inString = false
		case inString:
		case unicode.IsSpace(rune(r)):
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			resetCursor()
		case bytes.Contains(Token("():"), input[i:i+1]):
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			tokens = append(tokens, input[i:i+1])
			resetCursor()
		case r == ',':
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			resetCursor()
		default:
		}
	}
	if curTokenStartIndex < len(input) {
		tokens = append(tokens, input[curTokenStartIndex:])
	}
	return tokens
}

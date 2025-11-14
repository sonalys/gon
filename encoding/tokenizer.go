package encoding

import (
	"bytes"
	"unicode"
)

type Token struct {
	content  []byte
	pos, end int
}

func tokenize(input []byte) []Token {
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
				return Token{
					content: input[curTokenStartIndex : i+1],
					pos:     curTokenStartIndex,
					end:     i + 1,
				}
			}
			return Token{
				content: input[curTokenStartIndex:i],
				pos:     curTokenStartIndex,
				end:     i,
			}
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
		case len(input) > i+2 && bytes.Equal(input[i:i+2], []byte("//")):
			inComment = true
		case r == '"' && !inString:
			inString = true
		case r == '"' && inString:
			tokens = append(tokens, getCurrent(true))
			inString = false
		case inString:
		case unicode.IsSpace(rune(r)):
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			resetCursor()
		case bytes.Contains([]byte("():"), input[i:i+1]):
			if curLength > 0 {
				tokens = append(tokens, getCurrent(false))
			}
			tokens = append(tokens, Token{
				content: input[i : i+1],
				pos:     i,
				end:     i + 1,
			})
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
		tokens = append(tokens, Token{
			content: input[curTokenStartIndex:],
			pos:     curTokenStartIndex,
			end:     len(input),
		})
	}
	return tokens
}

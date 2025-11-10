package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_tokenize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := `if(
		  condition: equal(
		      first: true,second: friend.name
			  else: "third",
		  ),
		)`

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			Token("if"),
			Token("("),
			Token("condition"),
			Token(":"),
			Token("equal"),
			Token("("),
			Token("first"),
			Token(":"),
			Token("true"),
			Token("second"),
			Token(":"),
			Token("friend.name"),
			Token("else"),
			Token(":"),
			Token("\"third\""),
			Token(")"),
			Token(")"),
		}

		t.Logf("%+v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("inlined", func(t *testing.T) {
		input := `if(condition: equal(first: true,second: friend.name else: "third",),)`

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			Token("if"),
			Token("("),
			Token("condition"),
			Token(":"),
			Token("equal"),
			Token("("),
			Token("first"),
			Token(":"),
			Token("true"),
			Token("second"),
			Token(":"),
			Token("friend.name"),
			Token("else"),
			Token(":"),
			Token("\"third\""),
			Token(")"),
			Token(")"),
		}

		t.Logf("%+v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("no space", func(t *testing.T) {
		input := `if(condition:equal(first:true,second:friend.name,else:"third",),)`

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			Token("if"),
			Token("("),
			Token("condition"),
			Token(":"),
			Token("equal"),
			Token("("),
			Token("first"),
			Token(":"),
			Token("true"),
			Token("second"),
			Token(":"),
			Token("friend.name"),
			Token("else"),
			Token(":"),
			Token("\"third\""),
			Token(")"),
			Token(")"),
		}

		t.Logf("%+v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("comment should be ignored", func(t *testing.T) {
		input := "// comment\nif()"

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			Token("if"),
			Token("("),
			Token(")"),
		}

		t.Logf("%+v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})
}

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
			{content: []uint8{0x69, 0x66}, pos: 0, end: 2},
			{content: []uint8{0x28}, pos: 2, end: 3},
			{content: []uint8{0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e}, pos: 8, end: 17},
			{content: []uint8{0x3a}, pos: 17, end: 18},
			{content: []uint8{0x65, 0x71, 0x75, 0x61, 0x6c}, pos: 19, end: 24},
			{content: []uint8{0x28}, pos: 24, end: 25},
			{content: []uint8{0x66, 0x69, 0x72, 0x73, 0x74}, pos: 34, end: 39},
			{content: []uint8{0x3a}, pos: 39, end: 40},
			{content: []uint8{0x74, 0x72, 0x75, 0x65}, pos: 41, end: 45},
			{content: []uint8{0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64}, pos: 46, end: 52},
			{content: []uint8{0x3a}, pos: 52, end: 53},
			{content: []uint8{0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x2e, 0x6e, 0x61, 0x6d, 0x65}, pos: 54, end: 65},
			{content: []uint8{0x65, 0x6c, 0x73, 0x65}, pos: 71, end: 75},
			{content: []uint8{0x3a}, pos: 75, end: 76},
			{content: []uint8{0x22, 0x74, 0x68, 0x69, 0x72, 0x64, 0x22}, pos: 77, end: 84},
			{content: []uint8{0x29}, pos: 90, end: 91},
			{content: []uint8{0x29}, pos: 95, end: 96},
		}

		t.Logf("%#v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("inlined", func(t *testing.T) {
		input := `if(condition: equal(first: true,second: friend.name else: "third",),)`

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			{content: []uint8{0x69, 0x66}, pos: 0, end: 2},
			{content: []uint8{0x28}, pos: 2, end: 3},
			{content: []uint8{0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e}, pos: 3, end: 12},
			{content: []uint8{0x3a}, pos: 12, end: 13},
			{content: []uint8{0x65, 0x71, 0x75, 0x61, 0x6c}, pos: 14, end: 19},
			{content: []uint8{0x28}, pos: 19, end: 20},
			{content: []uint8{0x66, 0x69, 0x72, 0x73, 0x74}, pos: 20, end: 25},
			{content: []uint8{0x3a}, pos: 25, end: 26},
			{content: []uint8{0x74, 0x72, 0x75, 0x65}, pos: 27, end: 31},
			{content: []uint8{0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64}, pos: 32, end: 38},
			{content: []uint8{0x3a}, pos: 38, end: 39},
			{content: []uint8{0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x2e, 0x6e, 0x61, 0x6d, 0x65}, pos: 40, end: 51},
			{content: []uint8{0x65, 0x6c, 0x73, 0x65}, pos: 52, end: 56},
			{content: []uint8{0x3a}, pos: 56, end: 57},
			{content: []uint8{0x22, 0x74, 0x68, 0x69, 0x72, 0x64, 0x22}, pos: 58, end: 65},
			{content: []uint8{0x29}, pos: 66, end: 67},
			{content: []uint8{0x29}, pos: 68, end: 69},
		}

		t.Logf("%#v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("no space", func(t *testing.T) {
		input := `if(condition:equal(first:true,second:friend.name,else:"third",),)`

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			{content: []uint8{0x69, 0x66}, pos: 0, end: 2},
			{content: []uint8{0x28}, pos: 2, end: 3},
			{content: []uint8{0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e}, pos: 3, end: 12},
			{content: []uint8{0x3a}, pos: 12, end: 13},
			{content: []uint8{0x65, 0x71, 0x75, 0x61, 0x6c}, pos: 13, end: 18},
			{content: []uint8{0x28}, pos: 18, end: 19},
			{content: []uint8{0x66, 0x69, 0x72, 0x73, 0x74}, pos: 19, end: 24},
			{content: []uint8{0x3a}, pos: 24, end: 25},
			{content: []uint8{0x74, 0x72, 0x75, 0x65}, pos: 25, end: 29},
			{content: []uint8{0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64}, pos: 30, end: 36},
			{content: []uint8{0x3a}, pos: 36, end: 37},
			{content: []uint8{0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x2e, 0x6e, 0x61, 0x6d, 0x65}, pos: 37, end: 48},
			{content: []uint8{0x65, 0x6c, 0x73, 0x65}, pos: 49, end: 53},
			{content: []uint8{0x3a}, pos: 53, end: 54},
			{content: []uint8{0x22, 0x74, 0x68, 0x69, 0x72, 0x64, 0x22}, pos: 54, end: 61},
			{content: []uint8{0x29}, pos: 62, end: 63},
			{content: []uint8{0x29}, pos: 64, end: 65},
		}

		t.Logf("%#v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("comment should be ignored", func(t *testing.T) {
		input := "// comment\nif()"

		tokens := tokenize([]byte(input))
		expectedTokens := []Token{
			{content: []uint8{0x69, 0x66}, pos: 11, end: 13},
			{content: []uint8{0x28}, pos: 13, end: 14},
			{content: []uint8{0x29}, pos: 14, end: 15},
		}

		t.Logf("%#v", tokens)

		assert.Equal(t, expectedTokens, tokens)
	})
}

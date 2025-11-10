package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DecodePretty(t *testing.T) {
	input := `if(
		  condition: equal(
		      first: true
		      second: friend.name
		  ),
		  then: call(
		      name: "reply"
		      target: 10.23
			  else: 102
		  ),
		  else: call("whoAreYou")
		)`

	got, err := Decode([]byte(input), DefaultExpressionCodex)
	require.NoError(t, err)

	err = HumanEncode(t.Output(), got)
	require.NoError(t, err)
}

func Test_DecodeInlined(t *testing.T) {
	input := `if(condition: equal(first: myName, second: friend.name),then: call(name:"reply",target: 10.23,else: 102),else: call("whoAreYou"))`
	got, err := Decode([]byte(input), DefaultExpressionCodex)
	require.NoError(t, err)

	err = HumanEncode(t.Output(), got)
	require.NoError(t, err)
}

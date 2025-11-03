package goncoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Decode(t *testing.T) {
	input := `if(
		  // comment block
		  condition: equal(
		      first: myName,
		      second: friend.name,
		  ),
		  then: call(
		      name: "reply",
		      target: 10.23,
			  else: 102,
		  ),
		  else: call("whoAreYou")
		)`

	got, err := Decode([]byte(input))
	require.NoError(t, err)

	err = Encode(t.Output(), got)
	require.NoError(t, err)

	t.Fail()
}

package main_test

import (
	"fmt"

	"github.com/sonalys/gon"
	"github.com/sonalys/gon/encoding"
)

func Example_objectAccessRule() {
	type File struct {
		UID  int64  `gon:"uid"`
		GID  int64  `gon:"gid"`
		Path string `gon:"path"`
	}

	file := &File{
		UID:  1000,
		GID:  100,
		Path: "/shared/file_001.txt",
	}

	scope, err := gon.
		NewScope().
		WithDefinitions(gon.Definitions{
			"file": gon.Literal(file),
		})
	if err != nil {
		panic(err)
	}

	fileAccessPolicy := `if(
		condition: or(
			equal(file.uid, 1023),
			equal(file.gid, 102),
			hasPrefix(file.path, "/shared")
		),
		then: "pass",
		else: "fail"
	)`

	policy, err := encoding.Decode([]byte(fileAccessPolicy), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	value, err := scope.Compute(policy)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	// Output:
	// pass
}

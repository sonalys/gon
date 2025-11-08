package main_test

import (
	"fmt"
	"os"

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
		then: true,
		else: false
	)`

	policy, err := encoding.Decode([]byte(fileAccessPolicy), encoding.DefaultExpressionCodex)
	if err != nil {
		panic(err)
	}

	isAuthorized, err := scope.Compute(policy)
	if err != nil {
		panic(err)
	}

	err = encoding.HumanEncode(os.Stdout, policy)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nIs authorized? %v", isAuthorized)

	// Output:
	// if(
	// 	condition: or(
	// 		equal(
	// 			first: file.uid,
	// 			second: 1023
	// 		),
	// 		equal(
	// 			first: file.gid,
	// 			second: 102
	// 		),
	// 		hasPrefix(
	// 			text: file.path,
	// 			prefix: "/shared"
	// 		)
	// 	),
	// 	then: true,
	// 	else: false
	// )
	// Is authorized? true
}

package gon

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Encode(w io.Writer, root Expression) error {
	return encodeBody(w, root, 0)
}

func encodeBody(w io.Writer, root Expression, indentation int) error {
	print := func(indentation int, mask string, args ...any) {
		if indentation > 0 {
			fmt.Fprint(w, strings.Repeat("\t", indentation))
		}
		fmt.Fprintf(w, mask, args...)
	}

	switch root.Type() {
	case ExpressionTypeInvalid:
		return errors.New("cannot encode invalid expression type")
	case ExpressionTypeOperation:
		name, args := root.Name()

		print(0, "%s(", name)

		for i, arg := range args {
			if i == 0 && args[0].Key != "" || i > 0 {
				print(0, "\n")
				print(indentation+1, "")
			}
			if arg.Key != "" {
				print(0, "%s: ", arg.Key)
			}

			if err := encodeBody(w, arg.Value, indentation+1); err != nil {
				return err
			}
		}

		if len(args) > 1 {
			print(0, "\n")
			print(indentation, ")")
			break
		}

		print(0, ")")
	case ExpressionTypeReference:
		_, args := root.Name()
		key := args[0].Value.(interface{ Any() any }).Any()
		print(0, "%v", key)
	case ExpressionTypeValue:
		value := root.(interface{ Any() any }).Any()
		if str, ok := value.(string); ok {
			value = strconv.Quote(str)
		}
		print(0, "%v", value)
	}
	return nil
}

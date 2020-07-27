package cli

import (
	"fmt"

	prettyJSON "github.com/hokaccha/go-prettyjson"
)

func Print(input interface{}) error {
	formatter := prettyJSON.NewFormatter()
	bytes, err := formatter.Marshal(input)
	fmt.Println(string(bytes))
	if err != nil {
		return err
	}
	return nil
}

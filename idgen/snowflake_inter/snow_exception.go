package snowflake

import "fmt"

type idGeneratorException struct {
	message string
	error   error
}

func (e idGeneratorException) IdGeneratorException(message ...interface{}) {
	fmt.Println(message)
}

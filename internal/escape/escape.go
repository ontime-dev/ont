package escape

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func ErrorWithZeroRC(err error) {
	fmt.Println(err.Error())
	os.Exit(0)
}

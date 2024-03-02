package escape

import (
	"fmt"
	"os"
)

func Error(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

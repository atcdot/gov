package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	writeToFile bool
}

func (l Logger) Log(a any) {
	fmt.Println(a)

	if l.writeToFile {
		f, err := os.OpenFile("/tmp/gov.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0644))
		if err != nil {
			fmt.Println(err)
		}

		fileContent := fmt.Sprintf("%v\n", a)

		_, err = f.Write([]byte(fileContent))
		if err != nil {
			fmt.Println(err)
		}
	}
}

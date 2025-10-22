package utils

import (
	"log"
	"os"
)

func WaitForEnter(prompt string) {
	os.Stdout.WriteString(prompt)
	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		if b[0] == '\n' {
			break
		}
	}
}

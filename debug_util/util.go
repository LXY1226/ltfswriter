package util

import (
	"encoding/hex"
	"os"
)

func DumpHex(dat []byte) {
	hex.Dumper(os.Stdout).Write(dat)
	os.Stdout.WriteString("\n")
}

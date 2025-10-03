package main

import (
	"fmt"
	"os"

	"github.com/LXY1226/ltfswriter/tape"
)

func main() {
	drive, err := tape.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	drive.CheckParts()
	drive.DumpCapacity()
	drive.Locate10PartBlock(tape.Locate10FlagWithPart, 1, 1)
	pos, err := drive.ReadPosition()
	if err != nil {
		panic(err)
	}
	fmt.Println(pos)
	for i := byte(0); i < 10; i++ {
		drive.Locate10PartBlock(tape.Locate10FlagWithPart, i, 0)
		f, err := os.OpenFile(fmt.Sprintf("part-a%02d.bin", i), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		err = drive.WriteTo(f)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

package main

import (
	"encoding/hex"
	"github.com/LXY1226/ltfswriter/ltfs"
	"github.com/LXY1226/ltfswriter/tape"
	"log"
	"os"
)

func main() {
	drive, err := tape.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = drive.MTSetOptions(tape.MTSTBOOLEANS |
		tape.MTSTBUFFERWRITES | tape.MTSTASYNCWRITES | tape.MTSTCANBSR | tape.MTSTCANPARTITIONS |
		tape.MTSTNOWAITEOF | tape.MTSTSCSI2LOGICAL |
		tape.MTSTDEBUGGING | tape.MTSTSILI | tape.MTSTSYSV)
	if err != nil {
		panic(err)
	}
	//err = drive.MTSetOptions(tape.MTSTDEFBLKSIZE | 0xfffffff)
	//if err != nil {
	//	panic(err)
	//}

	//err = drive.MTSetBlock(0)
	//if err != nil {
	//	panic(err)
	//}
	vol, err := ltfs.Open(drive)
	log.Println(err)
	//if err != nil {
	//	panic(err)
	//}
	log.Println(vol.Vol1Label.String())
	log.Println(vol.LatestIndex.GenerationNumber)
	log.Println(vol.LatestIndex.Directory.Name)

	//buf := make([]byte, 8<<20)
	//for range 100 {
	//	start := time.Now()
	//	n, err := io.ReadFull(drive.File, buf)
	//	end := time.Now()
	//	log.Println(n, "of", cap(buf), "read in", end.Sub(start))
	//	//if n != 0 {
	//	//	if n > 128 {
	//	//		hex.Dumper(os.Stdout).Write(buf[:128])
	//	//	} else {
	//	//		hex.Dumper(os.Stdout).Write(buf[:n])
	//	//	}
	//	//	os.Stdout.WriteString("\n")
	//	//}
	//	if err != nil {
	//		log.Println(err)
	//		if err == io.EOF {
	//			// TODO check status EOF?EOM
	//			continue
	//		}
	//	}
	//}
	//drive.CheckParts()
	//drive.DumpCapacity()
	//drive.Locate10PartBlock(tape.Locate10FlagWithPart, 1, 1)
	//pos, err := drive.ReadPosition()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(pos)
	//for i := byte(0); i < 10; i++ {
	//	drive.Locate10PartBlock(tape.Locate10FlagWithPart, i, 0)
	//	f, err := os.OpenFile(fmt.Sprintf("part-a%02d.bin", i), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	//	err = drive.WriteTo(f)
	//	if err != nil {
	//		panic(err)
	//	}
	//	f.Close()
	//}
}

func dumpHex(dat []byte) {
	hex.Dumper(os.Stdout).Write(dat)
	os.Stdout.WriteString("\n")
}

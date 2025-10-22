package main

import (
	"archive/tar"
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/LXY1226/ltfswriter/tape"
	"github.com/LXY1226/ltfswriter/utils"
	"github.com/ncw/directio"
	"github.com/zeebo/blake3"
)

type Task struct {
	Tapes []string `json:"tapes"`
}

func LoadJson[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var v T
	err = decoder.Decode(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

var sch = tape.NewMediaChanger("/dev/sch0")

func main() {
	task, err := LoadJson[Task]("run.json")
	if err != nil {
		log.Fatal(err)
	}
	zstdCmd := exec.Command("zstd", "-d", "-v")
	zstdOut, err := zstdCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	zstdIn, err := zstdCmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	zstdCmd.Stderr = os.Stderr
	go func() {
		err := zstdCmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// async zstd (channel)
	tarReader := tar.NewReader(zstdOut)
	// async untar + hasher(blake3) + logger(buffered writer)
	out, err := os.Create("run_out.tsv")
	if err != nil {
		log.Fatal(err)
	}
	bufOut := bufio.NewWriter(out)
	defer bufOut.Flush()
	go func() {
		for {
			tarHeader, err := tarReader.Next()
			if err != nil && err != tar.ErrInsecurePath {
				panic(err)
			}
			bufOut.Write([]byte(tarHeader.Name))
			bufOut.WriteByte('\t')
			bufOut.WriteString(strconv.FormatInt(tarHeader.Size, 10))
			hasher := blake3.New()
			n, err := io.CopyBuffer(hasher, tarReader, make([]byte, 1024*1024))
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n != tarHeader.Size {
				log.Println("inconsistent size for tar:", n, tarHeader.Size)
			}
			bufOut.WriteByte('\t')
			bufOut.WriteString(hex.EncodeToString(hasher.Sum(nil)))
			bufOut.WriteByte('\n')
		}
	}()
	for _, tapeTag := range task.Tapes {
		log.Println("Ready for", tapeTag)
		TryLoadByTag(tapeTag, 0)
		// open drive
		drive := EnsureOpenDrive(tapeTag, "/dev/st0")
		//drive.MTSeek()
		// read out (512KB block size) & drop to zstd
		buf := directio.AlignedBlock(1024 * 1024)
		written := int64(0)
		for {
			nr, er := drive.Read(buf)
			if nr > 0 {
				//fmt.Println(nr)
				nw, ew := zstdIn.Write(buf[0:nr])
				if nw < 0 || nr < nw {
					nw = 0
					if ew == nil {
						panic(errors.New("invalid write result"))
						//ew = errInvalidWrite
					}
				}
				written += int64(nw)
				if ew != nil {
					err = ew
					break
				}
				if nr != nw {
					panic(io.ErrShortWrite)
				}
			}
			if er != nil {
				if er != io.EOF {
					panic(err)
				}
				break
			}
		}
		log.Println(tapeTag, "read", written, "bytes")
		// close drive
		drive.Close()
		// unload
		sch.Unload(0)
	}
}

func EnsureOpenDrive(tag string, drivePath string) *tape.Drive {
	for {
		drive, err := tape.Open(drivePath)
		if err != nil {
			log.Println("Error opening drive:", err)
			goto retry
		}

		err = drive.MTSetOptions(tape.MTSTBOOLEANS |
			tape.MTSTBUFFERWRITES | tape.MTSTASYNCWRITES | tape.MTSTCANBSR | tape.MTSTCANPARTITIONS |
			tape.MTSTNOWAITEOF | tape.MTSTSCSI2LOGICAL |
			tape.MTSTDEBUGGING | tape.MTSTSILI | tape.MTSTSYSV)
		if err != nil {
			panic(err)
		}
		return drive
	retry:
		utils.WaitForEnter(fmt.Sprintln("Please change the tape manually for", tag, "into", drivePath))
	}
}

func TryLoadByTag(tag string, driveID int) bool {
	tap, err := sch.GetLibraryInv(tag)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(tap) == 0 {
		log.Println("media", tag, "not found")
		return false
	}
	if tap[0].Drive == -1 {
		log.Println("Loading", tag, "from", tap[0].SlotID, "to drive", driveID)
		err = sch.LoadTo(tap[0].SlotID, driveID)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	if tap[0].Drive != driveID {
		log.Println("media", tag, "already loaded in different drive", driveID, "trying unload")
		err = sch.Unload(tap[0].Drive)
		if err != nil {
			log.Println(err)
			return false
		}
		err = sch.LoadTo(tap[0].SlotID, driveID)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return true
}

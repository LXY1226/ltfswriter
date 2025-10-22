//go:build linux

package tape

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func NewMediaChanger(devPath string) *MediaChanger {
	return &MediaChanger{devPath: devPath}
}

type MediaChanger struct {
	devPath string
}

type MediaChangerTape struct {
	SlotID int
	Drive  int
	Tag    string
}

/*
root@pve:~# mtx -f /dev/sch0 status
  Storage Changer /dev/sch0:2 Drives, 24 Slots ( 1 Import/Export )
Data Transfer Element 0:Full (Storage Element 22 Loaded):VolumeTag = 000057L5
Data Transfer Element 1:Empty
      Storage Element 1:Empty
      Storage Element 2:Empty
      Storage Element 3:Empty
      Storage Element 4:Empty
      Storage Element 5:Empty
      Storage Element 6:Empty
      Storage Element 7:Empty
      Storage Element 8:Empty
      Storage Element 9:Empty
      Storage Element 10:Empty
      Storage Element 11:Empty
      Storage Element 12:Empty
      Storage Element 13:Full :VolumeTag=000052L5
      Storage Element 14:Empty
      Storage Element 15:Full :VolumeTag=000055L5
      Storage Element 16:Full :VolumeTag=000054L5
      Storage Element 17:Full :VolumeTag=000051L5
      Storage Element 18:Empty
      Storage Element 19:Empty
      Storage Element 20:Full :VolumeTag=000053L5
      Storage Element 21:Full :VolumeTag=000050L5
      Storage Element 22:Empty
      Storage Element 23:Empty
      Storage Element 24 IMPORT/EXPORT:Empty
*/

// ErrMediaNotFound indicates specified media not existed in the changer
var ErrMediaNotFound = errors.New("media not found")

func in[T comparable](a T, arr []T) bool {
	for _, b := range arr {
		if b == a {
			return true
		}
	}
	return false
}

func (sch *MediaChanger) GetLibraryInv(only ...string) ([]MediaChangerTape, error) {
	cmd := exec.Command("mtx", "-f", sch.devPath, "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var tapes []MediaChangerTape
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.Contains(line, "Full") {
			continue
		}
		tap := MediaChangerTape{Drive: -1}
		i := strings.LastIndexByte(line, '=')
		tap.Tag = strings.TrimSpace(line[i+1:])
		if len(only) > 0 && !in(tap.Tag, only) {
			continue
		}

		elements := strings.Split(line, "Element ")
		if len(elements) < 2 {
			log.Print("invalid line from GetLibraryInv: ", line)
			continue
		}

		i = strings.Index(elements[1], ":")
		tap.SlotID, err = strconv.Atoi(elements[1][:i])
		if err != nil {
			log.Print("invalid line from GetLibraryInv: ", line)
			continue
		}
		if len(elements) > 2 {
			tap.Drive = tap.SlotID
			i = strings.Index(elements[2], " ")
			tap.SlotID, err = strconv.Atoi(elements[2][:i])
			if err != nil {
				log.Print("invalid line from GetLibraryInv: ", line)
				continue
			}
		}
		tapes = append(tapes, tap)
		if len(tapes) == len(only) {
			break
		}
	}
	return tapes, nil
}

func (sch *MediaChanger) LoadTo(position int, driveID int) error {
	return exec.Command("mtx", "-f", sch.devPath, "load", strconv.Itoa(position), "drive"+strconv.Itoa(driveID)).Run()
}

func (sch *MediaChanger) Unload(driveID int) error {
	return exec.Command("mtx", "-f", sch.devPath, "unload", "drive"+strconv.Itoa(driveID)).Run()
}

func (sch *MediaChanger) UnloadTo(position int, driveID int) error {
	panic("not implemented")
}

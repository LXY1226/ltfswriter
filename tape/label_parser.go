package tape

import (
	"errors"
	"strings"
)

type VOL1Label struct {
	VolID            [6]byte // [6]char
	VolAccessibility byte
	ImplID           [13]byte // string
	OwnerID          [14]byte // string
}

func ParseVol1Label(dat []byte) (lab VOL1Label, err error) {
	if len(dat) < 80 {
		return lab, errors.New("vol1 label too short")
	}
	if string(dat[:4]) != "VOL1" {
		return lab, errors.New("bad vol1 label head")
	}
	copy(lab.VolID[:], dat[4:])
	lab.VolAccessibility = dat[10]
	copy(lab.ImplID[:], dat[24:])
	copy(lab.OwnerID[:], dat[37:])
	// TODO other checks
	//if len(dat) != 80 {
	//	log.Println("Warning: bad vol1")
	//}
	return lab, nil
}

func (lab *VOL1Label) String() string {
	sb := new(strings.Builder)
	sb.WriteString("Volume ID:")
	sb.WriteString(trimSuffixByte(lab.VolID[:], 0))
	sb.WriteString("\nImplementation ID:")
	sb.WriteString(trimSuffixByte(lab.ImplID[:], ' '))
	sb.WriteString("\nOwner ID:")
	sb.WriteString(trimSuffixByte(lab.OwnerID[:], ' '))
	return sb.String()
}

func trimSuffixByte(b []byte, v byte) string {
	s := string(b)
	i := strings.LastIndexByte(s, v)
	if i == -1 {
		return s
	}
	return s[:i]
}

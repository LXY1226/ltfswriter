package tape

import (
	"fmt"
	"io"
)

const shortLimit = 1 << 20

func (d Drive) MTMustReadShortFile() []byte {
	b, err := d.MTReadShortFile()
	if err != io.EOF {
		panic(err)
	}
	return b
}

func (d Drive) MTReadShortFile() ([]byte, error) {
	b := make([]byte, shortLimit)
	n, err := d.MTReadFull(b)
	if err == nil {
		return b, fmt.Errorf("file too big for %d bytes", shortLimit)
	}
	if n < shortLimit/2 {
		b2 := make([]byte, n)
		copy(b2, b[:n])
		b = b2
	}
	return b, err
}

func (d Drive) MTReadFull(buf []byte) (n int, err error) {
	for n < len(buf) && err == nil {
		var nn int
		nn, err = d.File.Read(buf[n:])
		n += nn
	}
	return
}

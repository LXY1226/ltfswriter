package tape

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
)

const senseBufferSize = 32

type senseError []byte

var senseKeyDesc = [16]string{
	0x00: "NO SENSE",
	0x01: "RECOVERED ERROR",
	0x02: "NOT READY",
	0x03: "MEDIUM ERROR",
	0x04: "HARDWARE ERROR",
	0x05: "ILLEGAL REQUEST",
	0x06: "UNIT ATTENTION",
	0x07: "DATA PROTECT",
	0x08: "BLANK CHECK",
	0x09: "VENDOR SPECIFIC",
	0x0a: "COPY ABORTED",
	0x0b: "ABORTED COMMAND",
	0x0d: "VOLUME OVERFLOW",
	0x0e: "MISCOMPARE",
	0x0f: "COMPLETED",
}

func (s senseError) Error() string {
	sb := new(bytes.Buffer)
	sb.WriteString("SCSI sense error: ")
	sb.WriteString(senseKeyDesc[s[2]])
	sb.WriteString(" ")
	b := hex.AppendEncode(sb.Bytes(), s[12:12]) // Additional Sense Code and Qualifier
	b = append(b, '/')
	b = hex.AppendEncode(b, s[13:13])
	return string(b)
}

func (d Drive) TestUnitReady() error {
	d.scsiCmd([]byte{ScsiOpTestUnitReady, 0, 0, 0, 0, 0}, 60_000)
	return nil
}

func (d Drive) CheckParts() {
	dat, err := d.scsiRead([]byte{
		ScsiOpModeSense6, 0,
		ModePageMediumPartitions, 0,
		0xff, 0, // bufLen
	}, 0xff, 60_000)
	if err != nil {
		panic(err)
	}
	if len(dat) < 0xf {
		log.Fatalln("CheckParts: short read", len(dat))
	}
	if dat[0xf] != 1 {
		log.Fatalln("CheckParts: expected 2 partitions, got", dat[0xf]+1)
	}
}

func (d Drive) DumpCapacity() {
	//dat, err := d.scsiRead([]byte{
	//	ScsiOpLogSense, 0,
	//	LogPageSupportedPages, 0, 0,
	//	0, 0, // ParameterPointer
	//	0x0, 255, // bufLen
	//	0}, 255, 60_000)
	//if err != nil {
	//	panic(err)
	//}
	//dumpHex(dat)
	/*
		00000000  00 00 00 15 00 02 03 0c  0d 11 12 13 14 15 16 17  |................|
		00000010  1b 2e 30 31 32 33 34 35  3e 00 00 00 00 00 00 00  |..012345>.......|
	*/
	dat, err := d.scsiRead([]byte{
		ScsiOpLogSense, 0,
		LogPageTapeCapacity, 0, 0,
		0, 0, // ParameterPointer
		0x0, 0x80, // bufLen
		0}, 0x80, 60_000)
	if err != nil {
		panic(err)
	}
	dumpHex(dat)

	/*CurrentCumulativeValues*/
	//0b01_000000 |
}

func (d Drive) LocateBlock(block uint32) error {
	err := d.scsiCmd([]byte{
		ScsiOpLocate10, 0, 0,
		byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block),
		0, 0, 0,
	}, 600_000)
	return err
}

const (
	Locate10FlagWithPart = 0b0000_0010
)

func (d Drive) Locate10PartBlock(flag byte, part byte, block uint32) error {
	err := d.scsiCmd([]byte{
		ScsiOpLocate10, flag, 0,
		byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block),
		0, part, 0,
	}, 600_000)
	return err
}

const (
	Locate16FlagDestObjID  = 0b00_000_000
	Locate16FlagDestFileID = 0b00_001_000
	Locate16FlagDestEOD    = 0b00_011_000
	Locate16FlagWithPart   = 0b00_000_010
)

//func (d Drive) ReadBlockLimits()

func (d Drive) Locate16(flag byte, part byte, logicalID uint64) error {
	err := d.scsiCmd([]byte{
		ScsiOpLocate16, flag, 0, part,
		byte(logicalID >> 56), byte(logicalID >> 48), byte(logicalID >> 40), byte(logicalID >> 32),
		byte(logicalID >> 24), byte(logicalID >> 16), byte(logicalID >> 8), byte(logicalID),
		0, 0, 0, 0,
	}, 600_000)
	return err
}

func (d Drive) Read() ([]byte, error) {
	return d.scsiRead([]byte{
		ScsiOpRead, 0b0000_0010, 0, 4, 00, 0,
	}, 256*1024, 600_000)
}

type PositionData struct {
	Partition uint32
	Block     uint64
	File      uint64
}

func (d Drive) ReadPosition() (PositionData, error) {
	dat, err := d.scsiRead([]byte{
		ScsiOpReadPosition, 6,
		0, 0, 0, 0, 0, // reserved
		0, 0, // must be 0
		0,
	}, 64, 60_000)
	if err != nil {
		return PositionData{}, err
	}
	// BOP
	if dat[0]&0x80 == 0x80 {
		return PositionData{Partition: binary.BigEndian.Uint32(dat[4:])}, nil
	}
	dumpHex(dat)
	return PositionData{
		Partition: binary.BigEndian.Uint32(dat[4:]),
		Block:     binary.BigEndian.Uint64(dat[8:]),
		File:      binary.BigEndian.Uint64(dat[16:]),
	}, nil
}

func dumpHex(dat []byte) {
	hex.Dumper(os.Stdout).Write(dat)
	os.Stdout.WriteString("\n")
}

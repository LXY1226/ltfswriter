//go:build linux

package tape

import (
	"os"
	"unsafe"
)

type Drive struct {
	*os.File           // File handle
	fd         uintptr // File descriptor
	fixedBlock bool    // Whether to use fixed block size
}

type MtOp struct {
	op    int16 // Operation ID
	_     int16 // Padding to match C structures
	count int32 // Operation count
}

type MtStatus struct {
	Type int64 // Type of magtape device
	// Residual count: may be one of:
	// - number of bytes ignored
	// - number of files not skipped
	// - number of records not skipped
	ResID int64
	DsReg int64 // Device dependent status register
	GStat int64 // Generic (device independent) status
	ErReg int64 // Error register

	// The next two fields are not always used
	FileNo int32 // Current file number
	BlkNo  int32 // Current block number
}

// MTRESET resets the drive in case of problems
func (d Drive) MTRESET() error { return d.mtioctop(MTRESET, 0) }

// MTFSFM spaces forward over count position, positioned before filemark
func (d Drive) MTFSFM(count int32) error { return d.mtioctop(MTFSFM, count) }

// MTFSF spaces forward over count filemarks, positioned after filemark.
func (d Drive) MTFSF(count int32) error { return d.mtioctop(MTFSF, count) }

// MTBSF spaces backward over count filemarks, positioned before filemark
func (d Drive) MTBSF(count int32) error { return d.mtioctop(MTBSF, count) }

// MTBSFM spaces backward over count filemark position, positioned after filemark
func (d Drive) MTBSFM(count int32) error { return d.mtioctop(MTBSFM, count) }

// MTFSR spaces forward over count records
func (d Drive) MTFSR(count int32) error { return d.mtioctop(MTFSR, count) }

// MTBSR spaces backward over count records
func (d Drive) MTBSR(count int32) error { return d.mtioctop(MTBSR, count) }

// MTWEOF writes end-of-file record (mark)
func (d Drive) MTWEOF(count int32) error { return d.mtioctop(MTWEOF, count) }

// MTREW rewinds the tape
func (d Drive) MTREW() error { return d.mtioctop(MTREW, 0) }

// MTOFFL rewinds and puts the drive offline (eject tape)
func (d Drive) MTOFFL() error { return d.mtioctop(MTOFFL, 0) }

// MTNOP performs no operation, only sets status (read with MTIOCGET)
func (d Drive) MTNOP() error { return d.mtioctop(MTNOP, 0) }

// MTRETEN retensions the tape
func (d Drive) MTRETEN() error { return d.mtioctop(MTRETEN, 0) }

// MTEOM goes to end of recorded media (for appending files), positioned after the last filemark
func (d Drive) MTEOM() error { return d.mtioctop(MTEOM, 0) }

// MTERASE erases tape -- be careful!
func (d Drive) MTERASE() error { return d.mtioctop(MTERASE, 0) }

// MTRAS1 runs self test 1 (nondestructive)
func (d Drive) MTRAS1() error { return d.mtioctop(MTRAS1, 0) }

// MTRAS2 runs self test 2 (destructive)
func (d Drive) MTRAS2() error { return d.mtioctop(MTRAS2, 0) }

// MTRAS3 reserved for self test 3
func (d Drive) MTRAS3() error { return d.mtioctop(MTRAS3, 0) }

// MTWEOFI writes end-of-file record (mark) in immediate mode
func (d Drive) MTWEOFI(count int32) error { return d.mtioctop(MTWEOFI, count) }

// MTSetBlock sets block length (SCSI)
func (d Drive) MTSetBlock(size int32) error { return d.mtioctop(MTSETBLK, size) }

// MTSetDensity sets tape density (SCSI)
func (d Drive) MTSetDensity(density int32) error { return d.mtioctop(MTSETDENSITY, density) }

// MTSeek seeks to block (Tandberg, etc.)
func (d Drive) MTSeek(pos int32) error { return d.mtioctop(MTSEEK, pos) }

// MTTell gets current block position (Tandberg, etc.)
func (d Drive) MTTell() (pos int32, err error) {
	panic("not implemented")
	// Note: MTTELL requires special handling, this is just a placeholder implementation
	// Specific implementation needs to be adjusted according to driver documentation
	err = d.mtioctop(MTTELL, 0)
	return
}

// MTSetBuffer sets the drive buffering according to SCSI-2
func (d Drive) MTSetBuffer(size int32) error { return d.mtioctop(MTSETDRVBUFFER, size) }

// MTFSS spaces forward over setmarks
func (d Drive) MTFSS(count int32) error { return d.mtioctop(MTFSS, count) }

// MTBSS spaces backward over setmarks
func (d Drive) MTBSS(count int32) error { return d.mtioctop(MTBSS, count) }

// MTWSM writes setmarks
func (d Drive) MTWSM(count int32) error { return d.mtioctop(MTWSM, count) }

// MTLOCK locks the drive door
func (d Drive) MTLOCK() error { return d.mtioctop(MTLOCK, 0) }

// MTUNLOCK unlocks the drive door
func (d Drive) MTUNLOCK() error { return d.mtioctop(MTUNLOCK, 0) }

// MTLOAD executes the SCSI load command
func (d Drive) MTLOAD() error { return d.mtioctop(MTLOAD, 0) }

// MTUNLOAD executes the SCSI unload command
func (d Drive) MTUNLOAD() error { return d.mtioctop(MTUNLOAD, 0) }

// MTCompression controls compression with SCSI mode page 15
func (d Drive) MTCompression(enable int32) error { return d.mtioctop(MTCOMPRESSION, enable) }

// MTSwitchPart changes the active tape partition
func (d Drive) MTSwitchPart(part int32) error { return d.mtioctop(MTSETPART, part) }

// MTMkPart formats the tape with one or two partitions
func (d Drive) MTMkPart(partitions int32) error { return d.mtioctop(MTMKPART, partitions) }

// MTSetOptions sets drive options (using MTSETDRVBUFFER)
func (d Drive) MTSetOptions(op int32) error { return d.mtioctop(MTSETDRVBUFFER, op) }

// MTGetPos gets the current tape position
func (d Drive) MTGetPos() (pos int64, err error) {
	err = ioctl(d.fd, MTIOCPOS, uintptr(unsafe.Pointer(&pos)))
	return
}

// MTGetStatus gets tape status information
func (d Drive) MTGetStatus() (status MtStatus, err error) {
	err = ioctl(d.fd, MTIOCGET, uintptr(unsafe.Pointer(&status)))
	return
}

// mtioctop executes magnetic tape operation commands internally
// op: operation type, count: operation count
func (d Drive) mtioctop(op int16, count int32) (err error) {
	return ioctl(d.fd, MTIOCTOP, uintptr(unsafe.Pointer(&MtOp{op: op, count: count})))
}

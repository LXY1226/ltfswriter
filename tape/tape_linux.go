package tape

import (
	"unsafe"
)

const (

	//MTIOCTOP do a mag tape op
	MTIOCTOP = 0x40086d01
	//MTIOCGET get tape status
	MTIOCGET = 0x80306d02
	//MTIOCPOS get tape position
	MTIOCPOS = 0x80086d03

	//MTRESET +reset drive in case of problems
	MTRESET = 0
	//MTFSF Forward space over FileMark, position at first record of next file
	MTFSF = 1
	//MTBSF Backward space FileMark (position before FM)
	MTBSF = 2
	//MTFSR Forward space record
	MTFSR = 3
	//MTBSR Backward space record
	MTBSR = 4
	//MTWEOF Write an end-of-file record (mark)
	MTWEOF = 5
	//MTREW Rewind
	MTREW = 6
	//MTOFFL Rewind and put the drive offline (eject?)
	MTOFFL = 7
	//MTNOP No op, set status only (read with MTIOCGET)
	MTNOP = 8
	//MTRETEN Retension tape
	MTRETEN = 9
	//MTBSFM +backward space FileMark, position at FM
	MTBSFM = 10
	//MTFSFM +forward space FileMark, position at FM
	MTFSFM = 11
	//MTEOM Goto end of recorded media (for appending files)
	//MTEOM positions after the last FM, ready for appending another file
	MTEOM = 12
	//MTERASE Erase tape -- be careful!
	MTERASE = 13

	//MTRAS1 Run self test 1 (nondestructive)
	MTRAS1 = 14
	//MTRAS2 Run self test 2 (destructive)
	MTRAS2 = 15
	//MTRAS3 Reserved for self test 3
	MTRAS3 = 16

	//MTSETBLK Set block length (SCSI)
	MTSETBLK = 20
	//MTSETDENSITY Set tape density (SCSI)
	MTSETDENSITY = 21
	//MTSEEK Seek to block (Tandberg, etc.)
	MTSEEK = 22
	//MTTELL Tell block (Tandberg, etc.)
	MTTELL = 23
	//MTSETDRVBUFFER Set the drive buffering according to SCSI-2
	//Ordinary buffered operation with code 1
	MTSETDRVBUFFER = 24

	//MTFSS Space forward over setmarks
	MTFSS = 25
	//MTBSS Space backward over setmarks
	MTBSS = 26
	//MTWSM Write setmarks
	MTWSM = 27

	//MTLOCK Lock the drive door
	MTLOCK = 28
	//MTUNLOCK Unlock the drive door
	MTUNLOCK = 29
	//MTLOAD Execute the SCSI load command
	MTLOAD = 30
	//MTUNLOAD Execute the SCSI unload command
	MTUNLOAD = 31
	//MTCOMPRESSION Control compression with SCSI mode page 15
	MTCOMPRESSION = 32
	//MTSETPART Change the active tape partition
	MTSETPART = 33
	//MTMKPART Format the tape with one or two partitions
	MTMKPART = 34

	//MTISUNKNOWN unknown
	MTISUNKNOWN = 0x01
	//MTISQIC02 Generic QIC-02 tape streamer
	MTISQIC02 = 0x02
	//MTISWT5150 Wangtek 5150EQ, QIC-150, QIC-02
	MTISWT5150 = 0x03
	//MTISARCHIVE5945L2 Archive 5945L-2, QIC-24, QIC-02?
	MTISARCHIVE5945L2 = 0x04
	//MTISCMSJ500 CMS Jumbo 500 (QIC-02?)
	MTISCMSJ500 = 0x05
	//MTISTDC3610 Tandberg 6310, QIC-24
	MTISTDC3610 = 0x06
	//MTISARCHIVEVP60I Archive VP60i, QIC-02
	MTISARCHIVEVP60I = 0x07
	//MTISARCHIVE2150L Archive Viper 2150L
	MTISARCHIVE2150L = 0x08
	//MTISARCHIVE2060L Archive Viper 2060L
	MTISARCHIVE2060L = 0x09
	//MTISARCHIVESC499 Archive SC-499 QIC-36 controller
	MTISARCHIVESC499 = 0x0A
	//MTISQIC02ALLFEATURES Generic QIC-02 with all features
	MTISQIC02ALLFEATURES = 0x0F
	//MTISWT5099EEN24 Wangtek 5099-een24, 60MB, QIC-24
	MTISWT5099EEN24 = 0x11
	//MTISTEACMT2ST Teac MT-2ST 155mb drive, Teac DC-1 card (Wangtek type)
	MTISTEACMT2ST = 0x12
	//MTISEVEREXFT40A Everex FT40A (QIC-40)
	MTISEVEREXFT40A = 0x32
	//MTISDDS1 DDS device without partitions
	MTISDDS1 = 0x51
	//MTISDDS2 DDS device with partitions
	MTISDDS2 = 0x52
	//MTISONSTREAMSC OnStream SCSI tape drives (SC-x0) and SCSI emulated (DI, DP, USB)
	MTISONSTREAMSC = 0x61
	//MTISSCSI1 Generic ANSI SCSI-1 tape unit
	MTISSCSI1 = 0x71
	//MTISSCSI2 Generic ANSI SCSI-2 tape unit
	MTISSCSI2 = 0x72

	//MTISFTAPEFLAG QIC-40/80/3010/3020 ftape supported drives
	//20bit vendor ID + 0x800000 (see vendors.h in ftape distribution)
	MTISFTAPEFLAG = 0x800000

	//MTSTBLKSIZESHIFT blocksize shift
	MTSTBLKSIZESHIFT = 0
	//MTSTBLKSIZEMASK blocksize mask
	MTSTBLKSIZEMASK = 0xffffff
	//MTSTDENSITYSHIFT density shift
	MTSTDENSITYSHIFT = 24
	//MTSTDENSITYMASK density mask
	MTSTDENSITYMASK = 0xff000000

	//MTSTSOFTERRSHIFT soft error shift
	MTSTSOFTERRSHIFT = 0
	//MTSTSOFTERRMASK soft error mask
	MTSTSOFTERRMASK = 0xffff

	//MTSTOPTIONS MTSETDRVBUFFER options
	MTSTOPTIONS = 0xf0000000
	//MTSTBOOLEANS MTSETDRVBUFFER booleans
	MTSTBOOLEANS = 0x10000000
	//MTSTSETBOOLEANS MTSETDRVBUFFER set booleans
	MTSTSETBOOLEANS = 0x30000000
	//MTSTCLEARBOOLEANS MTSETDRVBUFFER clear booleans
	MTSTCLEARBOOLEANS = 0x40000000
	//MTSTWRITETHRESHOLD MTSETDRVBUFFER write threshold
	MTSTWRITETHRESHOLD = 0x20000000
	//MTSTDEFBLKSIZE MTSETDRVBUFFER default blocksize
	MTSTDEFBLKSIZE = 0x50000000
	//MTSTDEFOPTIONS MTSETDRVBUFFER default options
	MTSTDEFOPTIONS = 0x60000000
	//MTSTTIMEOUTS timeouts
	MTSTTIMEOUTS = 0x70000000
	//MTSTSETTIMEOUT set timeout
	MTSTSETTIMEOUT = (MTSTTIMEOUTS | 0x000000)
	//MTSTSETLONGTIMEOUT set long timeout
	MTSTSETLONGTIMEOUT = (MTSTTIMEOUTS | 0x100000)
	//MTSTSETCLN set cln
	MTSTSETCLN = 0x80000000

	//MTSTBUFFERWRITES buffered writes
	MTSTBUFFERWRITES = 0x1
	//MTSTASYNCWRITES async writes
	MTSTASYNCWRITES = 0x2
	//MTSTREADAHEAD read ahead
	MTSTREADAHEAD = 0x4
	//MTSTDEBUGGING debugging
	MTSTDEBUGGING = 0x8
	//MTSTTWOFM write two filemarks
	MTSTTWOFM = 0x10
	//MTSTFASTMTEOM send MTEOM directly to drive
	MTSTFASTMTEOM = 0x20
	//MTSTAUTOLOCK auto lock
	MTSTAUTOLOCK = 0x40
	//MTSTDEFWRITES apply settings to drive defaults
	MTSTDEFWRITES = 0x80
	//MTSTCANBSR correct readahaead backspace position
	MTSTCANBSR = 0x100
	//MTSTNOBLKLIMS dont use READ BLOCK LIMITS
	MTSTNOBLKLIMS = 0x200
	//MTSTCANPARTITIONS enable partitions
	MTSTCANPARTITIONS = 0x400
	//MTSTSCSI2LOGICAL use logical block addresses
	MTSTSCSI2LOGICAL = 0x800
	//MTSTSYSV sysv
	MTSTSYSV = 0x1000
	//MTSTNOWAIT no wait
	MTSTNOWAIT = 0x2000
	//MTSTSILI SILI
	MTSTSILI = 0x4000
	//MTSTNOWAITEOF nowait_filemark
	MTSTNOWAITEOF = 0x8000

	//MTSTCLEARDEFAULT clear default
	MTSTCLEARDEFAULT = 0xfffff
	//MTSTDEFDENSITY default density
	MTSTDEFDENSITY = (MTSTDEFOPTIONS | 0x100000)
	//MTSTDEFCOMPRESSION default compression
	MTSTDEFCOMPRESSION = (MTSTDEFOPTIONS | 0x200000)
	//MTSTDEFDRVBUFFER default buffering
	MTSTDEFDRVBUFFER = (MTSTDEFOPTIONS | 0x300000)

	//MTSTHPLOADEROFFSET arguments for the special HP changer load command
	MTSTHPLOADEROFFSET = 10000
)

// MtOp is structure for MTIOCTOP - magnetic tape operation command
type MtOp struct {
	// Operation ID
	op int16

	// Padding to match C structures
	_ int16

	// Operation count
	count int32
}

func (d Drive) MTSetBuffer(size int32) error  { return d.mtioctop(MTSETDRVBUFFER, size) }
func (d Drive) MTSetOptions(op int32) error   { return d.mtioctop(MTSETDRVBUFFER, op) }
func (d Drive) MTSwitchPart(part int32) error { return d.mtioctop(MTSETPART, part) }
func (d Drive) MTSeek(pos int32) error        { return d.mtioctop(MTSEEK, pos) }
func (d Drive) MTSetBlock(size int32) error   { return d.mtioctop(MTSETBLK, size) }
func (d Drive) MTGetPos() (pos int64, err error) {
	err = ioctl(d.fd, MTIOCPOS, uintptr(unsafe.Pointer(&pos)))
	return
}

func (d Drive) mtioctop(op int16, count int32) (err error) {
	return ioctl(d.fd, MTIOCTOP, uintptr(unsafe.Pointer(&MtOp{op: op, count: count})))
}

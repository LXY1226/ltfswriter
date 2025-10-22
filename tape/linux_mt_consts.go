//go:build linux

package tape

const (
	MTIOCTOP = 0x40086d01 // Do a mag tape op
	MTIOCGET = 0x80306d02 // Get tape status
	MTIOCPOS = 0x80086d03 // Get tape position
)

/* Magnetic Tape operations [Not all operations supported by all drivers]: */
// #define MTRESET 0	/* +reset drive in case of problems */
// #define MTFSF	1	/* forward space over FileMark,
// 			 * position at first record of next file
// 			 */
// #define MTBSF	2	/* backward space FileMark (position before FM) */
// #define MTFSR	3	/* forward space record */
// #define MTBSR	4	/* backward space record */
// #define MTWEOF	5	/* write an end-of-file record (mark) */
// #define MTREW	6	/* rewind */
// #define MTOFFL	7	/* rewind and put the drive offline (eject?) */
// #define MTNOP	8	/* no op, set status only (read with MTIOCGET) */
// #define MTRETEN 9	/* retension tape */
// #define MTBSFM	10	/* +backward space FileMark, position at FM */
// #define MTFSFM  11	/* +forward space FileMark, position at FM */
// #define MTEOM	12	/* goto end of recorded media (for appending files).
// 			 * MTEOM positions after the last FM, ready for
// 			 * appending another file.
// 			 */
// #define MTERASE 13	/* erase tape -- be careful! */

// #define MTRAS1  14	/* run self test 1 (nondestructive) */
// #define MTRAS2	15	/* run self test 2 (destructive) */
// #define MTRAS3  16	/* reserved for self test 3 */

// #define MTSETBLK 20	/* set block length (SCSI) */
// #define MTSETDENSITY 21	/* set tape density (SCSI) */
// #define MTSEEK	22	/* seek to block (Tandberg, etc.) */
// #define MTTELL	23	/* tell block (Tandberg, etc.) */
// #define MTSETDRVBUFFER 24 /* set the drive buffering according to SCSI-2 */
// 			/* ordinary buffered operation with code 1 */
// #define MTFSS	25	/* space forward over setmarks */
// #define MTBSS	26	/* space backward over setmarks */
// #define MTWSM	27	/* write setmarks */

// #define MTLOCK  28	/* lock the drive door */
// #define MTUNLOCK 29	/* unlock the drive door */
// #define MTLOAD  30	/* execute the SCSI load command */
// #define MTUNLOAD 31	/* execute the SCSI unload command */
// #define MTCOMPRESSION 32/* control compression with SCSI mode page 15 */
// #define MTSETPART 33	/* Change the active tape partition */
// #define MTMKPART  34	/* Format the tape with one or two partitions */
// #define MTWEOFI	35	/* write an end-of-file record (mark) in immediate mode */
const (
	MTRESET = 0  // Reset drive in case of problems
	MTFSF   = 1  // Forward space over FileMark, position at first record of next file
	MTBSF   = 2  // Backward space FileMark (position before FM)
	MTFSR   = 3  // Forward space record
	MTBSR   = 4  // Backward space record
	MTWEOF  = 5  // Write an end-of-file record (mark)
	MTREW   = 6  // Rewind
	MTOFFL  = 7  // Rewind and put the drive offline (eject?)
	MTNOP   = 8  // No op, set status only (read with MTIOCGET)
	MTRETEN = 9  // Retension tape
	MTBSFM  = 10 // Backward space FileMark, position at FM
	MTFSFM  = 11 // Forward space FileMark, position at FM
	MTEOM   = 12 // Goto end of recorded media (for appending files), positions after the last FM
	MTERASE = 13 // Erase tape -- be careful!

	MTRAS1 = 14 // Run self test 1 (nondestructive)
	MTRAS2 = 15 // Run self test 2 (destructive)
	MTRAS3 = 16 // Reserved for self test 3

	MTSETBLK       = 20 // Set block length (SCSI)
	MTSETDENSITY   = 21 // Set tape density (SCSI)
	MTSEEK         = 22 // Seek to block (Tandberg, etc.)
	MTTELL         = 23 // Tell block (Tandberg, etc.)
	MTSETDRVBUFFER = 24 // Set the drive buffering according to SCSI-2

	MTFSS = 25 // Space forward over setmarks
	MTBSS = 26 // Space backward over setmarks
	MTWSM = 27 // Write setmarks

	MTLOCK        = 28 // Lock the drive door
	MTUNLOCK      = 29 // Unlock the drive door
	MTLOAD        = 30 // Execute the SCSI load command
	MTUNLOAD      = 31 // Execute the SCSI unload command
	MTCOMPRESSION = 32 // Control compression with SCSI mode page 15
	MTSETPART     = 33 // Change the active tape partition
	MTMKPART      = 34 // Format the tape with one or two partitions
	MTWEOFI       = 35 // Write an end-of-file record (mark) in immediate mode
)

const (
	MTISUNKNOWN          = 0x01 // Unknown
	MTISQIC02            = 0x02 // Generic QIC-02 tape streamer
	MTISWT5150           = 0x03 // Wangtek 5150EQ, QIC-150, QIC-02
	MTISARCHIVE5945L2    = 0x04 // Archive 5945L-2, QIC-24, QIC-02?
	MTISCMSJ500          = 0x05 // CMS Jumbo 500 (QIC-02?)
	MTISTDC3610          = 0x06 // Tandberg 6310, QIC-24
	MTISARCHIVEVP60I     = 0x07 // Archive VP60i, QIC-02
	MTISARCHIVE2150L     = 0x08 // Archive Viper 2150L
	MTISARCHIVE2060L     = 0x09 // Archive Viper 2060L
	MTISARCHIVESC499     = 0x0A // Archive SC-499 QIC-36 controller
	MTISQIC02ALLFEATURES = 0x0F // Generic QIC-02 with all features
	MTISWT5099EEN24      = 0x11 // Wangtek 5099-een24, 60MB, QIC-24
	MTISTEACMT2ST        = 0x12 // Teac MT-2ST 155mb drive, Teac DC-1 card (Wangtek type)
	MTISEVEREXFT40A      = 0x32 // Everex FT40A (QIC-40)
	MTISDDS1             = 0x51 // DDS device without partitions
	MTISDDS2             = 0x52 // DDS device with partitions
	MTISONSTREAMSC       = 0x61 // OnStream SCSI tape drives (SC-x0) and SCSI emulated (DI, DP, USB)
	MTISSCSI1            = 0x71 // Generic ANSI SCSI-1 tape unit
	MTISSCSI2            = 0x72 // Generic ANSI SCSI-2 tape unit

	MTISFTAPEFLAG = 0x800000 // QIC-40/80/3010/3020 ftape supported drives (20bit vendor ID + 0x800000)

	MTSTBLKSIZESHIFT = 0          // Blocksize shift
	MTSTBLKSIZEMASK  = 0xffffff   // Blocksize mask
	MTSTDENSITYSHIFT = 24         // Density shift
	MTSTDENSITYMASK  = 0xff000000 // Density mask

	MTSTSOFTERRSHIFT = 0      // Soft error shift
	MTSTSOFTERRMASK  = 0xffff // Soft error mask

	MTSTOPTIONS        = 0xf0000000                // MTSETDRVBUFFER options
	MTSTBOOLEANS       = 0x10000000                // MTSETDRVBUFFER booleans
	MTSTSETBOOLEANS    = 0x30000000                // MTSETDRVBUFFER set booleans
	MTSTCLEARBOOLEANS  = 0x40000000                // MTSETDRVBUFFER clear booleans
	MTSTWRITETHRESHOLD = 0x20000000                // MTSETDRVBUFFER write threshold
	MTSTDEFBLKSIZE     = 0x50000000                // MTSETDRVBUFFER default blocksize
	MTSTDEFOPTIONS     = 0x60000000                // MTSETDRVBUFFER default options
	MTSTTIMEOUTS       = 0x70000000                // Timeouts
	MTSTSETTIMEOUT     = (MTSTTIMEOUTS | 0x000000) // Set timeout
	MTSTSETLONGTIMEOUT = (MTSTTIMEOUTS | 0x100000) // Set long timeout
	MTSTSETCLN         = 0x80000000                // Set cln

	MTSTBUFFERWRITES  = 0x1    // Buffered writes
	MTSTASYNCWRITES   = 0x2    // Async writes
	MTSTREADAHEAD     = 0x4    // Read ahead
	MTSTDEBUGGING     = 0x8    // Debugging
	MTSTTWOFM         = 0x10   // Write two filemarks
	MTSTFASTMTEOM     = 0x20   // Send MTEOM directly to drive
	MTSTAUTOLOCK      = 0x40   // Auto lock
	MTSTDEFWRITES     = 0x80   // Apply settings to drive defaults
	MTSTCANBSR        = 0x100  // Correct readahead backspace position
	MTSTNOBLKLIMS     = 0x200  // Don't use READ BLOCK LIMITS
	MTSTCANPARTITIONS = 0x400  // Enable partitions
	MTSTSCSI2LOGICAL  = 0x800  // Use logical block addresses
	MTSTSYSV          = 0x1000 // SysV
	MTSTNOWAIT        = 0x2000 // No wait
	MTSTSILI          = 0x4000 // SILI
	MTSTNOWAITEOF     = 0x8000 // No wait filemark

	MTSTCLEARDEFAULT   = 0xfffff                     // Clear default
	MTSTDEFDENSITY     = (MTSTDEFOPTIONS | 0x100000) // Default density
	MTSTDEFCOMPRESSION = (MTSTDEFOPTIONS | 0x200000) // Default compression
	MTSTDEFDRVBUFFER   = (MTSTDEFOPTIONS | 0x300000) // Default buffering

	MTSTHPLOADEROFFSET = 10000 // Arguments for the special HP changer load command
)

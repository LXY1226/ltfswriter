package tape

import (
	"log"
	"syscall"
	"unsafe"
)

type Drive struct {
	fd uintptr
}

func Open(path string) (*Drive, error) {
	fd, err := syscall.Open(path, syscall.O_RDWR|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	return &Drive{fd: uintptr(fd)}, nil
}

//typedef struct sg_io_hdr {
//      int interface_id;       		/* [i] 'S' for SCSI generic (required) */
//		int dxfer_direction;    		/* [i] data transfer direction  */
//		unsigned char cmd_len;  		/* [i] SCSI command length */
//		unsigned char mx_sb_len;		/* [i] max length to write to sbp */
//		unsigned short iovec_count;     /* [i] 0 implies no sgat list */
//		unsigned int dxfer_len; 		/* [i] byte count of data transfer */
//		/* dxferp points to data transfer memory or scatter gather list */
//		void __user *dxferp;    		/* [i], [device -> *i or *i -> device] */
//		unsigned char __user *cmdp;		/* [i], [*i] points to command to perform */
//		void __user *sbp;       		/* [i], [*o] points to sense_buffer memory */
//		unsigned int timeout;   		/* [i] MAX_UINT->no timeout (unit: millisec) */
//		unsigned int flags;     		/* [i] 0 -> default, see SG_FLAG... */
//		int pack_id;            		/* [i->o] unused internally (normally) */
//		void __user *usr_ptr;   		/* [i->o] unused internally */
//		unsigned char status;       	/* [o] scsi status */
//		unsigned char masked_status;	/* [o] shifted, masked scsi status */
//		unsigned char msg_status;		/* [o] messaging level data (optional) */
//		unsigned char sb_len_wr; 		/* [o] byte count actually written to sbp */
//		unsigned short host_status; 	/* [o] errors from host adapter */
//		unsigned short driver_status;   /* [o] errors from software driver */
//		int resid;         		        /* [o] dxfer_len - actual_transferred */
//		/* unit may be nanoseconds after SG_SET_GET_EXTENDED ioctl use */
//		unsigned int duration;  		/* [o] time taken by cmd (unit: millisec) */
//		unsigned int info;      		/* [o] auxiliary information */
//} sg_io_hdr_t;

type sgioHdr struct {
	interfaceID    int32
	dxferDirection int32
	cmdLen         uint8
	mxSbLen        uint8
	iovecCount     uint16
	dxferLen       uint32
	dxferp         uintptr
	cmdp           uintptr
	sbp            uintptr
	timeout        uint32
	flags          uint32
	packID         int32
	usrPtr         uintptr
	status         uint8
	maskedStatus   uint8
	msgStatus      uint8
	sbLenWr        uint8
	hostStatus     uint16
	driverStatus   uint16
	resid          int32

	duration uint
	info     uint
}

func newScsiCmd(cmd []byte, timeout uint32) sgioHdr {
	sbp := make(senseError, senseBufferSize)
	return sgioHdr{
		interfaceID: 'S',

		cmdLen:  uint8(len(cmd)),
		cmdp:    uintptr(unsafe.Pointer(&cmd[0])), // unsafe.SliceData(cmd)?
		sbp:     uintptr(unsafe.Pointer(&sbp[0])),
		mxSbLen: senseBufferSize,

		timeout: timeout,
	}
}

func (d Drive) scsiCmd(cmd []byte, timeout uint32) error {
	hdr := newScsiCmd(cmd, timeout)
	hdr.dxferDirection = -1 // SG_DXFER_NONE
	return scsi(d.fd, &hdr)
}

func (d Drive) scsiRead(cmd []byte, recvLen uint32, timeout uint32) ([]byte, error) {
	buf := make([]byte, recvLen)
	hdr := newScsiCmd(cmd, timeout)
	hdr.dxferDirection = -3 // SG_DXFER_FROM_DEV
	hdr.dxferLen = recvLen
	hdr.dxferp = uintptr(unsafe.Pointer(&buf[0]))
	err := scsi(d.fd, &hdr)
	//log.Println(hdr.resid)
	return buf[:hdr.resid], err
}

func (d Drive) scsiWrite(cmd, buf []byte, timeout uint32) error {
	hdr := newScsiCmd(cmd, timeout)
	hdr.dxferDirection = -2 // SG_DXFER_TO_DEV
	hdr.dxferLen = uint32(len(buf))
	hdr.dxferp = uintptr(unsafe.Pointer(&buf[0]))
	return scsi(d.fd, &hdr)
}

func scsi(fd uintptr, hdr *sgioHdr) error {
	err := sgIO(fd, hdr)
	if hdr.sbLenWr != 0 && *(*byte)(unsafe.Pointer(hdr.sbp)) != 0 {
		log.Printf("scsi: cmd=%x status=%x host=%d driver=%d resid=%d sb=%x\n",
			unsafe.Slice((*byte)(unsafe.Pointer(hdr.cmdp)), hdr.cmdLen),
			hdr.status, hdr.hostStatus, hdr.driverStatus, hdr.resid,
			unsafe.Slice((*byte)(unsafe.Pointer(hdr.sbp)), hdr.sbLenWr))

		return senseError(unsafe.Slice((*byte)(unsafe.Pointer(hdr.sbp)), hdr.sbLenWr))
	}
	return err
}

func sgIO(fd uintptr, arg *sgioHdr) (err error) {
	_, _, e1 := syscall.Syscall(syscall.SYS_IOCTL, fd, 0x2285, uintptr(unsafe.Pointer(arg))) // SG_IO
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

//func rawIoctl(fd uintptr, cmd, arg uintptr) error {
//	return 0, nil
//}

var (
	errEAGAIN error = syscall.EAGAIN
	errEINVAL error = syscall.EINVAL
	errENOENT error = syscall.ENOENT
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.EAGAIN:
		return errEAGAIN
	case syscall.EINVAL:
		return errEINVAL
	case syscall.ENOENT:
		return errENOENT
	}
	return e
}

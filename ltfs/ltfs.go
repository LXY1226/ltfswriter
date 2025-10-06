package ltfs

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/LXY1226/ltfswriter/tape"
	"io"
	"log"
)

type Label struct {
	XMLName    xml.Name `xml:"ltfslabel"`
	Version    string   `xml:"version,attr"`
	Creator    string   `xml:"creator"`
	Formattime string   `xml:"formattime"`
	Volumeuuid string   `xml:"volumeuuid"`
	Location   struct {
		Partition string `xml:"partition"`
	} `xml:"location"`
	Partitions struct {
		Index string `xml:"index"`
		Data  string `xml:"data"`
	} `xml:"partitions"`
	Blocksize   string `xml:"blocksize"`
	Compression string `xml:"compression"`
}

type Index struct {
	XMLName                    xml.Name  `xml:"ltfsindex"`
	Version                    string    `xml:"version,attr"`
	Creator                    string    `xml:"creator"`
	VolumeUUID                 string    `xml:"volumeuuid"`
	GenerationNumber           int       `xml:"generationnumber"`
	UpdateTime                 string    `xml:"updatetime"`
	Location                   Location  `xml:"location"`
	PreviousGenerationLocation Location  `xml:"previousgenerationlocation"`
	AllowPolicyUpdate          bool      `xml:"allowpolicyupdate"`
	VolumeLockState            string    `xml:"volumelockstate"`
	HighestFileUID             int       `xml:"highestfileuid"`
	Directory                  Directory `xml:"directory"`
}

type Location struct {
	Partition  string `xml:"partition"`
	StartBlock int    `xml:"startblock"`
}

type Directory struct {
	Name         string   `xml:"name"`
	ReadOnly     bool     `xml:"readonly"`
	CreationTime string   `xml:"creationtime"`
	ChangeTime   string   `xml:"changetime"`
	ModifyTime   string   `xml:"modifytime"`
	AccessTime   string   `xml:"accesstime"`
	BackupTime   string   `xml:"backuptime"`
	FileUID      int      `xml:"fileuid"`
	Contents     Contents `xml:"contents"`
}

type Contents struct {
	Files       []File      `xml:"file"`
	Directories []Directory `xml:"directory"`
}

type File struct {
	Name               string              `xml:"name"`
	Length             int64               `xml:"length"`
	ReadOnly           bool                `xml:"readonly"`
	OpenForWrite       bool                `xml:"openforwrite"`
	CreationTime       string              `xml:"creationtime"`
	ChangeTime         string              `xml:"changetime"`
	ModifyTime         string              `xml:"modifytime"`
	AccessTime         string              `xml:"accesstime"`
	BackupTime         string              `xml:"backuptime"`
	FileUID            int                 `xml:"fileuid"`
	ExtendedAttributes *ExtendedAttributes `xml:"extendedattributes"`
	ExtentInfo         *ExtentInfo         `xml:"extentinfo"`
}

type ExtendedAttributes struct {
	XAttrs []XAttr `xml:"xattr"`
}

type XAttr struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}

type ExtentInfo struct {
	Extents []Extent `xml:"extent"`
}

type Extent struct {
	FileOffset int64  `xml:"fileoffset"`
	Partition  string `xml:"partition"`
	StartBlock int64  `xml:"startblock"`
	ByteOffset int64  `xml:"byteoffset"`
	ByteCount  int64  `xml:"bytecount"`
}

type Volume struct {
	Vol1Label   tape.VOL1Label
	Label       Label
	LatestIndex Index
	Indexes     []struct {
		Position tape.PositionData
		Index    Index
	}
}

// Open read LTFSVolume from drive
// TODO interface
func Open(drive *tape.Drive) (*Volume, error) {
	aVol1, aLabel, aIndex, err := readPartHead(drive, 0)
	if err != nil {
		panic(err)
	}
	vol := new(Volume)
	vol.Vol1Label, err = tape.ParseVol1Label(aVol1)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(aLabel, &vol.Label)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(aIndex, &vol.LatestIndex)
	for {
		//pos, err := drive.MTGetPos()
		//if err == nil {
		//	panic(err)
		//}
		//status, err := drive.MTGetStatus()
		dat, err := drive.MTReadShortFile()
		if err != io.EOF {
			if err := drive.MTBSFM(1); err != nil {
				panic(err)
			}
			log.Println(err)
			i := bytes.Index(dat, []byte(`<ltfsindex`))
			if i == -1 {
				return vol, nil
			}
			dat = make([]byte, 16<<20) // big buf!
			n, err := drive.MTReadFull(dat)
			if err != io.EOF {
				if err := drive.MTBSFM(1); err != nil {
					panic(err)
				}
				return vol, fmt.Errorf("index too big")
			}
			dat = dat[:n]
		}
		var idx Index
		err = xml.Unmarshal(dat, &idx)
		if err != nil {
			return vol, err
		}
		vol.Indexes = append(vol.Indexes, struct {
			Position tape.PositionData
			Index    Index
		}{Position: tape.PositionData{
			Partition: 0,
			Block:     0,
			File:      0,
		}, Index: idx})
		if idx.GenerationNumber > vol.LatestIndex.GenerationNumber {
			vol.LatestIndex = idx
		}
	}
}

func PartToSCSIPart(c byte) int32 {
	return int32(c - 'A')
}

// readPartHead read and check partition has valid starting
func readPartHead(drive *tape.Drive, part int32) (vol1, label, index []byte, err error) {
	err = drive.MTSeek(0)
	if err != nil {
		panic(err)
	}
	err = drive.MTSwitchPart(part)
	if err != nil {
		panic(err)
	}
	vol1 = drive.MTMustReadShortFile()
	label = drive.MTMustReadShortFile()
	index = drive.MTMustReadShortFile()
	return
}

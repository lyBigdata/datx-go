package datx

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func NewBaseStation(fn string) (*BaseStation, error) {

	db := &BaseStation{}

	if err := db.load(fn); err != nil {
		return nil, err
	}

	return db, nil
}

type BaseStation struct {
	file *os.File

	index []byte
	data  []byte
}

func (db *BaseStation) load(fn string) error {
	var err error
	db.file, err = os.Open(fn)
	if err != nil {
		return err
	}

	b4 := make([]byte, 4)

	db.file.Read(b4)

	off := int(binary.BigEndian.Uint32(b4))
	db.file.Seek(262148, 0)

	db.index = make([]byte, off-262148-262144)
	db.file.Read(db.index)

	db.data, err = ioutil.ReadAll(db.file)
	if err != nil {
		return err
	}
	//	fmt.Println(len(db.data))
	return nil
}

func (db *BaseStation) Find(s string) ([]string, error) {

	ipv := net.ParseIP(s)
	if ipv == nil {
		return nil, fmt.Errorf("%s", "ip format error.")
	}

	low := 0
	high := int(len(db.index) / 13)
	mid := 0

	val := binary.BigEndian.Uint32(ipv.To4())

	for low <= high {
		mid = int((low + high) / 2)
		pos := mid * 13

		start := binary.BigEndian.Uint32(db.index[pos : pos+4])
		end := binary.BigEndian.Uint32(db.index[pos+4 : pos+8])

		if val < start {
			high = mid - 1
		} else if val > end {
			low = mid + 1
		} else {

			off := int(binary.LittleEndian.Uint32(db.index[pos+8 : pos+12]))

			return strings.Split(string(db.data[off:off+int(db.index[pos+12])]), "\t"), nil
		}
	}
	return nil, nil
}

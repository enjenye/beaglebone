//Beaglebone support package.

package beaglebone

import (
	"errors"
	"syscall"
)

type Bone int

func NewBone() (Bone, error) {
	b, err := syscall.Open("/dev/logibone_mem", syscall.O_RDWR|syscall.O_SYNC, 777)
	if err != nil {
		return 0, err
	}
	return Bone(b), nil
}

func (b Bone) EndBone() {
	syscall.Close(int(b))
	return
}

func (b Bone) ReadInt16(addr int) (uint16, error) {
	data := make([]byte, 2)
	count, err := syscall.Pread(int(b), data, int64(addr))
	if err != nil {
		return 0, err
	}
	if count != 2 {
		return 0, errors.New("wrong number of bytes read")
	}
	val := (uint16(data[1]) << 8) | uint16(data[0])
	return val, nil
}

func (b Bone) WriteInt16(addr int, val uint16) error {
	data := make([]byte, 2)
	data[0] = byte((val) & 0x00FF)
	data[1] = byte((val >> 8) & 0x00FF)
	count, err := syscall.Pwrite(int(b), data, int64(addr))
	if err != nil {
		return err
	}
	if count != 2 {
		return errors.New("wrong number of bytes written")
	}
	return nil
}

func (b Bone) ReadInt32(addr int) (uint32, error) {
	msb, err := b.ReadInt16(addr)
	if err != nil {
		return 0, err
	}
	lsb, err := b.ReadInt16(addr + 1)
	if err != nil {
		return 0, err
	}
	val := (uint32(msb) << 16) | uint32(lsb)
	return val, nil
}

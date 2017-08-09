package scel

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"unicode/utf16"
)

type Scel struct {
	Data []byte

	PyTable map[uint16]string
}

func NewScel(data []byte) *Scel {
	return &Scel{
		Data:    data,
		PyTable: make(map[uint16]string),
	}
}

var SougouTag = []byte{0x40, 0x15, 0x00, 0x00, 0x44, 0x43, 0x53, 0x01, 0x01, 0x00, 0x00, 0x00}

func (s *Scel) IsValid() bool {
	data := s.Data

	if !reflect.DeepEqual(data[:12], SougouTag) {
		return true
	}

	return false
}

// UTF-16
func byte2str(data []byte) (string, error) {
	var c uint16
	var result []uint16
	var err error

	r := bytes.NewReader(data)

	for {
		err = binary.Read(r, binary.LittleEndian, &c)
		if err != nil {
			break
		}
		if c != 0 {
			result = append(result, c)
		}
	}

	if err != io.EOF {
		return "", err
	}

	return string(utf16.Decode(result)), nil
}

func (s *Scel) Run() (err error) {
	if ok := s.IsValid(); !ok {
		return errors.New("wrong scel format")
	}

	err = s.genPyTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *Scel) genPyTable() error {
	var index, l uint16
	var err error

	data := s.Data[0x1540:0x2628]

	if !reflect.DeepEqual(data[0:4], []byte{0x9D, 0x01, 0x00, 0x00}) {
		return errors.New("wrong sogou pytable")
	}

	data = data[4:]

	r := bytes.NewReader(data)

	for {
		// index
		err = binary.Read(r, binary.LittleEndian, &index)
		if err != nil {
			break
		}

		// l
		err = binary.Read(r, binary.LittleEndian, &l)
		if err != nil {
			break
		}

		b := make([]byte, l)
		_, err = r.Read(b)
		if err != nil {
			break
		}

		pinyin, err := byte2str(b)
		if err != nil {
			break
		}

		s.PyTable[index] = pinyin

		if err != nil {
			break
		}

	}

	if err != io.EOF {
		return err
	}

	return nil
}

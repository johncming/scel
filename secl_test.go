package scel

import (
	"io/ioutil"
	"testing"
)

func testData(t *testing.T) []byte {
	d, err := ioutil.ReadFile("sample.scel")
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func TestByte2str(t *testing.T) {
	data := testData(t)

	cases := map[string][]byte{
		"军事词汇大全【官方推荐】": data[0x130:0x338],
		"军事": data[0x338:0x540],
		"官方推荐，词库来源于网友上传！": data[0x540:0xd40],
	}

	for expect, kace := range cases {
		res, err := byte2str(kace)
		if err != nil {
			t.Fatal(err)
		}

		if res != expect {
			t.Fail()
		}

	}

}

func TestPyTable(t *testing.T) {
	data := testData(t)

	cases := map[uint16]string{
		390: "zhe",
		204: "min",
	}

	s := &Scel{
		Data:    data,
		PyTable: make(map[uint16]string),
	}

	err := s.genPyTable()
	if err != nil {
		t.Fatal(err)
	}

	for index, py := range cases {
		if py != s.PyTable[index] {
			t.Fail()
		}
	}

}

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const (
	FROM = "testdata/input.txt"
	TO   = "/tmp/temp_sviat_v_guss_hw_07.txt"
)

var tests = []struct {
	id               int
	fromPath, toPath string
	offset, limit    int64
	expected         string
}{
	{1, FROM, TO, 0, 0, "testdata/out_offset0_limit0.txt"},
	{2, FROM, TO, 0, 10, "testdata/out_offset0_limit10.txt"},
	{3, FROM, TO, 0, 1000, "testdata/out_offset0_limit1000.txt"},
	{4, FROM, TO, 0, 10000, "testdata/out_offset0_limit10000.txt"},
	{5, FROM, TO, 100, 1000, "testdata/out_offset100_limit1000.txt"},
	{6, FROM, TO, 6000, 1000, "testdata/out_offset6000_limit1000.txt"},
	{7, "/dev/urandom", TO, 0, 0, ""}, // неизвестна длина (например, /dev/urandom)
	{8, FROM, TO, 10000, 0, ""},       // offset больше, чем размер файла - невалидная ситуация
}

func TestCopy(t *testing.T) {
	for _, test := range tests {
		err := Copy(test.fromPath, test.toPath, test.offset, test.limit)
		if err != nil {
			switch test.id {
			case 7:
				if errors.Is(err, ErrUnsupportedFile) {
					continue
				}
			case 8:
				if errors.Is(err, ErrOffsetExceedsFileSize) {
					continue
				}
			default:
				t.Errorf("%v", err)
			}
		}

		// copied file
		file, err := os.Open(test.toPath)
		if err != nil {
			t.Errorf("%v", err)
		}
		reality, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("%v", err)
		}
		file.Close()

		// expected file
		file, err = os.Open(test.expected)
		if err != nil {
			t.Errorf("%v", err)
		}
		expectation, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("%v", err)
		}
		file.Close()

		err = os.Remove(TO)
		if err != nil {
			fmt.Println(err)
		}

		if !reflect.DeepEqual(reality, expectation) {
			t.Errorf("Copied file doesn't equal to %v", test.fromPath)
		}
	}
}

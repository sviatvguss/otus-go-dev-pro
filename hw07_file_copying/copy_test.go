package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

const FROM = "testdata/input.txt"
const TO = "/tmp/temp.txt"

var tests = []struct {
	fromPath, toPath string
	offset, limit    int64
	expected         string
}{
	{FROM, TO, 0, 0, "testdata/out_offset0_limit0.txt"},
	{FROM, TO, 0, 10, "testdata/out_offset0_limit10.txt"},
	{FROM, TO, 0, 1000, "testdata/out_offset0_limit1000.txt"},
	{FROM, TO, 0, 10000, "testdata/out_offset0_limit10000.txt"},
	{FROM, TO, 100, 1000, "testdata/out_offset100_limit1000.txt"},
	{FROM, TO, 6000, 1000, "testdata/out_offset6000_limit1000.txt"},
}

func TestCopyRegular(t *testing.T) {
	for _, test := range tests {
		err := Copy(test.fromPath, test.toPath, test.offset, test.limit)
		if err != nil {
			t.Errorf("%v", err)
		}

		//copied file
		file, err := os.Open(test.toPath)
		if err != nil {
			t.Errorf("%v", err)
		}
		reality, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("%v", err)
		}
		file.Close()

		//expected file
		file, err = os.Open(test.expected)
		if err != nil {
			t.Errorf("%v", err)
		}
		expectation, err := ioutil.ReadAll(file)
		if err != nil {
			t.Errorf("%v", err)
		}
		file.Close()

		if !reflect.DeepEqual(reality, expectation) {
			t.Errorf("Copied file doesn't equal to %v", test.fromPath)
		}
	}
}

package cfs

import (
	"os"
	"scc/decen"
)

func Open(fileName string, encryptor *decen.AESCBC) (*CryptoFileSystem, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	fs := &CryptoFileSystem{
		f:      f,
		de:     encryptor,
		header: nil,
	}
	err = fs.readHeader()
	if err != nil {
		return nil, err
	}
	return fs, nil
}

package cfs

import (
	"os"
	"scc/decen"
)

func Create(fileName string, encryptor *decen.AESCBC) (*CryptoFileSystem, error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, err
	}
	fs := &CryptoFileSystem{
		f:      f,
		de:     encryptor,
		header: &CryptoFsHeader{},
	}
	fs.genHeader()

	return fs, nil
}

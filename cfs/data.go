package cfs

import (
	"errors"
	"os"
	"scc/decen"
)

const (
	cryptoFsHeaderLen = 1 << 12
	cryptoBlockSize   = 1600
)

var (
	ErrDecode = errors.New("decoding error! maybe password is incorrect")
)

type (
	CryptoFileSystem struct {
		f  *os.File
		de *decen.AESCBC

		header    *CryptoFsHeader
		headerLen int
	}
	CryptoFsHeader struct {
		Modified bool
		Items    []*CryptoFsItem
	}
	CryptoFsItem struct {
		Name      string
		Offset    int
		Len       int
		Removed   bool
		OriginLen int
	}
)

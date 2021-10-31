package cfs

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
)

func (x *CryptoFileSystem) AddFile(filePath string) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var (
		n             int
		b, dest       []byte
		item          *CryptoFsItem
		fileOffset    int
		fileLen       int
		fileCryptoLen int
	)
	b = make([]byte, cryptoBlockSize-x.de.OutputAdd())
	dest = make([]byte, cryptoBlockSize)
	// calc new file offset
	for _, fsItem := range x.header.Items {
		fileOffset += fsItem.Len
	}
	ivbuf := make([]byte, x.de.IVLen())
	for {
		clearByteArray(b)
		clearByteArray(dest)
		n, err = f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if n == 0 {
			break
		}
		if _, err = io.ReadFull(rand.Reader, ivbuf); err != nil {
			return err
		}
		x.de.Encrypt(b, dest, ivbuf)
		_, err = x.f.Seek(int64(4+cryptoFsHeaderLen+(fileOffset)+fileCryptoLen), io.SeekStart)
		if err != nil {
			return err
		}
		_, err = x.f.Write(dest)
		if err != nil {
			return err
		}
		fileLen += n
		fileCryptoLen += cryptoBlockSize
	}
	err = x.f.Sync()
	if err != nil {
		return err
	}
	_, name := path.Split(filePath)
	item = &CryptoFsItem{
		Name:      name,
		Offset:    fileOffset,
		Len:       fileCryptoLen,
		OriginLen: fileLen,
	}
	x.header.Items = append(x.header.Items, item)
	x.header.Modified = true
	x.storeHeader()
	return nil
}

func clearByteArray(a []byte) {
	if len(a) == 0 {
		return
	}
	a[0] = 0
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

package cfs

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

func (x *CryptoFileSystem) genHeader() {
	x.header = &CryptoFsHeader{Modified: true}
}

func (x *CryptoFileSystem) storeHeader() error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}

	h, hl, err := x.serializeHeader()
	if err != nil {
		return err
	}

	_, err = x.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	x.headerLen = hl
	headerLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(headerLenBytes, uint32(x.headerLen))
	_, err = x.f.Write(headerLenBytes)
	if err != nil {
		return err
	}
	_, err = x.f.Write(h)
	if err != nil {
		return err
	}
	x.header.Modified = false
	return err
}

func (x *CryptoFileSystem) serializeHeader() ([]byte, int, error) {
	var (
		buf *bytes.Buffer
	)
	buf = bytes.NewBuffer(make([]byte, 0, cryptoBlockSize))
	enc := gob.NewEncoder(buf)
	err := enc.Encode(x.header)
	if err != nil {
		return nil, 0, err
	}
	hbytes := buf.Bytes()
	//log.Printf("gob buf=%v", hbytes)
	l := x.de.AdjustInputSize(len(hbytes))
	if l > cryptoFsHeaderLen {
		panic("header is more than 1MB! too many files in SCC?")
	}
	input := make([]byte, l)
	copy(input, hbytes)
	dst := make([]byte, cryptoFsHeaderLen)
	ivbuf := make([]byte, x.de.IVLen())
	if _, err := io.ReadFull(rand.Reader, ivbuf); err != nil {
		//log.Fatalln("Unable to get rand data:", err)
		return nil, 0, err
	}
	l = x.de.Encrypt(input, dst, ivbuf)
	//log.Printf("aes buf=%v", dst[:l])
	return dst, len(hbytes), nil
}

func (x *CryptoFileSystem) readHeader() error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	_, err := x.f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	var buf = make([]byte, 4)
	_, err = x.f.Read(buf)
	if err != nil {
		return err
	}
	x.headerLen = int(binary.LittleEndian.Uint32(buf))
	buf = make([]byte, x.de.AdjustInputSize(x.headerLen)+x.de.OutputAdd())
	_, err = x.f.Read(buf)
	if err != nil {
		return err
	}
	dst := make([]byte, x.de.AdjustInputSize(x.headerLen))
	_, err = x.de.Decrypt(buf, dst)
	if err != nil {
		return err
	}
	x.header = &CryptoFsHeader{}
	decoder := gob.NewDecoder(bytes.NewBuffer(dst))
	err = decoder.Decode(x.header)
	if err != nil {
		return ErrDecode
	}
	x.header.Modified = false
	return nil
}

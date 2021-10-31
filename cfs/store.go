package cfs

import "fmt"

func (x *CryptoFileSystem) Store() error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	if x.header == nil {
		return fmt.Errorf("CryptoFsHeader is nil")
	}
	var (
		err error
	)
	// store header
	if x.header.Modified {
		err = x.storeHeader()
		if err != nil {
			return err
		}
	}
	err = x.f.Close()
	return err
}

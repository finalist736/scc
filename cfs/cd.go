package cfs

import "os"

func (x *CryptoFileSystem) CD(dir string) error {
	return os.Chdir(dir)
}

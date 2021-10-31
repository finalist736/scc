package cfs

import "fmt"

func (x *CryptoFileSystem) RM(fileName string) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	for _, item := range x.header.Items {
		if item.Name != fileName {
			continue
		}
		item.Removed = true
	}
	x.header.Modified = true
	return x.Vacuum()
}

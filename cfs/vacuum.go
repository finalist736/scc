package cfs

import (
	"fmt"
	"io"
)

func (x *CryptoFileSystem) Vacuum() error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	var (
		isRemovedFound bool
		buf            = make([]byte, cryptoBlockSize)
		err            error
		fileLen        int
		totalFileLen   int
		newOffset      int
	)
	totalFileLen = 4 + cryptoFsHeaderLen
	for _, item := range x.header.Items {
		totalFileLen += item.Len
	}
	for {
		isRemovedFound = false
		for index, item := range x.header.Items {
			if !item.Removed {
				continue
			}
			totalFileLen -= item.Len
			newOffset = item.Offset
			isRemovedFound = true
			for _, fsItem := range x.header.Items[index+1:] {
				fileLen = 0
				for fileLen < fsItem.Len {

					clearByteArray(buf)
					_, err = x.f.Seek(int64(4+cryptoFsHeaderLen+fsItem.Offset+fileLen), io.SeekStart)
					if err != nil {
						return err
					}
					_, err = x.f.Read(buf)
					if err != nil {
						return err
					}
					_, err = x.f.Seek(int64(4+cryptoFsHeaderLen+newOffset+fileLen), io.SeekStart)
					_, err = x.f.Write(buf)
					if err != nil {
						return err
					}
					fileLen += cryptoBlockSize
				}
				fsItem.Offset = newOffset
				newOffset = newOffset + fsItem.Len
			}
			x.header.Items = append(x.header.Items[:index], x.header.Items[index+1:]...)
			x.header.Modified = true
			break
		}
		if !isRemovedFound {
			break
		}
		err = x.f.Truncate(int64(totalFileLen))
		if err != nil {
			return err
		}
	}

	return nil
}

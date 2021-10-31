package cfs

import (
	"fmt"
	"github.com/chzyer/readline"
	"scc/helpers"
)

func (x *CryptoFileSystem) DU(instance *readline.Instance) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	var (
		size       int
		originSize int
		overSize   int
	)
	for _, item := range x.header.Items {
		if item.Removed {
			overSize += item.Len
		} else {
			size += item.Len
			originSize += item.OriginLen
		}
	}
	str := fmt.Sprintf("crypted(origin) %s(%s) %d B(%d B)\noversize %s\n",
		helpers.Bytes(uint64(size)), helpers.Bytes(uint64(originSize)),
		size, originSize, helpers.Bytes(uint64(overSize)))
	_, err := instance.Write([]byte(str))
	return err
}

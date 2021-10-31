package cfs

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
)

func (x *CryptoFileSystem) PWD(instance *readline.Instance) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	_, err = instance.Write([]byte(fmt.Sprintf("%s\n", wd)))
	return err
}

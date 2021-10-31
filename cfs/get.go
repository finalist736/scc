package cfs

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"path"
)

func (x *CryptoFileSystem) GetFile(instance *readline.Instance, file string) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	var (
		err error
		str string
		wd  string
	)
	for _, item := range x.header.Items {
		if item.Name != file {
			continue
		}

		wd, err = os.Getwd()
		if _, err = os.Stat(path.Join(wd, file)); !os.IsNotExist(err) {
			instance.SetPrompt(fmt.Sprintf("file %s is exists, replace? [y/N] ", file))
			str, err = instance.Readline()
			if err != nil {
				return err
			}
			switch str {
			case "Y", "y", "д", "Д":
			default:
				return nil
			}
		}
		var (
			originBlockLen       = cryptoBlockSize - x.de.OutputAdd()
			input                = make([]byte, cryptoBlockSize)
			output               = make([]byte, originBlockLen)
			cryptoLen, originLen int
			f                    *os.File
			n                    int
		)
		cryptoLen = 0
		originLen = item.OriginLen

		f, err = os.Create(item.Name)
		if err != nil {
			return err
		}
		defer f.Close()
		for cryptoLen < item.Len {
			_, err = x.f.Seek(int64(4+cryptoFsHeaderLen+item.Offset+cryptoLen), io.SeekStart)
			if err != nil {
				return err
			}
			clearByteArray(input)
			clearByteArray(output)

			_, err = x.f.Read(input)
			if err != nil {
				return err
			}
			_, err = x.de.Decrypt(input, output)
			if err != nil {
				return err
			}
			n = min(originLen, originBlockLen)
			_, err = f.Write(output[:n])
			if err != nil {
				return err
			}
			cryptoLen += cryptoBlockSize
			originLen -= n
		}
		return nil
	}
	instance.Write([]byte(fmt.Sprintf("file %s not found\n", file)))
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

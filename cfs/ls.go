package cfs

import (
	"fmt"
	"github.com/chzyer/readline"
	"text/tabwriter"
)

func (x *CryptoFileSystem) LS(instance *readline.Instance) error {
	if x == nil {
		return fmt.Errorf("CryptoFileSystem is nil")
	}
	if x.f == nil {
		return fmt.Errorf("os.File is nil")
	}
	var (
	//result string
	)
	tw := new(tabwriter.Writer)
	tw.Init(instance.Stdout(), 8, 8, 0, '\t', 0)
	fmt.Fprintf(tw, "\n%s\t%s\t%s", "FileName", "Crypted size", "Origin size")
	for _, item := range x.header.Items {
		if item.Removed {
			continue
		}
		fmt.Fprintf(tw, "\n%s\t%d\t%d", item.Name, item.Len, item.OriginLen)
		//result += fmt.Sprintf("%s\t%d\t%d\n", item.Name, item.Len, item.OriginLen)
	}
	tw.Flush()
	//_, err := instance.Write([]byte(fmt.Sprintf("file list:\n%s", result)))
	return nil
}

package termm

import (
	"fmt"
	"github.com/chzyer/readline"
	"io/ioutil"
	"log"
	"os"
	"scc/cfs"
	"scc/decen"
	"strconv"
	"strings"
	"syscall"
)

var (
	cfsFolderName = "/.cfs/"
)

type Terminal struct {
	c *cfs.CryptoFileSystem
}

func (x *Terminal) Stop() {
	log.Println("stopping")
	x.c.Store()
}

func (x *Terminal) Start() error {

	var (
		err         error
		line        string
		projectName string
	)
	uhd, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfsFolder := uhd + cfsFolderName
	err = os.Mkdir(cfsFolder, 0700)
	if err != nil {
		_, ok := err.(*os.PathError)
		if !ok {
			return err
		}
	}
	files, err := ioutil.ReadDir(cfsFolder)
	if err != nil {
		return err
	}
	config := &readline.Config{
		Prompt:         "",
		AutoComplete:   nil,
		UniqueEditLine: false,
	}
	instance, err := readline.NewEx(config)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		instance.Write([]byte("no cfs files detected, lets create one\n"))
		instance.SetPrompt("enter file name without .cfs: ")
		line, err = instance.Readline()
		if err != nil {
			return err
		}
		fileName := line
		projectName = fileName
		instance.SetPrompt("enter key")
		key, err := readline.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}
		fileName = fmt.Sprintf("%s%s.cfs", cfsFolder, fileName)
		encryptor, err := decen.NewAesCbc(key)
		x.c, err = cfs.Create(fileName, encryptor)
		if err != nil {
			return err
		}
		config.AutoComplete = x.c
		instance.SetConfig(config)
	} else {
		for {
			instance.Write([]byte("0) Create New\n"))
			for i, file := range files {
				instance.Write([]byte(fmt.Sprintf("%d) %s\n", i+1, fileNameWithoutExtension(file.Name()))))
			}
			instance.SetPrompt("choice file: ")
			//line, _, err = stdinReader.ReadLine()
			line, err = instance.Readline()
			if err != nil {
				return err
			}
			sel, err := strconv.Atoi(line)
			if err != nil {
				return err
			}
			if sel > len(files) || sel < 0 {
				instance.Write([]byte("incorrect choice, try again\n"))
				continue
			}
			if sel == 0 {
				// TODO create func createNew()
				instance.SetPrompt("enter file name without .cfs: ")
				line, err = instance.Readline()
				if err != nil {
					return err
				}
				fileName := line
				projectName = fileName
				instance.SetPrompt("enter key")
				key, err := readline.ReadPassword(syscall.Stdin)
				if err != nil {
					return err
				}
				fileName = fmt.Sprintf("%s%s.cfs", cfsFolder, fileName)
				encryptor, err := decen.NewAesCbc(key)
				x.c, err = cfs.Create(fileName, encryptor)
				if err != nil {
					return err
				}
				config.AutoComplete = x.c
				instance.SetConfig(config)
			} else {
				projectName = fileNameWithoutExtension(files[sel-1].Name())
				fileName := fmt.Sprintf("%s%s", cfsFolder, files[sel-1].Name())
				instance.Write([]byte("enter key"))
				key, err := readline.ReadPassword(syscall.Stdin)
				if err != nil {
					return err
				}
				encryptor, err := decen.NewAesCbc(key)
				x.c, err = cfs.Open(fileName, encryptor)
				if err != nil {
					if err == cfs.ErrDecode {
						instance.Write([]byte(fmt.Sprintf("error: %s\n", err)))
						continue
					}
					return err
				}
				config.AutoComplete = x.c
				instance.SetConfig(config)
			}
			break
		}
	}
cmdLoop:
	for {
		instance.SetPrompt(fmt.Sprintf("%s > ", projectName))
		line, err = instance.Readline()
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		switch line {
		case "exit":
			break cmdLoop
		case "help":
			instance.Write([]byte(helpString))
		case "du":
			x.c.DU(instance)
		case "ls":
			x.c.LS(instance)
		case "pwd":
			x.c.PWD(instance)
		case "vacuum":
			x.c.Vacuum()
		case "":
		default:
			if strings.HasPrefix(line, "add ") {
				fullPathFile := getFileNameFromCommandLine(line)
				if fullPathFile == "" {
					instance.Write([]byte(fmt.Sprintf("add cmd error type help\n")))
					continue
				}
				err = x.c.AddFile(fullPathFile)
				if err != nil {
					return err
				}
			} else if strings.HasPrefix(line, "rm ") {
				fullPathFile := getFileNameFromCommandLine(line)
				if fullPathFile == "" {
					instance.Write([]byte(fmt.Sprintf("rm cmd error type help\n")))
					continue
				}
				err = x.c.RM(fullPathFile)
				if err != nil {
					instance.Write([]byte(fmt.Sprintf("rm error: %s\n", err)))
				}
			} else if strings.HasPrefix(line, "cd ") {
				fullPathFile := getFileNameFromCommandLine(line)
				if fullPathFile == "" {
					instance.Write([]byte(fmt.Sprintf("cd cmd error type help\n")))
					continue
				}
				err = x.c.CD(fullPathFile)
				if err != nil {
					instance.Write([]byte(fmt.Sprintf("cd error: %s\n", err)))
				}
			} else if strings.HasPrefix(line, "get ") {
				parts := strings.Split(line, " ")
				if len(parts) > 2 {
					parts = parts[:2]
				}
				if len(parts) != 2 {
					instance.Write([]byte("incorrect file name type help"))
					continue
				}
				fullPathFile := parts[1]
				err = x.c.GetFile(instance, fullPathFile)
				if err != nil {
					return err
				}
			} else {
				instance.Write([]byte("unknown command, type help\n"))
			}
		}
	}

	if x.c == nil {
		return nil
	}
	err = x.c.Store()
	return err
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func getFileNameFromCommandLine(line string) string {
	parts := strings.Split(line, " ")
	if len(parts) > 2 {
		parts = parts[:2]
	}
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

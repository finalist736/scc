package cfs

import (
	"github.com/chzyer/readline/runes"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	cmdList = []string{"exit ", "help ", "du ", "ls ", "add ", "get ", "rm ", "cd ", "vacuum "}
)

func (x *CryptoFileSystem) Do(line []rune, pos int) (newLine [][]rune, length int) {
	names := make([][]rune, 0)

	if runes.HasPrefix(line, stringToRunes("add ")) {
		var (
			scanPath string
			scanFile string
		)
		if len(line) > 0 && runes.Index('/', line) >= 0 {
			scanPath = string(line[4:])
			scanPath, scanFile = path.Split(scanPath)
		} else {
			scanPath, _ = os.Getwd()
			scanFile = string(line[4:])
		}
		files, _ := ioutil.ReadDir(scanPath)
		for _, f := range files {
			name := f.Name()
			fr := []rune(name)
			if runes.HasPrefix(fr, []rune(scanFile)) {
				fr = fr[len(scanFile):]
				if f.IsDir() {
					fr = append(fr, '/')
				} else {
					fr = append(fr, ' ')
				}
				names = append(names, fr)
			}
		}
		return names, pos - (4 + 0)
	} else if runes.HasPrefix(line, stringToRunes("rm ")) {
		parts := strings.Split(string(line), " ")
		if len(parts) < 2 {
			return names, pos
		}
		fileName := parts[1]
		currentPos := pos - 3
		for _, item := range x.header.Items {
			if item.Removed {
				continue
			}
			if strings.HasPrefix(item.Name, fileName) {
				names = append(names, stringToRunes(item.Name[currentPos:]))
			}
		}
		return names, currentPos
	} else if runes.HasPrefix(line, stringToRunes("cd ")) {
		parts := strings.Split(string(line), " ")
		if len(parts) < 2 {
			return names, pos
		}
		dir := parts[1]
		var (
			scanPath string
			scanFile string
		)
		if len(dir) > 0 && strings.Index(dir, "/") >= 0 {
			scanPath = dir
			scanPath, scanFile = path.Split(scanPath)
		} else {
			scanPath, _ = os.Getwd()
			scanFile = dir
		}
		files, _ := ioutil.ReadDir(scanPath)
		for _, file := range files {
			name := file.Name()
			if !file.IsDir() {
				continue
			}
			if strings.HasPrefix(name, scanFile) {
				name = (scanPath + name + "/")[pos-3:]
				names = append(names, stringToRunes(name))
			}
		}
		return names, pos - 3
	} else if runes.HasPrefix(line, stringToRunes("get ")) {
		parts := strings.Split(string(line), " ")
		if len(parts) < 2 {
			return names, pos
		}
		fileName := parts[1]
		currentPos := pos - 4
		for _, item := range x.header.Items {
			if item.Removed {
				continue
			}
			if strings.HasPrefix(item.Name, fileName) {
				names = append(names, stringToRunes(item.Name[currentPos:]))
			}
		}
		return names, currentPos
	} else {
		scanCmd := string(line[:pos])
		for _, s := range cmdList {
			if strings.HasPrefix(s, scanCmd) {
				names = append(names, stringToRunes(s[pos:]))
			}
		}
	}
	return names, pos
}

func stringToRunes(s string) []rune {
	var result []rune
	for i := range s {
		result = append(result, rune(s[i]))
	}
	return result
}

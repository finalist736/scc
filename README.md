# Simple Crypto Container
SCC is a pseudo file system which can store files into crypto container

## installing
```bash
go get github.com/finalist736/scc
cd ~/go/src/github.com/finalist736/scc
go install
```
## Using

```bash
$ exit
```
- finish work and exit

```bash
$ help 
```
- type help information

```bash
$ add ./path/to/file.any
```
- add file to scc

```bash
$ get file.any 
```
- get file from scc and put it into current working directory

```bash
$ rm file.any
```
- remove file.any from scc and run vacuum

```bash
$ du
```
- shows used space in bytes

```bash
$ ls
```
- shows file list

```bash
$ pwd
```
- prints working directory

```bash
$ cd /working/directory
```
- changes local working directory

```bash
$ vacuum
```
- move files to empty places if exists
- run automatically when rm called


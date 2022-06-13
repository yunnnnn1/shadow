package options

import (
	"flag"
	"fmt"
	"os"
)

type ExecPars struct {
	Method  string
	Rserver string
	User    string
	Cmd     string
}

type FilePars struct {
	Method     string
	User       string
	Rserver    string
	Cip        string
	Path       string
	FileName   string
	IsRecover  string
	Timeout    string
	limitspeed string
}

// Comminfo 代码版本控制
var Comminfo = `
Version  : 4.1.0  
Author   : Mingyu
`

func execprint() {
	fmt.Fprintf(os.Stderr, `Name    : Shadow Client%s
Usage:
[远程Linux执行] 
shell> ./Sclient EXEC -Rserver 192.168.2.178:3721 -user oracle -cmd "python -u /tmp/test.py"

Options:
`, Comminfo)
}

func fileprint() {
	fmt.Fprintf(os.Stderr, `Name    : Shadow Client%s
Usage:
[文件下载]
shell> ./Sclient FILE -User root -Rserver 192.168.2.178:3721 -Cip 192.168.2.178:8000 -FileName test -Path /temp/ -IsRecover N -timeout 60 -limitspeed 60

Options:
`, Comminfo)
}

func defaultprint() {
	fmt.Fprintf(os.Stderr, `Name    : Shadow Client%s
Options:  EXEC|FILE  
Usage of EXEC:
  -Rserver string
        Your Remote Server Address (default "192.168.2.178:3721")
  -cmd string
        Your Commands
  -user string
        Your Remote Commands Exec User

Usage of FILE:
  -User string
        Your Remote File Owner (default "root")
  -Rserver string
        Your Remote Server Address (default "192.168.2.178:3721")
  -Cip string
        Software Center Server Address (default "192.168.2.178:3721")
  -FileName string
        Download software name
  -IsRecover string
        Y/N, Choose Y,file will be recover (default "Y")
  -Path string
        Path for download file 
  -limitspeed string
        Download File Speed (default "60")
  -timeout string
        Download timeout (default "60")

`, Comminfo)
}

func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func Myflag() interface{} {
	e := ExecPars{}
	f := FilePars{}

	// 对于不同的子命令，可以定义不同的flag
	execmd := flag.NewFlagSet("EXEC", flag.ExitOnError)
	execRserver := execmd.String("Rserver", "192.168.2.178:3721", "Your Remote Server Address")
	execUser := execmd.String("user", "", "Your Remote Commands Exec User")
	execMyCmd := execmd.String("cmd", "", "Your Commands")

	fileCmd := flag.NewFlagSet("FILE", flag.ExitOnError)
	fileUser := fileCmd.String("User", "root", "Your Remote File Owner")
	fileRserver := fileCmd.String("Rserver", "192.168.2.178:3721", "Software Center Server Address")
	fileCip := fileCmd.String("Cip", "192.168.2.178:3721", "Software Center Server Address")
	filePath := fileCmd.String("Path", "", "Path for download file ")
	fileName := fileCmd.String("FileName", "", "Download software name")
	fileIsRecover := fileCmd.String("IsRecover", "Y", "Y/N, Choose Y,file will be recover")
	fileDownloadTimeout := fileCmd.String("timeout", "60", "Download timeout")
	fileDownloadSpeed := fileCmd.String("limitspeed", "60", "Download File Speed")

	if len(os.Args) < 2 {
		defaultprint()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "EXEC":
		if in("-h", os.Args[2:]) {
			execprint()
			execmd.Parse(os.Args[2:])
		} else if len(os.Args[2:]) == 6 {
			execmd.Parse(os.Args[2:])
			e.Method = "EXEC"
			e.Rserver = *execRserver
			e.User = *execUser
			e.Cmd = *execMyCmd
			return e
		} else {
			defaultprint()
			os.Exit(1)
		}

	case "FILE":
		if in("-h", os.Args[2:]) {
			fileprint()
			fileCmd.Parse(os.Args[2:])
		} else if len(os.Args[2:]) == 16 {
			fileCmd.Parse(os.Args[2:])
			f.Method = "FILE"
			f.User = *fileUser
			f.Rserver = *fileRserver
			f.Cip = *fileCip
			f.Path = *filePath
			f.FileName = *fileName
			f.IsRecover = *fileIsRecover
			f.Timeout = *fileDownloadTimeout
			f.limitspeed = *fileDownloadSpeed
			return f
		} else {
			defaultprint()
			os.Exit(1)
		}

	default:
		defaultprint()
		os.Exit(1)
	}
	return nil
}

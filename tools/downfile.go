package tools

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

const (
	downloadtimeout = 60
	downloadspeed   = 60
)

var (
	mytimeout int
	myspeed   int
	myfile    string
)

func DownloadFile(outchan chan string, execuser string, url string, fn string, isRecover string, mypath string, timeout int, limitspeed int) {

	myuser, _ := user.Lookup(execuser)
	UserID, _ := strconv.Atoi(myuser.Uid)
	GroupUserId, _ := strconv.Atoi(myuser.Gid)

	if isRecover == "Y" {
		myfile = fn
	} else {
		myfile = fmt.Sprintf("%s_1", fn)
	}
	if timeout > downloadtimeout {
		mytimeout = timeout
	} else {
		mytimeout = downloadtimeout
	}
	if limitspeed > downloadspeed {
		myspeed = limitspeed
	} else {
		myspeed = downloadspeed
	}

	durl := fmt.Sprintf("http://%s/%s", url, fn)
	filearry := strings.Split(myfile, "/")
	tmpfilename := filearry[len(filearry)-1]
	filename := fmt.Sprintf("%s/%s", mypath, tmpfilename)
	startmsg := fmt.Sprintf("Start Download File '%s' to '%s' ", myfile, mypath)
	outchan <- MydateType(1, startmsg)

	client := http.DefaultClient
	client.Timeout = time.Second * time.Duration(mytimeout)
	resp, err := client.Get(durl)
	if err != nil {
		errmsg := fmt.Sprintf("%s...", err)
		outchan <- MydateType(2, errmsg)
	}
	if resp.ContentLength <= 0 {
		lenmsg := fmt.Sprintf("File is Null, Can not downloading...")
		fmt.Println(MydateType(2, lenmsg))
		outchan <- MydateType(2, lenmsg)

	} else {
		raw := resp.Body
		defer raw.Close()
		reader := bufio.NewReaderSize(raw, 1024*myspeed)

		file, err := os.Create(filename)
		file.Chown(UserID, GroupUserId)
		if err != nil {
			fmt.Println(70)
			myerr := fmt.Sprintf("%s", err)
			outchan <- MydateType(2, myerr)
		}
		writer := bufio.NewWriter(file)

		var endRead bool
		buff := make([]byte, 1*1024)
		readed := 0
		written := 0

		go func() {
			for {
				nr, er := reader.Read(buff)
				if nr > 0 {
					readed += nr
					nw, ew := writer.Write(buff[0:nr])
					if nw > 0 {
						written += nw
					}
					if ew != nil {
						err = ew
						break
					}
					if nr != nw {
						err = io.ErrShortWrite
						break
					}
				}
				if er != nil {
					if er != io.EOF {
						err = er
					} else {
						endRead = true
					}
					break
				}
			}
			if err != nil {
				myerr := fmt.Sprintf("%s", err)
				outchan <- MydateType(2, myerr)
			}
		}()

		spaceTime := time.Second * 1
		ticker := time.NewTicker(spaceTime)
		lastWtn := 0
		stop := false

		for {
			select {
			case <-ticker.C:
				speed := written - lastWtn
				msg := fmt.Sprintf("Download Speed %s/%s", bytesToSize(speed), spaceTime.String())
				outchan <- MydateType(1, msg)
				if (written-lastWtn == 0) && endRead && (readed == written) {
					ticker.Stop()
					stop = true
					writer.Flush()
					break
				}
				lastWtn = written
			}
			if stop {
				break
			}
		}
	}
	outchan <- MydateType(1, "Download Finish!")
	close(outchan)

}

func bytesToSize(length int) string {
	var k = 1024 // or 1024
	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	if length == 0 {
		return "0 Bytes"
	}
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(k)))
	r := float64(length) / math.Pow(float64(k), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}

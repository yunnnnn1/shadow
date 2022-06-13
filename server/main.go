package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jingmingyu/shadow/options"
	shadow "github.com/jingmingyu/shadow/proto/sd"
	"github.com/jingmingyu/shadow/tools"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

type ShadowStreamService struct {
}

var (
	port *int
	h    *bool
)

func ExecMyCommand(outchan chan string, execuser string, cmdstr string) {
	var closeChanGP sync.WaitGroup
	var (
		cmd *exec.Cmd = nil
	)
	if os.Getuid() == 0 {
		if execuser == "root" {
			cmd = exec.Command("bash", "-c", cmdstr)
		} else {
			cmd = exec.Command("sudo", "su", "-", execuser, "-c", cmdstr)
		}
	} else {
		cmd = exec.Command("bash", "-c", cmdstr)
	}

	os.Setenv("NAME", execuser)

	stderrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}

	// 标准错误输出
	closeChanGP.Add(1)
	go func() {
		defer stderrPipe.Close()
		defer closeChanGP.Done()
		scanerr := bufio.NewScanner(stderrPipe)
		for scanerr.Scan() { // 命令在执行的过程中, 实时地获取其输出
			// check out
			outchan <- string(scanerr.Bytes())
		}
	}()
	//标准正确输出
	closeChanGP.Add(1)
	go func() {
		defer stdoutPipe.Close()
		defer closeChanGP.Done()
		scanout := bufio.NewScanner(stdoutPipe)
		for scanout.Scan() {
			outchan <- string(scanout.Bytes())
		}
	}()

	//if err := cmd.Run(); err != nil {
	//	log.Printf("[ERROR] :%v", err)
	//	fmt.Println(21231231231)
	//}
	cmd.Wait()
	if !cmd.ProcessState.Success() {
		outchan <- strconv.Itoa(12433)
	}
	closeChanGP.Wait()
	close(outchan)
}

func checkDsnOk(dsn string) error {

	_, err := net.DialTimeout("tcp", dsn, 3*time.Second)
	if err != nil {
		log.Printf("[ERROR] Software Center Err On  : %s  \n", dsn)
		return err
	}
	return nil
}

func (s *ShadowStreamService) ExecStreamCMD(req *shadow.ExecRequest, srv shadow.ShadowStreamService_ExecStreamCMDServer) error {

	ch1 := make(chan string)

	if req.Action == "EXEC" {
		log.Printf("[INFO] ExecUser: %v ,Command: %s\n", req.User, req.Cmd)
		go ExecMyCommand(ch1, req.User, req.Cmd)
		for val := range ch1 {
			err := srv.Send(&shadow.ExecStreamResponse{
				Stdout: val,
			})
			if err != nil {
				return err
			}
		}
	} else if req.Action == "FILE" {
		timeout, _ := strconv.Atoi(req.Timeout)
		limitspeed, _ := strconv.Atoi(req.Limitspeed)
		log.Printf("[INFO] DownloadFile: %s\n", req.Filename)
		err := checkDsnOk(req.Cip)
		if err != nil {
			msg := fmt.Sprintf("Software Center Err On  : %s", req.Cip)
			zzz := tools.MydateType(2, msg)
			srv.Send(&shadow.ExecStreamResponse{
				Stdout: zzz,
			})
		} else {
			go tools.DownloadFile(ch1, req.User, req.Cip, req.Filename, req.IsRecover, req.Remotepath, timeout, limitspeed)
			for val := range ch1 {
				err := srv.Send(&shadow.ExecStreamResponse{
					Stdout: val,
				})
				if err != nil {
					return err
				}
			}
		}

	} else {
		msg := fmt.Sprintf("can not exec!")
		err := srv.Send(&shadow.ExecStreamResponse{
			Stdout: msg,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	h = flag.Bool("h", false, "This Help")
	port = flag.Int("port", 3721, "Server port")
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `Name    : Shadow Server%s
Usage:
[Server Start] ./Sserver -port 3721
Options:
`, options.Comminfo)
	flag.PrintDefaults()
}

func main() {

	flag.Parse()
	if *h {
		flag.Usage()
	} else {
		myport := fmt.Sprintf(":%d", *port)

		grpcServer := grpc.NewServer()

		shadow.RegisterShadowStreamServiceServer(grpcServer, &ShadowStreamService{})

		lst, err := net.Listen("tcp", myport)

		if err != nil {
			log.Println(err)
		}

		log.Printf("[INFO] Shadow Server Start at Pid:%d Port%s\n", os.Getpid(), myport)

		err = grpcServer.Serve(lst)
		if err != nil {
			log.Println("[ERROR] err :%v", err)
		}

	}
}

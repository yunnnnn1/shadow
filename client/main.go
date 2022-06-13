package main

import (
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/jingmingyu/shadow/options"
	shadow "github.com/jingmingyu/shadow/proto/sd"
	"github.com/jingmingyu/shadow/tools"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"os"
	"strings"
)

func main() {
	enc := mahonia.NewEncoder("utf-8")

	myoptions := options.Myflag()

	newoptions, _ := json.Marshal(&myoptions)

	var paras map[string]interface{}
	_ = json.Unmarshal(newoptions, &paras)

	var remote_conn string
	var myuser string
	remote_conn = fmt.Sprintf("%s", paras["Rserver"])

	conn, err := grpc.Dial(remote_conn, grpc.WithInsecure())
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		fmt.Println(tools.MydateType(2, msg))
	}
	defer conn.Close()
	client := shadow.NewShadowStreamServiceClient(conn)

	if paras["Method"] == "EXEC" {
		myuser = fmt.Sprintf("%s", paras["User"])
		mycmd := fmt.Sprintf("%s", paras["Cmd"])
		req := &shadow.ExecRequest{
			Cmd:        mycmd,
			Action:     "EXEC",
			User:       myuser,
			Cip:        "",
			Filename:   "",
			IsRecover:  "",
			Remotepath: "",
			Timeout:    "",
			Limitspeed: "",
		}

		stream, err := client.ExecStreamCMD(context.Background(), req)
		if err != nil {
			msg := fmt.Sprintf("Client Response %s", err)
			fmt.Println(tools.MydateType(2, msg))
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				msg := fmt.Sprintf("Client Response %s", err)
				fmt.Println(tools.MydateType(2, msg))
				os.Exit(1235)
			}

			if enc.ConvertString(resp.GetStdout()) == "12433" {
				myout := strings.TrimRight(enc.ConvertString(resp.GetStdout()), "12433")
				fmt.Printf("%s\n", myout)
				os.Exit(1235)
			} else {
				fmt.Printf("%s\n", enc.ConvertString(resp.GetStdout()))
			}

		}
	} else if paras["Method"] == "FILE" {
		myuser = fmt.Sprintf("%s", paras["User"])
		mypath := fmt.Sprintf("%s", paras["Path"])
		mycenterIp := fmt.Sprintf("%s", paras["Cip"])
		myFIleName := fmt.Sprintf("%s", paras["FileName"])
		myIsRecover := fmt.Sprintf("%s", paras["IsRecover"])
		myTimeout := fmt.Sprintf("%s", paras["timeout"])
		myLimitspeed := fmt.Sprintf("%s", paras["limitspeed"])
		req := &shadow.ExecRequest{
			Action:     "FILE",
			User:       myuser,
			Cip:        mycenterIp,
			Filename:   myFIleName,
			IsRecover:  myIsRecover,
			Remotepath: mypath,
			Timeout:    myTimeout,
			Limitspeed: myLimitspeed,
		}

		stream, err := client.ExecStreamCMD(context.Background(), req)
		if err != nil {
			msg := fmt.Sprintf("Client Response %s", err)
			fmt.Println(tools.MydateType(2, msg))
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				msg := fmt.Sprintf("Client Response %s", err)
				fmt.Println(tools.MydateType(2, msg))
				os.Exit(1235)
			}
			fmt.Printf("%s\n", enc.ConvertString(resp.GetStdout()))

		}
	}

}

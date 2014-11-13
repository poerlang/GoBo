package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/funny/link"
	"os"
	"path/filepath"
)

type Action struct {
	kind string
	txt  string
	from string
	to   string
}

var all map[string]*link.Session = map[string]*link.Session{}

func main() {

	//守护进程，开始
	d := flag.Bool("d", false, "Whether or not to launch in the background(like a daemon)")
	flag.Parse()
	if *d {
		fmt.Println(os.Args[0] + " will run in background.")
		filePath, _ := filepath.Abs(os.Args[0]) //将命令行参数中执行文件路径转换成可用路径
		//cmd := exec.Command(filePath, os.Args[2:]...)
		//将其他命令传入生成出的进程
		//cmd.Stdin = os.Stdin //给新进程设置文件描述符，可以重定向到文件中
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//cmd.Start() //开始执行新进程，不等待新进程退出
		args := append([]string{filePath}, os.Args[2:]...)
		os.StartProcess(filePath, args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return
	}
	//守护进程，结束

	proto := link.PacketN(4, binary.BigEndian)
	server, _ := link.Listen("tcp", "0.0.0.0:9999", proto)
	fmt.Println("GoBo is online , and wait for Client's msg...[moketao]")
	server.AcceptLoop(func(s *link.Session) {
		fmt.Println("session start from " + s.Conn().RemoteAddr().String())
		s.ReadLoop(func(msg []byte) {
			var obj interface{}
			err := json.Unmarshal(msg, &obj)
			if err == nil {
				ob := obj.(map[string]interface{})
				a := Action{}

				//行为类型
				akind := ob["kind"]
				if akind != nil {
					a.kind = akind.(string)
				}

				//具体内容
				var hasTxt bool
				atxt := ob["txt"]
				if atxt != nil {
					a.txt = atxt.(string)
					hasTxt = true //记录是否有内容
				}

				//发消息的是谁
				afrom := ob["from"]
				if afrom != nil {
					a.from = afrom.(string)
					all[a.from] = s //记录用户
				}

				//针对哪个人
				ato := ob["to"]
				if ato != nil {
					a.to = ato.(string)
					if v, ok := all[a.to]; ok {
						if hasTxt {
							v.Send(link.Binary(msg)) //转发
						}
					}
				}

				fmt.Printf("new message: %s\n", msg)

			} else {
				fmt.Println("格式有误")
			}
		})
		fmt.Println("session closed")
	})
	fmt.Println("end")
}

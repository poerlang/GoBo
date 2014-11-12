package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/funny/link"
)

type Action struct {
	kind string
	txt  string
	from string
	to   string
}

var all map[string]*link.Session = map[string]*link.Session{}

func main() {
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

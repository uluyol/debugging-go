package main

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Server struct {
	mu     sync.Mutex
	notify chan int
}

func (s *Server) Notify(args *NotifyArgs, reply *NotifyReply) error {
	s.mu.Lock()
	s.notify <- 1
	return nil
}

func AsyncSendMesg() <-chan int {
	c := make(chan int)
	go func() {
		go call("127.0.0.1:9999", "Server.Notify", &NotifyArgs{}, &NotifyReply{})
		c <- 1
	}()
	return c
}

func main() {
	s := StartServer("127.0.0.1:9999")
	<-AsyncSendMesg()
	<-AsyncSendMesg()
	<-s.notify
	<-s.notify
}

func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	log.Print(err)
	return false
}

func StartServer(me string) *Server {
	s := &Server{
		notify: make(chan int),
	}

	rpcs := rpc.NewServer()
	rpcs.Register(s)

	l, e := net.Listen("tcp", me)
	if e != nil {
		log.Fatal("listen error: ", e)
	}

	// create a thread to accept RPC connections from clients.
	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			}
			if err != nil {
				log.Printf("Server(%v) accept: %v", me, err.Error())
			}
		}
	}()

	return s
}

type NotifyArgs struct{}
type NotifyReply struct{}

package main

import "net"

func NotifyOnMesg(notify chan<- int) {
	conn, _ := net.ListenPacket("udp", "127.0.0.1:9999")
	var buf [1]byte
	conn.ReadFrom(buf[:])
	notify <- 1
}

func AsyncSendMesg() <-chan int {
	c := make(chan int)
	go func() {
		conn, _ := net.Dial("udp", "127.0.0.1:9999")
		conn.Write([]byte{1})
		c <- 1
	}()
	return c
}

func main() {
	c := make(chan int)
	go NotifyOnMesg(c)
	<-AsyncSendMesg()
	<-AsyncSendMesg()
	<-c
	<-c
}

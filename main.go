package main

import (
	"fmt"
	"net"
)

const (
	PING = "PING"
	SET  = "SET"
)

type RClient struct {
	net.Conn
}

func NewClient(serverAddr string) *RClient {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil
	}
	return &RClient{conn}
}

func (c RClient) Ping() string {

	c.sendCmd(PING)
	return c.readResult()
}

func (c RClient) Set(key, val string) bool {
	cmd := fmt.Sprintf("%s %s %s", SET, key, val)
	c.sendCmd(cmd)
	res := c.readResult()

	fmt.Println(res)
	return true
}

const CRLF = "\r\n"

func (c RClient) sendCmd(cmd string) {
	c.Write([]byte(cmd + CRLF))
	return
}

func (c RClient) readResult() string {
	res := make([]byte, 100)
	_, err := c.Read(res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(res)
}

func main() {

	rc := NewClient("127.0.0.1:6379")

	rc.Set("hi", "world")
}

func scanToken() {

}

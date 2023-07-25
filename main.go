package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
)

const (
	PING = "PING"
	SET  = "SET"
	GET  = "GET"
)

// can optimize to
const (
	SIMPLE_STRING byte = '+'
	ERRORS        byte = '-'
	INTEGER       byte = ':'
	BULK_STRING   byte = '$'
	ARRAY         byte = '*'
)

type RClient struct {
	net.Conn
	rx chan string
}

func NewClient(serverAddr string) *RClient {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil
	}
	rClient := &RClient{conn, make(chan string)}
	go func() {
		rClient.readTx()
	}()
	return rClient
}

func (c RClient) Ping() (string, error) {

	c.sendCmd(PING)
	res, err := c.readResponse()

	return res.(string), err
}

func (c RClient) Set(key, val string) bool {
	cmd := fmt.Sprintf("%s %s %s", SET, key, val)
	c.sendCmd(cmd)
	res, _ := c.readResponse()

	if res.(string) == "OK" {
		return true
	}
	return false
}

func (c RClient) Get(key string) string {
	cmd := fmt.Sprintf("%s %s", GET, key)
	c.sendCmd(cmd)
	res, _ := c.readResponse()

	return res.(string)
}

const CRLF = "\r\n"

func (c RClient) sendCmd(cmd string) {
	c.Write([]byte(cmd + CRLF))
	return
}

func (c RClient) readTx() {
	sc := bufio.NewScanner(c)
	sc.Split(ScanCRLF)
	for sc.Scan() {
		c.rx <- sc.Text()
	}
}

func (c RClient) readResponse() (interface{}, error) {
	token := <-c.rx
	//type
	switch token[0] {
	case SIMPLE_STRING, ERRORS:
		return token[1:], nil
	case INTEGER:
		val, err := strconv.Atoi(token[1:])
		if err != nil {
			return 0, err
		}
		return val, nil
	case ARRAY:
		// TODO
	case BULK_STRING:
		str := <-c.rx
		return str, nil
	}
	return nil, nil
}

func main() {

	rc := NewClient("127.0.0.1:6379")
	fmt.Println(rc.Set("hi", "world"))
	fmt.Println(rc.Get("hi"))
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func ScanCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\r', '\n'}); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

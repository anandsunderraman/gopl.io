package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// os.Args[0] is the program name
	// os.Args[1:] contains the actual arguments
	args := os.Args[1:]
	clocks := parseArgs(args)

	//c := clocks[0]
	for _, c := range clocks {
		conn := c.Connect()
		defer conn.Close()
		//collect all connections and read from them
		//go mustCopy(os.Stdout, conn)
	}
	readClocks(clocks)

	for {
		time.Sleep(time.Minute)
	}

}

type Clock struct {
	Timezone string
	Address  string
	Conn     net.Conn
}

func (c *Clock) Connect() net.Conn {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		log.Fatal(err)
	}
	c.Conn = conn
	return conn
}

func readClocks(clocks []*Clock) {
	for _, clock := range clocks {
		go func(c *Clock) {
			s := bufio.NewScanner(c.Conn)
			for s.Scan() {
				fmt.Printf("%s %v \n", c.Timezone+":", s.Text())
			}
			if s.Err() != nil {
				fmt.Printf("lost %s %s\n", c.Timezone, s.Err())
			}
		}(clock)
	}
}

func parseArgs(args []string) []*Clock {
	var clocks []*Clock
	for _, arg := range args {
		//fmt.Println(arg)
		// Each argument is expected to be in the format "Timezone=ip:port"
		po := strings.Split(arg, "=")
		if len(po) != 2 {
			fmt.Printf("Invalid argument: %s, should be of the format Timezone=IP:PORT, example NewYork=localhost:8000\n", arg)
			continue
		}
		timezone := strings.ToLower(po[0])
		clocks = append(clocks, &Clock{timezone, po[1], nil})
	}
	return clocks
}

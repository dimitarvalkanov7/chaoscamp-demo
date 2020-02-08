package docker

import (
	"log"
	"net"
	"time"
)

const (
	host = "127.0.0.1"
)

func IsPortAvialable(port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		log.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		//fmt.Println("Opened", net.JoinHostPort(host, port))
		return false
	}
	return true
}

// func IsPortAvialable(host string, port string) {
// 	for _, port := range ports {
// 		timeout := time.Second
// 		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
// 		if err != nil {
// 			fmt.Println("Connecting error:", err)
// 		}
// 		if conn != nil {
// 			defer conn.Close()
// 			fmt.Println("Opened", net.JoinHostPort(host, port))
// 		}
// 	}
// }

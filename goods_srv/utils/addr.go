package utils

import (
	"go.uber.org/zap"
	"net"
)

func GetFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")

	if err != nil {
		zap.S().Panicf("err:", err)
		return 0
	}

	l, err := net.ListenTCP("tcp", addr)

	if err != nil {
		zap.S().Panicf("err:", err)
		return 0
	}

	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			zap.S().Panicf("err:", err)
			return
		}
	}(l)

	return l.Addr().(*net.TCPAddr).Port
}

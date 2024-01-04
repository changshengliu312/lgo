package lgo

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept tcp listener accept function
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// Server lgo http server 支持热重启, 长链接
type Server struct {
	server   *http.Server
	listener *net.TCPListener
	closing  bool
	stop     chan int
}

// ListenAndServe 启动http服务
func (srv *Server) ListenAndServe() error {
	go HandleSignals(srv)
	go srv.server.Serve(tcpKeepAliveListener{srv.listener})
	<-srv.stop
	return nil
}

// Fork 热重启fork子进程，老进程停止接收请求
func (srv *Server) Fork() (int, error) {
	if srv.closing {
		return 0, fmt.Errorf("tcp server already forked")
	}
	srv.closing = true
	file, err := srv.listener.File()
	if err != nil {
		return 0, fmt.Errorf("tcp server get conn file fail:%s", err)
	}
	return StartNewProcess(file.Fd(), 0)
}

// Shutdown 热重启，让老进程退出
func (srv *Server) Shutdown() error {
	srv.listener.SetDeadline(time.Now())
	time.Sleep(time.Second * 5)
	srv.stop <- 1

	return nil
}

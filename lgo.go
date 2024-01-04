package lgo

import (
	"fmt"
	"lgo/utils"
	"net"
	"net/http"
)

func Run(port uint32) {
	//DefaultServeMux.InitTimeout(time.Second * 10)

	listenIP := utils.ItoIP(utils.GetLocalIP())
	addr := fmt.Sprintf("%s:%d", listenIP, port)
	listenAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic("invalid listen addr")
	}
	ln, err := GetTCPListener(listenAddr)
	if err != nil {
		panic(err)
	}

	srv := &Server{
		server: &http.Server{
			Handler: DefaultServeMux,
			Addr:    addr,
		},
		listener: ln,
		closing:  false,
		stop:     make(chan int),
	}

	srv.ListenAndServe()
}

package lgo

import (
	"fmt"
	"io"
	"lgo/log"
	"lgo/utils"
	"net"
	"net/http"
	"time"
)

// DefaultLogWriter log writer
var DefaultLogWriter io.Writer

func Run(port uint32) {
	//DefaultServeMux.InitTimeout(time.Second * 10)
	logWriter, err := log.NewLogWriter("lgo", 1073741824, 10)
	if err != nil {
		panic(err)
	}
	//DefaultLogWriter = conf.Log.Writer
	log.SetOutput(logWriter)
	DefaultLogWriter = logWriter
	DefaultLogWriter.Write([]byte(fmt.Sprintf("==>-----------------lgo start at %s-----------------\n==>\n", time.Now())))

	listenIP := utils.GetLocalIP()
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

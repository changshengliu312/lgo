package lgo

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// 是否开启debug模式
var defaultEnableDebugmode = false

// graceful const
const (
	GracefulEnvironKey    = "IsGracefulCat"
	GracefulEnvironStr    = GracefulEnvironKey + "=1"
	GracefulTCPListenerFd = 3 // 0 stdin 1 stdout 2 stderr
	GracefulUDPListenerFd = 4
)

// SysLog 调试模式打印系统日志
func SysLog(v ...interface{}) {
	if defaultEnableDebugmode {
		log.Output(2, "[SYS] "+fmt.Sprintln(v...))
	}
}

// SysLogf 调试模式打印格式化系统日志
func SysLogf(format string, v ...interface{}) {
	if defaultEnableDebugmode {
		data := fmt.Sprintf(format, v...)
		log.Output(2, "[SYS] "+data)
	}
}

// GracefulRestart 支持热重启服务需实现的接口
type GracefulRestart interface {
	Fork() (int, error)
	Shutdown() error
}

// HandleSignals 监听信号
func HandleSignals(srv GracefulRestart) {
	signalChan := make(chan os.Signal)

	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGSEGV)
	SysLogf("server notify signal: SIGTERM SIGUSR2")

	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGTERM:
			//SIGTERM软件终止信号
			SysLogf("receive SIGTERM signal, shutdown server")
			srv.Shutdown()
		case syscall.SIGSEGV:
			//SIGSEGV段非法错误信号
			SysLogf("receive SIGSEGV signal, shutdown server")
		case syscall.SIGUSR2:
			//SIGUSR2用户信号
			SysLogf("receive SIGUSR2 signal, graceful restarting server")
			if pid, err := srv.Fork(); err != nil {
				//热重启fork失败
				SysLogf("start new process failed: %v, continue serving", err)
			} else {
				SysLogf("start new process succeed, the new pid is %d", pid)
				srv.Shutdown()
			}
		default:
		}
	}
}

// StartNewProcess fork子进程，传入listener复用fd
func StartNewProcess(tcpfd, udpfd uintptr) (int, error) {
	SysLogf("graceful start new process, tcp fd:%v, udp fd:%v", tcpfd, udpfd)
	envs := []string{}
	for _, value := range os.Environ() {
		if value != GracefulEnvironStr {
			envs = append(envs, value)
		}
	}
	envs = append(envs, GracefulEnvironStr)

	execSpec := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), tcpfd, udpfd},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return fork, nil
}

// GetTCPListener 获取tcp listener
func GetTCPListener(addr *net.TCPAddr) (*net.TCPListener, error) {
	var ln *net.TCPListener
	var err error

	if os.Getenv(GracefulEnvironKey) != "" {
		//正常热重启
		SysLogf("tcp server get listener from os file")
		file := os.NewFile(GracefulTCPListenerFd, "")
		var listener net.Listener
		listener, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return nil, err
		}
		var ok bool
		if ln, ok = listener.(*net.TCPListener); !ok {
			return nil, fmt.Errorf("net.FileListener is not TCPListener")
		}
	} else {
		//异常启动量
		SysLogf("tcp server create listener from tcpaddr")
		ln, err = net.ListenTCP("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.ListenTCP error: %v", err)
			return nil, err
		}
	}
	return ln, nil
}

// GetUDPListener 获取udp listener
func GetUDPListener(addr *net.UDPAddr) (*net.UDPConn, error) {
	var ln *net.UDPConn
	var err error

	if os.Getenv(GracefulEnvironKey) != "" {
		//正常热重启
		SysLogf("udp server get listener from os file")
		file := os.NewFile(GracefulUDPListenerFd, "")
		var listener net.Conn
		listener, err = net.FileConn(file)
		if err != nil {
			err = fmt.Errorf("net.FileConn error: %v", err)
			return nil, err
		}
		var ok bool
		if ln, ok = listener.(*net.UDPConn); !ok {
			return nil, fmt.Errorf("net.FileConn is not UDPConn")
		}
	} else {
		//异常启动量
		SysLogf("udp server create listener from udpaddr")
		ln, err = net.ListenUDP("udp", addr)
		if err != nil {
			err = fmt.Errorf("net.ListenUDP error: %v", err)
			return nil, err
		}
	}
	return ln, nil
}

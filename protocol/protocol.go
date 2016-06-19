package protocol

import (
	"net"
	"net/rpc"

	"github.com/op/go-logging"
	"github.com/subgraph/go-procsnitch"
)

var log = logging.MustGetLogger("procsnitchd_protocol")

// SetLogger allows setting a custom go-logging instance
func SetLogger(logger *logging.Logger) {
	log = logger
}

type TCPDescriptor struct {
	SrcPort, DstPort uint16
	DstAddr          net.IP
}

type ProcsnitchRPC struct {
	procInfo procsnitch.ProcInfo
}

func NewProcsnitchRPC(procInfo procsnitch.ProcInfo) *ProcsnitchRPC {
	rpc := ProcsnitchRPC{
		procInfo: procInfo,
	}
	return &rpc
}

func (t *ProcsnitchRPC) LookupUNIXSocketProcess(socketFile *string, info *procsnitch.Info) error {
	newInfo := t.procInfo.LookupUNIXSocketProcess(*socketFile)
	*info = *newInfo
	return nil
}

func (t *ProcsnitchRPC) LookupTCPSocketProcess(tcpDescriptor *TCPDescriptor, info *procsnitch.Info) error {
	newInfo := t.procInfo.LookupTCPSocketProcess(tcpDescriptor.SrcPort, tcpDescriptor.DstAddr, tcpDescriptor.DstPort)
	*info = *newInfo
	return nil
}

func (t *ProcsnitchRPC) LookupUDPSocketProcess(srcPort *uint16, info *procsnitch.Info) error {
	newInfo := t.procInfo.LookupUDPSocketProcess(*srcPort)
	*info = *newInfo
	return nil
}

func ConnectionHandlerFactory(procInfo procsnitch.ProcInfo) func(conn net.Conn) error {
	return func(conn net.Conn) error {
		s := NewProcSnitchSession(conn, procInfo)
		return s.Start()
	}
}

type ProcSnitchSession struct {
	conn      net.Conn
	rpcServer *rpc.Server
}

func NewProcSnitchSession(conn net.Conn, procInfo procsnitch.ProcInfo) *ProcSnitchSession {
	p := ProcSnitchSession{
		conn:      conn,
		rpcServer: rpc.NewServer(),
	}
	rpc := NewProcsnitchRPC(procInfo)
	p.rpcServer.Register(rpc)
	return &p
}

func (s *ProcSnitchSession) Start() error {
	s.rpcServer.ServeConn(s.conn)
	return nil
}

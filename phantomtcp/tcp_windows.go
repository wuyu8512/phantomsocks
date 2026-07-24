package phantomtcp

import (
	"net"
	"sync"
	"syscall"
	"time"
)

var TFOCookies sync.Map
var TFOPayload [64][]byte
var TFOSynID uint8 = 0

func DialWithOption(laddr, raddr *net.TCPAddr, ttl, mss int, tcpfastopen, keepalive bool, timeout time.Duration) (net.Conn, error) {
	if tcpfastopen || keepalive {
		d := net.Dialer{Timeout: timeout, LocalAddr: laddr,
			Control: func(network, address string, c syscall.RawConn) error {
				err := c.Control(func(fd uintptr) {
					f := syscall.Handle(fd)
					if tcpfastopen {
						if raddr.IP.To4() == nil {
							syscall.SetsockoptInt(f, syscall.IPPROTO_IPV6, syscall.IPV6_UNICAST_HOPS, int(TFOSynID)|64)
						} else {
							syscall.SetsockoptInt(f, syscall.IPPROTO_IP, syscall.IP_TTL, int(TFOSynID)|64)
						}
						TFOSynID++
					}
					if keepalive {
						syscall.SetsockoptInt(f, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1)
					}
				})
				return err
			}}
		return d.Dial("tcp", raddr.String())
	} else {
		d := net.Dialer{Timeout: timeout, LocalAddr: laddr}
		return d.Dial("tcp", raddr.String())
	}
}

func GetOriginalDST(conn *net.TCPConn) (*net.TCPAddr, error) {
	LocalAddr := conn.LocalAddr()
	LocalTCPAddr := LocalAddr.(*net.TCPAddr)

	if ip4 := LocalTCPAddr.IP.To4(); ip4 != nil {
		if ip4[0] == 127 && ip4[1] == 255 {
			ip4[0] = VirtualAddrPrefix
			ip4[1] = 0
			LocalTCPAddr.IP = ip4
			RemoteTCPAddr := conn.RemoteAddr().(*net.TCPAddr).IP.To4()
			LocalTCPAddr.Port = int(RemoteTCPAddr[2])<<8 | int(RemoteTCPAddr[3])
		}
	}

	return LocalTCPAddr, nil
}

func SendWithOption(conn net.Conn, payload, oob []byte, tos, ttl int) error {
	return nil
}

func (outbound *Outbound)SendWithFakePayload(conn net.Conn, fakepayload, realpayload []byte) error {
	return nil
}

func GetTCPState(conn net.Conn) (uint8, error) {
	return 0, nil
}

func TProxyTCP(address string) {
}

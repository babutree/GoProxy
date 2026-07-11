package proxy

import (
	"io"
	"net"
	"testing"
	"time"
)

// TestSocks5BypassDialsLocalTargetDirectly 端到端验证 bypass 直连：
// 上游 store 为空（无任何可用节点），若目标是内网/本地地址，SOCKS5 仍应直连成功。
// 因为没有任何上游节点，连接能建立且回显成功，只可能来自 bypass 直连路径——
// 这直接证明内网目标绕过了上游选路。
func TestSocks5BypassDialsLocalTargetDirectly(t *testing.T) {
	// 本地 echo 服务器，模拟内网服务（127.0.0.1）。
	echoAddr, echoPort := startLocalEcho(t)

	store := newProxyTestStore() // 故意为空：无任何上游节点
	server := newSocks5TestServer(store, proxyTestConfig(0))

	clientConn, serverConn := net.Pipe()
	t.Cleanup(func() { _ = clientConn.Close() })
	go server.handleConnection(serverConn)

	writeSocks5NoAuthHandshake(t, clientConn)
	writeSocks5DomainRequest(t, clientConn, "127.0.0.1", echoPort)

	// 读 SOCKS5 应答：VER REP RSV ATYP + BND.ADDR(IPv4 4) + BND.PORT(2) = 10 字节。
	reply := make([]byte, 10)
	if err := readFullDeadline(t, clientConn, reply); err != nil {
		t.Fatalf("read socks5 reply: %v", err)
	}
	if reply[0] != 0x05 || reply[1] != 0x00 {
		t.Fatalf("socks5 reply = %#v, want VER=0x05 REP=0x00 (bypass 直连成功); echoAddr=%s", reply[:2], echoAddr)
	}

	// 通过隧道发数据，本地 echo 应原样返回，证明直连真实可用。
	payload := []byte("ping-bypass")
	if err := writeDeadline(t, clientConn, payload); err != nil {
		t.Fatalf("write payload: %v", err)
	}
	got := make([]byte, len(payload))
	if err := readFullDeadline(t, clientConn, got); err != nil {
		t.Fatalf("read echo: %v", err)
	}
	if string(got) != string(payload) {
		t.Fatalf("echo = %q, want %q", got, payload)
	}
}

// startLocalEcho 起一个 127.0.0.1 上的 TCP echo 服务器，返回其地址与端口。
func startLocalEcho(t *testing.T) (string, uint16) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen echo: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				io.Copy(conn, conn)
			}(c)
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := net.LookupPort("tcp", portStr)
	return ln.Addr().String(), uint16(port)
}

func readFullDeadline(t *testing.T, conn net.Conn, buf []byte) error {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetReadDeadline(time.Time{})
	_, err := io.ReadFull(conn, buf)
	return err
}

func writeDeadline(t *testing.T, conn net.Conn, buf []byte) error {
	t.Helper()
	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetWriteDeadline(time.Time{})
	_, err := conn.Write(buf)
	return err
}

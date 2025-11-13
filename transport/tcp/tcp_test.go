package tcp

import (
	"sync"
	"testing"
	"time"
	"wwww/transport"
)

// go test -v ./transport/tcp
// go test -v ./transport/tcp -timeout 5s

func TestTCP(t *testing.T) {
	server := NewTCPServer()
	err := server.Listen("127.0.0.1", 18080)
	defer server.Close()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n := 0
	var srvConn transport.TransportConn

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srvConn, err = server.Accept()
		if err != nil {
			t.Fatal(err)
		}
		n, err = srvConn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("received data: %s", string(buf[:n]))
	}()

	client := NewTCPClient()
	cltConn, err := client.Dial("127.0.0.1:18080")
	if err != nil {
		t.Fatal(err)
	}
	// sleep 1 second
	time.Sleep(1 * time.Second)
	cltConn.Write([]byte("hello"))

	// check the data
	if string(buf[:n]) != "hello" {
		t.Fatalf("data mismatch: %s != %s", string(buf[:n]), "hello")
	}
}

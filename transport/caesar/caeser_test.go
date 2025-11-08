package caesar

import (
	"sync"
	"testing"
	"time"
	"wwww/transport/tcp"
)

// go test -v ./transport/caesar
// go test -v ./transport/caesar -timeout 5s

func TestCaesar(t *testing.T) {
	debugHook := func(bytein, byteout []byte, msg string) {
		// print the bytein and byteout in hex
		t.Logf("debugHook: %s, bytein: %v, byteout: %v", msg, string(bytein), string(byteout))
	}
	server := NewCaesarServer(3, tcp.NewTCPServer(), &debugHook)
	err := server.Listen("127.0.0.1", 18080)
	defer server.Close()
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	buf := make([]byte, 1024)
	n := 0

	wg.Add(1)
	go func() {
		defer wg.Done()
		srvConn := <-server.Accept()
		n, err = srvConn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("received data: %s", string(buf[:n]))
	}()

	client := NewCaesarClient(3, tcp.NewTCPClient(), &debugHook)
	cltConn, err := client.Dial("127.0.0.1:18080")
	if err != nil {
		t.Fatal(err)
	}

	// sleep 1 second
	time.Sleep(1 * time.Second)
	cltConn.Write([]byte("hello"))

	// wait for the server to receive the data
	wg.Wait()
	// check the data
	if string(buf[:n]) != "hello" {
		t.Fatalf("data mismatch: %s != %s", string(buf[:n]), "hello")
	}
}

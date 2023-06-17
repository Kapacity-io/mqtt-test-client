package listeners

import (
	"bufio"
	"net"
	"strings"
	"sync"
)

type MessageHandlerFunc func(topic string, message string)

type TCPListener struct {
	addr     string
	listener net.Listener
	handler  MessageHandlerFunc
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewTCPListener(addr string, handler MessageHandlerFunc) (*TCPListener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &TCPListener{
		addr:     addr,
		listener: listener,
		handler:  handler,
		stopChan: make(chan struct{}),
	}, nil
}

func (l *TCPListener) Start() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()

		for {
			conn, err := l.listener.Accept()
			if err != nil {
				select {
				case <-l.stopChan:
					// The listener was stopped, so just return
					return
				default:
					// There was an error accepting the connection
					continue
				}
			}

			l.wg.Add(1)
			go func(c net.Conn) {
				l.handleConnection(c)
				l.wg.Done()
			}(conn)
		}
	}()
}

func (l *TCPListener) Stop() {
	close(l.stopChan)
	l.listener.Close()
	l.wg.Wait()
}

func (l *TCPListener) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read a line from the TCP connection
		line, err := reader.ReadString('\n')
		if err != nil {
			// There was an error reading the line
			return
		}

		line = strings.TrimSpace(line)

		// Assume the line is in the format "topic|message"
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			// The line was not in the expected format
			continue
		}

		topic := parts[0]
		message := parts[1]

		// Call the message handler
		l.handler(topic, message)
	}
}

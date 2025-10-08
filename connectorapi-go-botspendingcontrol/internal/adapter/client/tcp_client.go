package client

import (
	"connectorapi-go/internal/adapter/utils"
	"fmt"             
	// "io"
	"net"  // For TCP connections
	"time" // For timeouts
	"bufio"
)

// TCPSocketClient defines the interface for a TCP socket client.
type TCPSocketClient interface {
	SendAndReceive(address string, combinedPayloadString string) (string, error)
}

// BasicTCPSocketClient implements TCPSocketClient with length-prefixing and TIS-620 encoding.
type BasicTCPSocketClient struct {
	DialTimeout      time.Duration // Timeout for establishing the connection
	ReadWriteTimeout time.Duration // Timeout for read/write operations
}

// NewBasicTCPSocketClient creates a new instance of BasicTCPSocketClient.
func NewBasicTCPSocketClient(dialTimeout, readWriteTimeout time.Duration) *BasicTCPSocketClient {
	return &BasicTCPSocketClient{
		DialTimeout:      dialTimeout,
		ReadWriteTimeout: readWriteTimeout,
	}
}

// func (c *BasicTCPSocketClient) SendAndReceive(address string, combinedPayloadString string) (string, error) {
// 	fmt.Println("Connecting to the server...")
// 	conn, err := net.DialTimeout("tcp", address, c.DialTimeout)
// 	if err != nil {
// 		return "", fmt.Errorf("ER040: " + err.Error())
// 	}
// 	defer conn.Close()

// 	fmt.Println("Encoding request to CP874...")
// 	encodedRequest, err := utils.Utf8ToCP874(combinedPayloadString)
// 	if err != nil {
// 		return "", fmt.Errorf("ER099: Failed to encode request to CP874: " + err.Error())
// 	}

// 	_, err = conn.Write(encodedRequest)
// 	if err != nil {
// 		return "", fmt.Errorf("ER060: " + err.Error())
// 	}

// 	fmt.Println("Request sent (UTF-8):", combinedPayloadString)

// 	fmt.Println("Reading response until \\r\\n...")
// 	reader := bufio.NewReader(conn)
// 	var fullResponse []byte

// 	for {
// 		conn.SetReadDeadline(time.Now().Add(c.DialTimeout))

// 		line, err := reader.ReadBytes('\n')
// 		if err != nil {
// 			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
// 				return "", fmt.Errorf("ER060: read timeout")
// 			}
// 			if err == io.EOF {
// 				break
// 			}
// 			return "", fmt.Errorf("ER060: " + err.Error())
// 		}

// 		fullResponse = append(fullResponse, line...)

// 		if len(fullResponse) >= 2 && fullResponse[len(fullResponse)-2] == '\r' && fullResponse[len(fullResponse)-1] == '\n' {
// 			break
// 		}
// 	}

// 	decoded, err := utils.DecodeCP874(fullResponse)
// 	if err != nil {
// 		return "", fmt.Errorf("ER099: Failed to decode response from CP874: " + err.Error())
// 	}

// 	fmt.Println("Final result (UTF-8):", decoded)
// 	return decoded, nil
// }

func (c *BasicTCPSocketClient) SendAndReceive(address string, combinedPayloadString string) (string, error) {
	fmt.Println("Connecting to the server...")
	conn, err := net.DialTimeout("tcp", address, c.DialTimeout)
	if err != nil {
		return "", fmt.Errorf("ER040: " + err.Error())
	}
	defer conn.Close()

	fmt.Println("Encoding request to CP874...")
	encodedRequest, err := utils.Utf8ToCP874(combinedPayloadString)
	if err != nil {
		return "", fmt.Errorf("ER099: Failed to encode request to CP874: " + err.Error())
	}

	_, err = conn.Write(encodedRequest)
	if err != nil {
		return "", fmt.Errorf("ER060: " + err.Error())
	}

	fmt.Println("Request sent (UTF-8):", combinedPayloadString)

	conn.SetReadDeadline(time.Now().Add(c.ReadWriteTimeout))

	reader := bufio.NewReader(conn)
	responseLine, err := reader.ReadBytes('\n')
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "", fmt.Errorf("ER050: read timeout")
		}
		return "", fmt.Errorf("ER060: failed to read: %v", err)
	}

	decoded, err := utils.DecodeCP874(responseLine)
	if err != nil {
		return "", fmt.Errorf("ER099: Failed to decode response from CP874: " + err.Error())
	}

	fmt.Println("Final result (UTF-8):", decoded)
	return decoded, nil
}
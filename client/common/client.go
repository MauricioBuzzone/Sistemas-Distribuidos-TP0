package common

import (
	"net"
	"time"
	"os"
	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		c.conn.Close()
		log.Fatalf(
	        "action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// Close the client socket
func (c *Client) CloseSocket() {	
	if c.conn != nil {
		c.conn.Close()
	}
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	// Create the connection the server in every loop iteration. Send an
	c.createClientSocket()
	
	bet := Bet{
		ID:            c.config.ID,
		FirstName:     os.Getenv("NOMBRE"),
		LastName:	   os.Getenv("APELLIDO"),
		Document:	   os.Getenv("DOCUMENTO"),
		Birthdate:	   os.Getenv("NACIMIENTO"),
		Number:        os.Getenv("NUMERO"),
	}
	data := serializeBet(bet)
	sendBet(c.conn, data)

	msg, err := readMessage(c.conn)
	
	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
		msg[1],
		msg[2],
	)

	c.conn.Close()
	log.Infof("action: release_socket | result: success")

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	log.Infof("action: finished | result: success | client_id: %v", c.config.ID)
}

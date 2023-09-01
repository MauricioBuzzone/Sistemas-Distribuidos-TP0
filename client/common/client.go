package common

import (
	"bufio"
	"net"
	"time"

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

	// Send messages
	// Create the connection the server in every loop iteration. Send an
	log.Infof("Llegue hasta aca ")
	
	c.createClientSocket()
	
	log.Infof("Llegue hasta aca ")
	bet := Bet{
		ID:            c.config.ID,
		FirstName:     "FirstName",
		LastName:	   "LastName",
		Document:	   "Document",
		Birthdate:	   "Birthdate",
		Number:        "Number",
	}
	
	sendBet(c.conn, bet)


	log.Infof("action: apuesta_enviada | result: success | dni: %v  | numero: %v",
	bet.Document,
	bet.Number,
)
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	
	c.conn.Close()
	log.Infof("action: release_socket | result: success")

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)

	log.Infof("action: finished | result: success | client_id: %v", c.config.ID)
}

package common

import (
	"net"
	"os"
	"io"
	"os/signal"
    "syscall"
    "sync"
	"encoding/csv"
    "fmt"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	MaxPackageSize int
	BatchSize      int
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
	signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGTERM)
    var wg sync.WaitGroup
    connFinishChan := make(chan bool)
    c.createClientSocket()
    wg.Add(1)

    go func() {
        
		c.sendBets()
		c.conn.Close()
		log.Infof("action: release_socket | result: success")
		log.Infof("action: finish_client | result: success | client_id: %v", c.config.ID)
	
		connFinishChan <- true
        wg.Done()
    }()
    select {
    case <-signalChan:
        c.conn.Close()
    case <-connFinishChan:
        c.conn.Close()
    }
    wg.Wait()
}


func (c *Client) sendBets() error {
    filename := fmt.Sprintf("agency-%s.csv", c.config.ID)
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.FieldsPerRecord = 5

	betsAmount := 0
	sizePackage := 0
	data := serializeField(c.config.ID)
	for {
        betData, err := reader.Read()
        if err != nil {
            if err == io.EOF {
				c.sendBatch(data)
				break
            }
			return err
        }
		firstName:= betData[0]
        lastName:=  betData[1]
        document:=  betData[2]
        birthdate:= betData[3]    
		number:=    betData[4]
		
		bet := serializeBet(firstName,lastName,document,birthdate,number)
		
		// Rise the max package size or the batch is complete
		if sizePackage + len(bet) > c.config.MaxPackageSize || betsAmount == c.config.BatchSize {
			
			c.sendBatch(data)
			data = serializeField(c.config.ID)
			betsAmount = 0
			sizePackage = 0
		}
		data = append(data,bet...)
		betsAmount +=1
		sizePackage += len(bet)
    }
	err = sendMessage(c.conn, data, END)
	if err != nil{
		log.Infof("action: send_final_message | result: fail ")
		return err
	}
	return nil
}

func (c *Client) sendBatch(data []byte) error {
	err := sendMessage(c.conn, data, BET)
	if err != nil{
		log.Infof("action: send_batch | result: fail ")
		return err
	}
	
	msg, err := readMessage(c.conn)
	if err != nil{
		log.Infof("action: recv_confirm | result: fail ")
		return err
	}
	log.Infof("action: apuestas_enviadas | result: success | amount: %v",
	msg[1],
	)

	return nil
}
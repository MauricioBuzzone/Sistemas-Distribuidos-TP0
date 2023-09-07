package common

import (
	"net"
	"os"
	"io"
	"os/signal"
    "syscall"
	"encoding/csv"
    "fmt"
	"time"

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
	on bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		on: true,
		conn: nil,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGTERM)

    go func() {
		<-signalChan
		log.Infof("action: release_signal_chan | result: success")
		close(signalChan)
		c.on = false
    }()

	// The client sends all the bets to the servers
	c.sendBets()

	// The client inquires about the lottery winners. 
	// In case of not receiving a response, it will 
	// retry the inquiry later (exponential backoff).
	c.checkWinners()

	log.Infof("action: finish_client | result: success | client_id: %v", c.config.ID)
}

func (c *Client)checkWinners() {
	waitingTime := 1
	data := serializeField(c.config.ID)
	for c.on{
		err := c.createClientSocket()
		if err != nil{
			log.Fatalf("action: connect | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: consulta_ganadores | result: in_progress")
		err = sendMessage(c.conn, data, CHECK_WIN_TYPE)
		if err != nil{
			c.conn.Close()
			log.Infof("action: check_winners | result: fail | %v",err)
			return
		}
	
		msg, err := readMessage(c.conn)
		if err != nil{
			c.conn.Close()
			log.Infof("action:  | result: fail | %v",err)
			return
		}

		c.conn.Close()
		log.Infof("action: release_socket | result: success")

		if msg[0] == string(CHECK_WIN_TYPE) {
			// The bets from the other agencies are not loaded yet, waiting to check again.
			log.Infof("action: consulta_ganadores | result: in_progress | waitingTime: %v",waitingTime)
			time.Sleep(time.Duration(waitingTime) * time.Second)
			waitingTime = waitingTime * 2
		
		} else if msg[0] == string(WIN_TYPE) {
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v",
			len(msg[1:]),
			)
			log.Infof("action: consulta_ganadores | result: success | ganadores: %v",msg[1:])
			break
		
		} else{
			log.Fatalf("action: check_winners | result: fail | msg_server_type: %v",msg[0])
		}
	}
}

func (c *Client) sendBets() error {

	err := c.createClientSocket()
	if err != nil{
		log.Fatalf(
	        "action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

    filename := fmt.Sprintf("agency-%s.csv", c.config.ID)
    file, err := os.Open(filename)
    if err != nil {
		c.conn.Close()
        return err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.FieldsPerRecord = 5

	betsAmount := 0
	sizePackage := 0
	data := serializeField(c.config.ID)
	for c.on{
        betData, err := reader.Read()
        if err != nil {
            if err == io.EOF {
				c.sendBatch(data)
				break
            }
			c.conn.Close()
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
			
			err = c.sendBatch(data)
			if err != nil{
				c.conn.Close()
				return err
			}
			data = serializeField(c.config.ID)
			betsAmount = 0
			sizePackage = 0
		}
		data = append(data,bet...)
		betsAmount +=1
		sizePackage += len(bet)
    }
	data = serializeField(c.config.ID)
	err = sendMessage(c.conn, data, END_TYPE)
	if err != nil{
		c.conn.Close()
		log.Infof("action: send_final_message | result: fail ")
		return err
	}

	c.conn.Close()
	log.Infof("action: release_socket | result: success")

	return nil
}

func (c *Client) sendBatch(data []byte) error {
	err := sendMessage(c.conn, data, BET_TYPE)
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
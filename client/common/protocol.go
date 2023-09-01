package common

import(
	"net"
	"encoding/binary"
)

func sendBet(conn net.Conn, bet Bet) error {
	data := serializeBet(bet)
	sizeData := len(data)
	bytesWritten := 0

	sizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeBytes, uint32(len(data)))
	// Send the length of message
	_, err := conn.Write([]byte(sizeBytes))
	if err != nil {
		return err
	}

	// Send the bet serialized
	for bytesWritten < sizeData {
		n, err := conn.Write(data[bytesWritten:])
        if err != nil {
            return err
        }
        bytesWritten += n
	}

	return nil
}
	
func readMessage(conn net.Conn) (string, error) {
	// 
	sizeBytes := make([]byte, 4)
	n, err := conn.Read(sizeBytes)
	if err != nil {
		return "", err
	}

	totalSize := int(binary.BigEndian.Uint32(sizeBytes))
	bytesRead := 0
	field := ""

	for bytesRead < totalSize {
		// Read lenght file
		fieldLengthBytes := make([]byte, 4)
		_, err = conn.Read(fieldLengthBytes)
		if err != nil {
			return "", err
		}
		fieldLength := int(binary.BigEndian.Uint32(fieldLengthBytes))
		bytesRead += 4


		// Read field
		fieldData := make([]byte, fieldLength)
		n, err = conn.Read(fieldData)
		if err != nil {
			return "", err
		}
		field = string(fieldData)
		bytesRead += n
	}

	return field, nil
}	

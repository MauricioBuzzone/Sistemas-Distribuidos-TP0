package common

import(
	"net"
	"encoding/binary"
)

const LENGTH = 4
const BET = 'B'
const END = 'F'

func sendMessage(conn net.Conn, msj []byte, typeMsg byte) error {
	// Send the size of total msj
	sizeMsj := make([]byte, 4)
	binary.BigEndian.PutUint32(sizeMsj, uint32(len(msj)))
	_, err := sendAll(conn, sizeMsj)
	if err != nil {
		return err
	}

	// Send the type of message: Bet
	typeOfMessage := []byte{byte(typeMsg)}
	_, err = sendAll(conn, typeOfMessage)
	if err != nil {
		return err
	}

	// Send the bet serialized
	_, err = sendAll(conn, msj)
	if err != nil {
		return err
	}
	return nil
}
	
func sendAll(conn net.Conn,data []byte) (int, error) {
	totalBytes := len(data)
	bytesWritten := 0
	for bytesWritten < totalBytes {
		n, err := conn.Write(data[bytesWritten:])
		if err != nil {
			return bytesWritten, err
		}
		bytesWritten += n
	}
	return bytesWritten, nil
}

func read(conn net.Conn, bytes_to_read int) ([]byte, error) {
	bytes_readed := 0
	buffer := make([]byte, bytes_to_read)

	for bytes_readed < bytes_to_read{
		n, err := conn.Read(buffer[bytes_readed:])
		if err != nil {
			return buffer, err
		}
		bytes_readed += n
	}
	return buffer, nil
}

func readMessage(conn net.Conn) ([]string, error) {
	read_bytes := 0
    var fields []string
	
	// Read size message
	len_data, err := read(conn, LENGTH)
	if err != nil {
		return fields, err
	}
	len_msg := int(binary.BigEndian.Uint32(len_data))
		
	// Read type message
	type_data, err := read(conn, 1)
	if err != nil {
		return fields, err
	}
	type_msg := string(type_data)
	fields = append(fields,type_msg)
	
	// Read message payload
	for read_bytes < len_msg{
		len_field_data, err := read(conn, LENGTH)
		if err != nil {
			return fields, err
		}
		
		len_field := int(binary.BigEndian.Uint32(len_field_data))
		read_bytes += LENGTH

		field_data, err := read(conn, len_field)
		if err != nil {
			return fields, err
		}
		field := string(field_data)
		fields = append(fields,field)
		read_bytes += len(field)
	}
	return fields,nil
}	

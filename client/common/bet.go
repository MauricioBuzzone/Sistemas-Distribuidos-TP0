package common

import (
	"encoding/binary"
)

type Bet struct {
    ID            string
    FirstName     string
    LastName      string
    Document      string
    Birthdate     string
    Number        string
}

func serializeBet(firstName string,lastName string,document string,birthdate string,number string) []byte {

	data := []byte{}
    data = append(data, serializeField(firstName)...)
	data = append(data, serializeField(lastName)...)
	data = append(data, serializeField(document)...)
	data = append(data, serializeField(birthdate)...)
	data = append(data, serializeField(number)...)
	return data
}

func serializeField(field string) []byte {
	// For each field: First send his lenght and then the field it self
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(field)))
	return append(lengthBytes, []byte(field)...)
}
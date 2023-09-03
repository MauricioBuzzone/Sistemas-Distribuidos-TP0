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

func serializeBet(bet Bet) []byte {
    serialized := []byte{}
	serialized = append(serialized, serializeField(bet.ID)...)
    serialized = append(serialized, serializeField(bet.FirstName)...)
	serialized = append(serialized, serializeField(bet.LastName)...)
	serialized = append(serialized, serializeField(bet.Document)...)
	serialized = append(serialized, serializeField(bet.Birthdate)...)
	serialized = append(serialized, serializeField(bet.Number)...)
	return serialized
}

func serializeField(field string) []byte {
	// For each field: First send his lenght and then the field it self
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, uint32(len(field)))
	return append(lengthBytes, []byte(field)...)
}
package peer

import (
	"bytes"
	"encoding/gob"
)

type Payload struct {
	Name string
	Group string
}

func encodePayload(p *Payload) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodePayload(data []byte) (*Payload, error) {
	var payload *Payload
	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(payload); err != nil {
		return nil, err
	}

	return payload, nil
}

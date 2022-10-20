package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *JSONMap) MarshalBSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bson.MarshalExtJSON(b, false, false)
}

func (m *JSONMap) UnmarshalBSON(b []byte) error {
	var d bson.D
	err := bson.Unmarshal(b, &d)
	if err != nil {
		return err
	}
	temporaryBytes, err := bson.MarshalExtJSON(d, false, false)
	if err != nil {
		return err
	}
	return json.Unmarshal(temporaryBytes, m)
}

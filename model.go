package genericscrud

import (
	"bytes"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model[T any] struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`

	Data T `bson:"inline"`
}

type baseModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m *Model[T]) MarshalJSON() ([]byte, error) {
	base := baseModel{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	baseData, err := json.Marshal(base)
	if err != nil {
		return nil, err
	}

	baseData = bytes.Replace(baseData, []byte("}"), []byte(","), 1)

	data, err := json.Marshal(m.Data)
	if err != nil {
		return nil, err
	}

	data = bytes.Replace(data, []byte("{"), []byte(""), 1)

	return append(baseData, data...), nil
}

func (m *Model[T]) UnmarshalJSON(data []byte) error {
	var base baseModel

	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	var d T
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	m.ID = base.ID
	m.CreatedAt = base.CreatedAt
	m.UpdatedAt = base.UpdatedAt
	m.Data = d

	return nil
}

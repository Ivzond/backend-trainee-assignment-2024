package models

import "encoding/json"

type IntSlice []int64

func (is IntSlice) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	byteArray := src.([]byte)
	return json.Unmarshal(byteArray, &is)
}

type Banner struct {
	ID        int                    `json:"id"`
	TagIDs    IntSlice               `json:"tag_ids"`
	FeatureID int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

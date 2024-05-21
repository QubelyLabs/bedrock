package pipeline

import (
	"encoding/json"

	"github.com/rudderlabs/analytics-go/v4"
)

type Identify struct {
	analytics.Identify
	SourceId    string `json:"sourceId"`
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workspaceId"`
}

func (i *Identify) FromMap(data map[string]interface{}) error {
	// marshal data to bytes
	buf, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	// Create a new Identify instance
	// and unmarshall data to Identify
	err = json.Unmarshal(buf, i)
	if err != nil {
		return err
	}

	return nil
}

type Track struct {
	analytics.Track
	SourceId    string `json:"sourceId"`
	UserId      string `json:"userId"`
	WorkspaceId string `json:"workspaceId"`
}

func (t *Track) FromMap(data map[string]interface{}) error {
	// marshal data to bytes
	buf, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	// Create a new Identify instance
	// and unmarshall data to Identify
	err = json.Unmarshal(buf, t)
	if err != nil {
		return err
	}

	return nil
}

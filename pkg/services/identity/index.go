package identity

import (
	"fmt"
	"log"

	"github.com/qubelylabs/bedrock/pkg/request"
)

const (
	baseUrl = "https://localhost:5001"
)

func PreFlight(method string, url string, headers map[string]string) (bool, error) {
	payload := map[string]any{
		"method":  method,
		"url":     url,
		"headers": headers,
	}
	response, err := request.Post(fmt.Sprintf("%v/pre-flight", baseUrl), payload, nil, headers, 0)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if !response.Status {
		return false, nil
	}

	payload, ok := response.Data.(map[string]any)
	if !ok {
		return false, nil
	}

	status := payload["status"].(bool)

	return status, nil
}

func GetWorkspace(sourceId string) (string, string, error) {
	url := fmt.Sprintf("%v/workspace/%v/sourceId", baseUrl, sourceId)
	response, err := request.Get(url, nil, nil, 0)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	if !response.Status {
		return "", "", err
	}

	data := response.Data.(map[string]any)
	userId := data["userId"].(string)
	workspaceId := data["workspaceId"].(string)

	return userId, workspaceId, nil
}

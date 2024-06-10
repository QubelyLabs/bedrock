package pricing

import (
	"fmt"
	"log"

	"github.com/QubelyLabs/bedrock/pkg/request"
	"github.com/QubelyLabs/bedrock/pkg/util"
)

const (
	baseUrl = "http://localhost:8089"
)

// Enquiry makes an HTTP request to the pricing service to know the status of the users account
// It also caches the request for 10mins
func Enquiry(userId, workspaceId string) (bool, error) {
	url := fmt.Sprintf("%v/enquiry", baseUrl)
	payload := util.Object{
		"userId":      userId,
		"workspaceId": workspaceId,
	}

	response, err := request.Post(url, payload, nil, nil, 0)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if !response.Status {
		return false, nil
	}

	return true, nil
}

func Event(eventType string, payload any) {

	// if eventName == "" { // its an identify event
	// 	go actor.ProcessIdentityEvent(sourceId, userId, workspaceId, payload)
	// } else { // its a track event
	// 	go pricing.Event(rudderstack.TrackEvent, payload)
	// }
	// switch type {
	// case Iden:

	// }
	// err := rudderstack.Identity(analytics.Identify{

	// })
}

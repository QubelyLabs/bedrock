package pipeline

import (
	"log"
	"time"

	analytics "github.com/rudderlabs/analytics-go/v4"
)

type rudderStackService struct {
	client analytics.Client
}

func (s *rudderStackService) Identify(data Identify) error {
	return s.client.Enqueue(data)
}

func (s *rudderStackService) Track(data Track) error {
	return s.client.Enqueue(data)
}

func NewRudderStackService(dataPlaneUrl, key string) *rudderStackService {
	client, err := analytics.NewWithConfig(
		key,
		analytics.Config{
			DataPlaneUrl: dataPlaneUrl,
			Interval:     30 * time.Second,
			BatchSize:    100,
			Verbose:      true,
			DisableGzip:  false,
		},
	)
	if err != nil {
		log.Println(err)
	}

	return &rudderStackService{client}
}

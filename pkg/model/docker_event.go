package model

import (
	"github.com/docker/docker/api/types/events"
)

type DockerEvent struct {
	NodeId int64          `json:"node_id"`
	Event  events.Message `json:"event"`
}

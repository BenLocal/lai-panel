package model

import (
	"encoding/json"
	"time"
)

type Service struct {
	ID         int64     `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	AppID      int64     `db:"app_id" json:"app_id"`
	NodeID     int64     `db:"node_id" json:"node_id"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	Metadata   *string   `db:"metadata" json:"metadata"`
	DeployInfo *string   `db:"deploy_info" json:"deploy_info"`
}

type ServiceView struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	AppID    int64             `json:"app_id"`
	NodeID   int64             `json:"node_id"`
	Status   string            `json:"status,omitempty"`
	QAValues map[string]string `json:"qa_values"`
}

func (s *Service) ToView() *ServiceView {
	metadata := []*Metadata{}
	if s.Metadata != nil {
		json.Unmarshal([]byte(*s.Metadata), &metadata)
	}
	qa, ok := ToMetadataMap(metadata, "qa")
	if !ok {
		qa = make(map[string]string)
	}

	return &ServiceView{
		ID:       s.ID,
		Name:     s.Name,
		AppID:    s.AppID,
		NodeID:   s.NodeID,
		Status:   s.Status,
		QAValues: qa,
	}
}

func (v *ServiceView) ToModel() *Service {
	metadata := []*Metadata{
		{
			MetadataBase: MetadataBase{
				Name:       "qa",
				Properties: v.QAValues,
			},
		},
	}
	var metadataString *string
	metadataJson, _ := json.Marshal(metadata)
	if len(metadataJson) > 0 {
		s := string(metadataJson)
		metadataString = &s
	}
	return &Service{
		ID:       v.ID,
		Name:     v.Name,
		AppID:    v.AppID,
		NodeID:   v.NodeID,
		Status:   v.Status,
		Metadata: metadataString,
	}
}

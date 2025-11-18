package model

import "time"

type Service struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	AppID     int64     `db:"app_id" json:"app_id"`
	NodeID    int64     `db:"node_id" json:"node_id"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Metadata  *string   `db:"metadata" json:"metadata"`
}

type ServiceView struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	AppID    int64   `json:"app_id"`
	NodeID   int64   `json:"node_id"`
	Status   string  `json:"status"`
	Metadata *string `json:"metadata"`
}

func (s *Service) ToView() *ServiceView {
	return &ServiceView{
		ID:       s.ID,
		Name:     s.Name,
		AppID:    s.AppID,
		NodeID:   s.NodeID,
		Status:   s.Status,
		Metadata: s.Metadata,
	}
}

func (v *ServiceView) ToModel() *Service {
	return &Service{
		ID:       v.ID,
		Name:     v.Name,
		AppID:    v.AppID,
		NodeID:   v.NodeID,
		Status:   v.Status,
		Metadata: v.Metadata,
	}
}

package model

import (
	"encoding/json"
	"time"
)

type App struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Display       *string   `db:"display" json:"display"`
	Description   *string   `db:"description" json:"description"`
	DockerCompose *string   `db:"docker_compose" json:"docker_compose"`
	Version       string    `db:"version" json:"version"`
	Icon          string    `db:"icon" json:"icon"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
	QA            *string   `db:"qa" json:"qa"`
	Metadata      *string   `db:"metadata" json:"metadata"`
}

type AppQAItem struct {
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"default_value"`
	Options      []string `json:"options"`
	Required     bool     `json:"required"`
	Description  string   `json:"description"`
}

type AppView struct {
	ID            int64        `json:"id"`
	Name          string       `json:"name"`
	Display       *string      `json:"display"`
	Description   *string      `json:"description"`
	DockerCompose *string      `json:"docker_compose"`
	Version       string       `json:"version"`
	Icon          string       `json:"icon"`
	QA            []*AppQAItem `json:"qa"`
	Metadata      []*Metadata  `json:"metadata"`
}

func (a *App) ToView() *AppView {
	qa := []*AppQAItem{}
	if a.QA != nil {
		json.Unmarshal([]byte(*a.QA), &qa)
	}

	metadata := []*Metadata{}
	if a.Metadata != nil {
		json.Unmarshal([]byte(*a.Metadata), &metadata)
	}

	return &AppView{
		ID:            a.ID,
		Name:          a.Name,
		Display:       a.Display,
		Description:   a.Description,
		DockerCompose: a.DockerCompose,
		Version:       a.Version,
		Icon:          a.Icon,
		QA:            qa,
		Metadata:      metadata,
	}
}

func (a *AppView) ToModel() *App {
	var qaString *string
	qa, _ := json.Marshal(a.QA)
	if len(qa) > 0 {
		s := string(qa)
		qaString = &s
	}

	var metadataString *string
	metadata, _ := json.Marshal(a.Metadata)
	if len(metadata) > 0 {
		s := string(metadata)
		metadataString = &s
	}
	return &App{
		ID:            a.ID,
		Name:          a.Name,
		Description:   a.Description,
		DockerCompose: a.DockerCompose,
		Version:       a.Version,
		Icon:          a.Icon,
		QA:            qaString,
		Metadata:      metadataString,
		Display:       a.Display,
	}
}

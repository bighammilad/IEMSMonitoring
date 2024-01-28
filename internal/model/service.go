package model

import (
	"time"
)

type Service struct {
	ID            int                    `json:"id,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Address       string                 `json:"address,omitempty"`
	Method        string                 `json:"method,omitempty"`
	Header        map[string]string      `json:"header,omitempty"`
	Body          map[string]interface{} `json:"body,omitempty"`
	AccessLevel   AccessLevel            `json:"access_level,omitempty"`
	ExecutionTime int64                  `json:"execution_time,omitempty"`
	AllowedUsers  []int                  `json:"allowed_users,omitempty"`
}

type ErrorReport struct {
	ServiceName string    `json:"service_name,omitempty"`
	Log         string    `json:"log,omitempty"`
	OccurredAt  time.Time `json:"occurred_at,omitempty"`
}

type System struct {
	Services     []Service     `json:"services,omitempty"`
	Users        []User        `json:"users,omitempty"`
	ErrorReports []ErrorReport `json:"error_reports,omitempty"`
}

package model

import (
	"time"
)

type RequestErr struct {
	Code      ErrorCode      `json:"code,omitempty"`
	Message   string         `json:"message,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	TimeStamp time.Time      `json:"time_stamp,omitempty"`
	Path      string         `json:"path,omitempty"`
	RequestId string         `json:"request_id,omitempty"`
	Status    int            `json:"status,omitempty"`
	Fault     string         `json:"fault,omitempty"` // server | client
}

type Success struct {
	Data      any            `json:"data,omitempty"`
	Message   string         `json:"message,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
	RequestId string         `json:"request_id,omitempty"`
	TimeStamp time.Time      `json:"time_stamp,omitempty"`
}

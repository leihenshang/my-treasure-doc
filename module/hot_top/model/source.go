package model

import (
	"encoding/json"
	"time"
)


type HotData struct {
	Code        int       `json:"code"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Total       int       `json:"total"`
	Data        HotItems  `json:"data"`
	Params      any       `json:"params,omitempty"`
	UpdateTime  time.Time `json:"updateTime,omitempty"`
}

type HotItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Timestamp int64  `json:"timestamp"`
	Hot       int    `json:"hot"`
	URL       string `json:"url"`
	MobileURL string `json:"mobileUrl"`
	Author    string `json:"author,omitempty"`
	Desc      string `json:"desc,omitempty"`
}

type HotItems []*HotItem

func (h HotItems) ToJsonWithErr() string {
	jsonBytes, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}



func (h *HotData) IsUpdateTimeExpired(t time.Duration) bool {
	return time.Since(h.UpdateTime) > t
}

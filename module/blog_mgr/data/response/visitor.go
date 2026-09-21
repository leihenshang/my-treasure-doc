package response

import "time"

// VisitorStat 单个访客 IP 在时间范围内的聚合。
type VisitorStat struct {
	IP        string   `json:"ip"`
	Count     int64    `json:"count"`
	FirstSeen time.Time `json:"firstSeen"`
	LastSeen  time.Time `json:"lastSeen"`
	// Categories 该 IP 在时间范围内访问过的去重内容分类标识（如 post、diary）。
	Categories []string `json:"categories"`
}

// VisitorStats 是后台概览的访客统计：时间范围内的总量与按 IP 分组明细。
type VisitorStats struct {
	Days       int64         `json:"days"`
	TotalCount int64         `json:"totalCount"`
	DistinctIP int64         `json:"distinctIp"`
	List       []VisitorStat `json:"list"`
}

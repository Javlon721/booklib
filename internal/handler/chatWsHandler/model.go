package chatWsHandler

type UserStat struct {
	ID     string `json:"user_id"`
	Status string `json:"status"`
}

type ChatStats struct {
	ID        string     `json:"chat_id"`
	Connected int        `json:"connected"`
	Users     []UserStat `json:"users"`
}

type ShardStats struct {
	ChatStats *ChatStats `json:"chat_stats,omitempty"`
	ID        int        `json:"shard_id"`
	IsEmpty   bool       `json:"is_empty"`
}

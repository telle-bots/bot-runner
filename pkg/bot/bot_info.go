package bot

type Info struct {
	ID      int64  `bson:"_id"     json:"id"`
	Token   string `bson:"token"   json:"token"`
	Enabled bool   `bson:"enabled" json:"enabled"`
	TaskID  string `bson:"taskID"  json:"taskID,omitempty"`
}

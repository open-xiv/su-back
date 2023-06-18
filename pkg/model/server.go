package model

type ServerRecord struct {
	IP     string `json:"ip" bson:"ip"`
	Update int64  `json:"update" bson:"update"`
}

type ServerStatus struct {
	Status  string `json:"status" bson:"status"`
	Version string `json:"version" bson:"version"`
}

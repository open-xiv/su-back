package model

type Instance struct {
	ID    int    `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Level int    `json:"level" bson:"level"`
}

type Job struct {
	ID    int    `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Gauge string `json:"gauge" bson:"gauge"`
}

// Meta is struct for ffxiv (all instances and jobs)
type Meta struct {
	Instances []Instance `json:"instances" bson:"instances"`
	Jobs      []Job      `json:"jobs" bson:"jobs"`
}

type Oper struct {
	OpCode    string `json:"op_code" bson:"op_code"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

type Area struct {
	Op       Oper     `json:"op" bson:"op"`
	Instance Instance `json:"instance" bson:"instance"`
}

type Player struct {
	Op    Oper   `json:"op" bson:"op"`
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Job   Job    `json:"job" bson:"job"`
	Level int    `json:"level" bson:"level"`
}

type FightRecord struct {
	Area    Area     `json:"area" bson:"area"`
	Players []Player `json:"players" bson:"players"`
}

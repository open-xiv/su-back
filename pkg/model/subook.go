package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type PersonInfo struct {
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	AvatarURL string `json:"avatar_url" bson:"avatar_url"`
	Password  string `json:"password,omitempty" bson:"password"`
	Key       string `json:"key,omitempty" bson:"key"`
}

type SuMeta struct {
	Base  int `json:"base" bson:"base"`
	Total int `json:"total" bson:"total"`
}

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	Person       PersonInfo           `json:"person" bson:"person"`
	Meta         SuMeta               `json:"meta" bson:"meta"`
	FightIDs     []primitive.ObjectID `json:"fight_ids" bson:"fight_ids"`
	ServerRecord ServerRecord         `json:"server_record" bson:"server_record"`
}

type Fight struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	FightRecord  FightRecord        `json:"fight_record" bson:"fight_record"`
	ServerRecord ServerRecord       `json:"server_record" bson:"server_record"`
}

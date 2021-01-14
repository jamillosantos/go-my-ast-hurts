package models

type Test struct {
	ID   int64  `json:"id,uuidTest"`
	Name string `json:"name" bson:""`
}

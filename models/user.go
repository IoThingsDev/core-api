package models

type User struct {
	Id        string   `json:"id" bson:"_id,omitempty"`
	Firstname string   `json:"firstname" bson:"firstname"`
	Lastname  string   `json:"lastname" bson:"lastname"`
	Address   *Address `json:"address" bson:"address"`
}

const UsersCollection = "users"
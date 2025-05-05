package models

type Note struct {
	Id       string  `json:"id,omitempty" bson:"id,omitempty"`
	Name     *string `json:"name,omitempty" bson:"name,omitempty"`
	Content  *string `json:"content,omitempty" bson:"content,omitempty"`
	AuthorId uint    `json:"authorId,omitempty" bson:"authorId,omitempty"`
}

package entities

type PostStatus = string

const (
	PENDING  PostStatus = "PENDING"
	APPROVED PostStatus = "APPROVED"
	REJECTED PostStatus = "REJECTED"
)

type Post struct {
	Id          string     `json:"id,omitempty" bson:"_id,omitempty"`
	Status      PostStatus `json:"status,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

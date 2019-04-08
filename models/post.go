package models

import (
	"time"
)

type Post struct {
	Author string `json:"author"`
	Created time.Time `json:"created,omitempty"`
	Forum string `json:"forum,omitempty"`
	Id int `json:"id,omitempty"`
	IsEdited bool `json:"isEdited,omitempty"`
	Message string `json:"message"`
	Parent int `json:"parent,omitempty"`
	Thread int `json:"thread,omitempty"`
}

type Posts []Post

type PostsRelated struct {
	PostModel *Post `json:"post"`
	AuthorModel *User `json:"author,omitempty"`
	ThreadModel *Thread `json:"thread,omitempty"`
	ForumModel *Forum `json:"forum,omitempty"`
}

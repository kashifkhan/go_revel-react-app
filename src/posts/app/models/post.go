package models

import (
	"github.com/revel/revel"
)

type Post struct {
	Id        int64  `db:"Id" json:"Id"`
	Title     string `db:"Title" json:"Title"`
	Blog      string `db:"Blog" json:"Blog"`
	CreatedAt int64  `db:"CreatedAt" json:"CreatedAt"`
}

func (post *Post) Validate(v *revel.Validation) {

	v.Check(post.Title, revel.ValidRequired())
}

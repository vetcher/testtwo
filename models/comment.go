package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Text     string `json:"text"`
	AuthorID uint   `gorm:"index" json:"author"`
}

func PostComment(db *Database, c *Comment) (*Comment, error) {
	if err := db.db.Create(c); err.Error != nil {
		return nil, DBError(err.Error)
	} else {
		return c, nil
	}
}

// не понятно, что возвращать после Delete
func DeleteComment(db *Database, id uint) (bool, error) {
	var temp Comment
	if err := db.db.Where("id = ?", id).First(&temp); err.Error != nil {
		return false, DBError(err.Error)
	} else {
		if err := db.db.Delete(temp); err.Error != nil {
			return false, DBError(err.Error)
		}
		return true, nil
	}
}

func SelectCommentByID(db *Database, id uint) (*Comment, error) {
	var c Comment
	if err := db.db.Where("id = ?", id).First(&c); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return &c, nil
}

func LoadCommentsForUser(db *Database, u *User) ([]*Comment, error) {
	if err := db.db.Model(u).Related(&u.Comments, "AuthorID"); err.Error != nil {
		return nil, DBError(err.Error)
	}
	return u.Comments, nil
}

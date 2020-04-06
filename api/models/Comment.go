package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint32    `gorm:"not null" json:"userId"`
	PostID    uint64    `gorm:"not null" json:"postId"`
	Body      string    `gorm:"text;not null;" json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Comment) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate() error {
	if c.Body == "" {
		return errors.New("Required Body")
	}

	if c.UserID < 1 {
		return errors.New("Required UserID")
	}

	if c.PostID < 1 {
		return errors.New("Required PostID")
	}

	return nil
}

func (c *Comment) AddComment(db *gorm.DB) (*Comment, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}

	return c, nil
}

func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

// Delete comments that users have made if they delete their account
func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint32) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("user_id = ?", uid).Find(&comments).Delete(&comments)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}

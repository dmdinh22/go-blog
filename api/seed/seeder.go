package seed

import (
	"log"

	"github.com/dmdinh22/go-blog/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Username: "Tester McTesterson",
		Email:    "Tester.McTesterson@mailinator.com",
		Password: "p@$$w0rd",
	},
	models.User{
		Username: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "p@$$w0rd",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Test Title 1",
		Content: "Lorem Khaled Ipsum is a major key to success. I’m up to something. Wraith talk.",
	},
	models.Post{
		Title:   "Test Title 2",
		Content: "In life you have to take the trash out, if you have trash in your life, take it out, throw it away, get rid of it, major key.",
	},
}

var comments = []models.Comment{
	models.Comment{
		UserID: 1,
		PostID: 1,
		Body:   "You see the hedges, how I got it shaped up? It’s important to shape up your hedges, it’s like getting a haircut, stay fresh.",
	},
	models.Comment{
		UserID: 2,
		PostID: 2,
		Body:   "The key is to enjoy life, because they don’t want you to enjoy life. I promise you, they don’t want you to jetski, they don’t want you to smile. Bless up.",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}).Error

	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error

		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}

		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error

		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}

		err = db.Debug().Model(&models.Comment{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}

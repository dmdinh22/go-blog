package seed

import (
	"log"

	"github.com/dmdinh22/go-blog/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "Tester McTesterson",
		Email:    "Tester.McTesterson@mailinator.com",
		Password: "p@$$w0rd",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "p@$$w0rd",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Test Title 1",
		Content: "Lorem Khaled Ipsum is a major key to success. I’m up to something. Wraith talk. The key to more success is to get a massage once a week, very important, major key, cloth talk. Don’t ever play yourself. Mogul talk. I’m up to something. Cloth talk. The key is to enjoy life, because they don’t want you to enjoy life. I promise you, they don’t want you to jetski, they don’t want you to smile. Always remember in the jungle there’s a lot of they in there, after you overcome they, you will make it to paradise. A major key, never panic. Don’t panic, when it gets crazy and rough, don’t panic, stay calm."
	},
	models.Post{
		Title:   "Test Title 2",
		Content: "In life you have to take the trash out, if you have trash in your life, take it out, throw it away, get rid of it, major key. Surround yourself with angels. Lion! Hammock talk come soon. They key is to have every key, the key to open every door. To succeed you must believe. When you believe, you will succeed. Special cloth alert. I’m up to something. To be successful you’ve got to work hard, to make history, simple, you’ve got to make it."
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error

	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error

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
	}
}

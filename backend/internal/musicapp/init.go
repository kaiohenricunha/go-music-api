package musicapp

import (
	"gorm.io/gorm"
	"log"
)

func DbInit(db *gorm.DB) error {
	// Migrate the schema
	if err := db.AutoMigrate(&Song{}); err != nil {
		return err
	}

	// Seed initial data if necessary
	err := seedData(db)
	if err != nil {
		log.Printf("Failed to seed data: %v", err)
	}

	return nil
}

func seedData(db *gorm.DB) error {
	var count int64
	db.Model(&Song{}).Count(&count)
	if count == 0 {
		songs := []Song{
			{Name: "Burn It Down", Artist: "Linkin Park"},
			{Name: "Earth Song", Artist: "Michael Jackson"},
			{Name: "The Show Must Go On", Artist: "Queen"},
			{Name: "Bohemian Rhapsody", Artist: "Queen"},
			{Name: "Mucuripe", Artist: "Fagner e Belchior"},
			{Name: "Cleaning Out My Closet", Artist: "Eminem"},
			{Name: "Lose Yourself", Artist: "Eminem"},
			{Name: "Aquemini", Artist: "Outkast"},
			{Name: "Blood Moon", Artist: "Tim Henson"},
			{Name: "Computadores Fazem Arte", Artist: "Nação Zumbi"},
		}
		for _, song := range songs {
			err := db.Create(&song).Error
			if err != nil {
				return err
			}
		}
		log.Println("Seeded initial data successfully")
	} else {
		log.Println("Database already seeded")
	}
	return nil
}

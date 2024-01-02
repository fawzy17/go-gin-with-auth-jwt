package initializers

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB")
	database, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic("Failed to connect database")
	}

	// database.AutoMigrate(&models.Kelas{})
	// database.AutoMigrate(&models.Mahasiswa{})

	DB = database
}

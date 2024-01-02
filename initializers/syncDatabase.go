package initializers

import (
	"github.com/jwt-auth/first-try/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}

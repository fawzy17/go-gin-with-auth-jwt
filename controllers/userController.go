package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jwt-auth/first-try/initializers"
	"github.com/jwt-auth/first-try/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context) {
	// Get email & pass
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Failed to read body",
			Data:       nil,
		})

		return
	}

	// Hash the pass
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Failed to hash password",
			Data:       nil,
		})

		return
	}

	// Create user

	newUser := models.User{
		Email:    user.Email,
		Password: string(hash),
	}
	if err := initializers.DB.Create(&newUser).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Failed to create user",
			Data:       nil,
		})

		return
	}

	// Response
	ctx.JSON(http.StatusCreated, models.Response{
		StatusCode: http.StatusCreated,
		IsSuccess:  true,
		Message:    "Data berhasil dibuat",
		Data:       newUser,
	})
}

func LogIn(ctx *gin.Context) {
	// Get email and pass body
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Failed to read body",
			Data:       nil,
		})

		return
	}
	// Look up requested user
	var availableUser models.User
	initializers.DB.First(&availableUser, "email = ?", user.Email)

	if availableUser.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Invalid Email or Password",
			Data:       nil,
		})

		return
	}

	// Compare sent in pass with saved user pass hash
	if err := bcrypt.CompareHashAndPassword([]byte(availableUser.Password), []byte(user.Password)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Invalid Email or Password",
			Data:       nil,
		})

		return
	}

	// generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": availableUser.ID,
		// decide the expiration date of the token
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Invalid to create token",
			Data:       nil,
		})

		return
	}

	// sent it back
	ctx.JSON(http.StatusCreated, models.Response{
		StatusCode: http.StatusCreated,
		IsSuccess:  true,
		Message:    "Success",
		Data: map[string]interface{}{
			"email": availableUser.Email,
			"token": tokenString,
		},
	})
}

func Validate(ctx *gin.Context) {
	// Lanjutkan dengan logika controller Anda
	userId := GetUserId(ctx)
	if userId == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
			StatusCode: http.StatusBadRequest,
			IsSuccess:  false,
			Message:    "Invalid to find user",
			Data:       nil,
		})

		return
	}
	ctx.JSON(http.StatusOK, models.Response{
		StatusCode: http.StatusOK,
		IsSuccess:  true,
		Message:    "Validated",
		Data: map[string]interface{}{
			"user_id": userId,
		},
	})
}


func GetUserId(ctx *gin.Context) int {
	user, exists := ctx.Get("user")

	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
			StatusCode: http.StatusUnauthorized,
			IsSuccess:  false,
			Message:    http.StatusText(http.StatusUnauthorized),
			Data:       nil,
		})

		return 0
	}

	// Konversi nilai "user" ke tipe yang sesuai
	userObject, ok := user.(models.User)

	if !ok {
		// Jika tipe data tidak sesuai dengan yang diharapkan
		// Lakukan penanganan sesuai kebutuhan, misalnya memberikan respons Internal Server Error
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
			StatusCode: http.StatusInternalServerError,
			IsSuccess:  false,
			Message:    http.StatusText(http.StatusInternalServerError),
			Data:       nil,
		})

		return 0
	}

	// Sekarang Anda dapat menggunakan userObject sesuai kebutuhan dalam controller Anda
	// ...

	// Contoh: Menggunakan ID pengguna dari userObject
	userID := userObject.ID

	return int(userID);
}

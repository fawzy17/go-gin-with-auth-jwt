package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jwt-auth/first-try/initializers"
	"github.com/jwt-auth/first-try/models"
)

func RequireAuth(ctx *gin.Context) {
	// Get token from header 
    fullTokenString := ctx.GetHeader("Authorization")

    // Memeriksa apakah Authorization header tidak kosong dan dimulai dengan "Bearer "
    if fullTokenString == "" || !strings.HasPrefix(fullTokenString, "Bearer ") {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
            StatusCode: http.StatusUnauthorized,
            IsSuccess:  false,
            Message:    http.StatusText(http.StatusUnauthorized),
            Data:       nil,
        })

        return
    }

    // Mengambil nilai token tanpa "Bearer "
    tokenString := strings.TrimPrefix(fullTokenString, "Bearer ")

    // decode and validate
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
			StatusCode: http.StatusUnauthorized,
			IsSuccess:  false,
			Message:    err.Error(),
			Data:       nil,
		})

		return
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				StatusCode: http.StatusUnauthorized,
				IsSuccess:  false,
				Message:    http.StatusText(http.StatusUnauthorized),
				Data:       nil,
			})
	
			return
		}
		// Find the user with token sub
		var user models.User
		initializers.DB.First(&user,  claims["sub"])

		if user.ID == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				StatusCode: http.StatusUnauthorized,
				IsSuccess:  false,
				Message:    http.StatusText(http.StatusUnauthorized),
				Data:       nil,
			})
	
			return
		}
		// Attach to req
		ctx.Set("user", user)

		// Continue
		fmt.Println(claims["exp"], claims["sub"])
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
            StatusCode: http.StatusUnauthorized,
            IsSuccess:  false,
            Message:    http.StatusText(http.StatusUnauthorized),
            Data:       nil,
        })

        return
	}

}


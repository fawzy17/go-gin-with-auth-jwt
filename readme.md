1. init project
go mod init github.com/jwt-auth/first-try

2. import library
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get github.com/gin-gonic/gin 
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v4
go get github.com/joho/godotenv
go get github.com/githubnemo/CompileDaemon

3. install daemon first-try.exe
go install github.com/githubnemo/CompileDaemon

4. make main.go
5. compile daemon for auto refresh and run while changes applied
compiledaemon --command="./first-try"

6. make initializer to loadEnvVariable.go and connectToDb.go on .env add
DB = "root to connect the database"
PORT = (the database port)


7. make model userModel.go

8. make syncDatabase.go

9. make userController.go
SignUp func
LogIn func
10. on .env add JWT_SECRET_KEY with ur own key free to create randomly

11. on userController.go add
Validate func 
just to check route

12. make middleware requireAuth.go

13. complete the Validate func on userController.go
package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"movie-api/api/database"
	helper "movie-api/api/resource/user/helpers"
	models "movie-api/api/resource/user/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// Use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New()

var rootContext = context.Background()

// Handle password hashing
func HashPassword(password string) string {
	// GenerateFromPassword returns the bcrypt hash of the password at the given cost.
	// If the cost given is less than MinCost, the cost will be set to DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(hashedUserPassword, providedClearTextPassword string) (bool, string) {
	// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
	// Returns nil on success, or an error on failure.
	err := bcrypt.CompareHashAndPassword([]byte(providedClearTextPassword), []byte(hashedUserPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Password is incorrect!")
		check = false
	}

	return check, msg
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			log.Println("Error Here: ", err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Find user with email address in the user DB
		err := userCollection.FindOne(rootContext, bson.M{"email_address": user.Email_address}).Decode(&foundUser)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		// Password is invalid
		if !passwordIsValid {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}

		// Email address does not exist
		if foundUser.Email_address == nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Generate tokens
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email_address, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *&foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(rootContext, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// Return logged in user
		c.IndentedJSON(http.StatusOK, foundUser)
	}
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Bind JSON request body to User struct.
		// See https://github.com/iden3/go-iden3-servers/issues/6 for information
		if err := c.ShouldBindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Returns InvalidValidationError for bad validation input, nil or ValidationErrors ( []FieldError )
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
			return
		}

		// Check if there's a user with the same email address.
		// If count > 0, means email address already exist
		emailCount, err := userCollection.CountDocuments(rootContext, bson.M{"email_address": user.Email_address})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while checking for user email!"})
			return
		}

		if emailCount > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Email already exists!"})
		}

		// Hash password
		password := HashPassword(*user.Password)
		user.Password = &password

		// Check if there's a user with the same phone number.
		// If count > 0, means phone number already exist
		phoneNumberCount, err := userCollection.CountDocuments(rootContext, bson.M{"phone_number": user.Phone_number})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error occurred while checking for user phone number!"})
			return
		}

		if phoneNumberCount > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Phone number already exists!"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email_address, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(rootContext, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User was not created")
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		// Return user
		c.IndentedJSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetUsers() {}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get queried user by user_id
		userId := c.Param("user_id")

		// Handle verification of user type to user_id matching
		err := helper.MatchUserTypeToUid(c, userId)

		if err != nil {
			// Return error if it exists
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Get User model
		var user models.User

		// Find user by user_id in DB userCollection
		err = userCollection.FindOne(rootContext, bson.M{"user_id": userId}).Decode(&user)

		if err == mongo.ErrNoDocuments {
			// Log error
			fmt.Printf("No user was found with the user_id %s\n", userId)
			return
		}

		// Return error 500
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}

		// Return user with status 200
		c.IndentedJSON(http.StatusOK, user)
	}
}

package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"dsb/config"
	"dsb/models"
	"dsb/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and generates a JWT token.
func Login(username, password string) (*models.LoginResponse, error) {
	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var user models.User

	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No match")
			return nil, fmt.Errorf("user not found")
		}

		log.Printf("Error querying database: %v", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		fmt.Println("Invalid credentials")
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJwtToken(username, user.Regnumber)

	fmt.Printf("TOKEN:: %v\n", token) // Corrected print statement
	if err != nil {
		fmt.Printf("error while generating token: %v\n", err)
		return nil, err
	}

	res := &models.LoginResponse{
		Regnumber: user.Regnumber,
		Username:  user.Username,
		Token:     token,
	}

	return res, nil
}

// AddUser adds a new user to the database.
func AddUser(user *models.User) (*models.User, error) {
	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername checks if a user exists by username.
func FindUserByUsername(username string) (*models.User, error) {
	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		log.Printf("Error querying database: %v", err)
		return nil, err // Other errors
	}

	return &user, nil // User found
}

// FindUserByRegnumber checks if a user exists by registration number.
func FindUserByRegnumber(regnumber int) (*models.User, error) {
	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"regnumber": regnumber}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		log.Printf("Error querying database: %v", err)
		return nil, err // Other errors
	}

	return &user, nil // User found
}

// AddAdmin adds a new admin to the database.
func AddAdmin(admin *models.Admin) (*models.Admin, error) {
	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("admin")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, admin)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// AdminLogin validates admin credentials and returns the admin if successful.
func AdminLogin(role, password string) (*models.Admin, error) {
	if config.DatabaseClient == nil {
		log.Fatal("MongoDB client is nil")
	}

	collection := config.DatabaseClient.Database("DSBox").Collection("admin")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin models.Admin
	// Find the admin by role (you may want to adjust this based on your schema)
	err := collection.FindOne(ctx, bson.M{"role": role}).Decode(&admin)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("admin not found")
		}
		log.Printf("Error querying database: %v", err)
		return nil, err // Other errors
	}

	// Compare hashed password with provided password
	if err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &admin, nil // Admin authenticated successfully
}

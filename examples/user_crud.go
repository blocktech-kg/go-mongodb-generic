package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/blocktech-kg/go-mongodb-generic/mongodb"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string            `bson:"name"`
	Email     string            `bson:"email"`
	Age       int               `bson:"age"`
	IsActive  bool              `bson:"is_active"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

func main() {
	ctx := context.Background()
	
	// Connect to database
	db, err := mongodb.Connect(ctx, "mongodb://localhost:27017", "example_db")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	
	// Create service
	userCollection := db.Collection("users")
	userService := mongodb.NewGenericObjectDBCtrl[User](userCollection)
	
	// Example operations
	runCRUDExamples(ctx, userService)
}

func runCRUDExamples(ctx context.Context, userService mongodb.CRUDDBService[User]) {
	log.Println("Starting CRUD examples...")
	
	// 1. Create user
	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "Alice Johnson",
		Email:    "alice@example.com",
		Age:      28,
		IsActive: true,
	}
	
	if err := userService.Create(ctx, user); err != nil {
		log.Printf("Create failed: %v", err)
		return
	}
	log.Printf("Created user: %s", user.Name)
	
	// 2. Get by ID
	foundUser, err := userService.Get(ctx, user.ID)
	if err != nil {
		log.Printf("Get failed: %v", err)
		return
	}
	log.Printf("Found user: %s (%s)", foundUser.Name, foundUser.Email)
	
	// 3. Find by email
	emailFilter := map[string]any{"email": "alice@example.com"}
	userByEmail, err := userService.Find(ctx, emailFilter)
	if err != nil {
		log.Printf("Find failed: %v", err)
		return
	}
	log.Printf("Found by email: %s", userByEmail.Name)
	
	// 4. Check if exists
	existsFilter := map[string]any{"email": "alice@example.com"}
	existingUser, exists, err := userService.Exists(ctx, existsFilter)
	if err != nil {
		log.Printf("Exists check failed: %v", err)
		return
	}
	if exists {
		log.Printf("User exists: %s", existingUser.Name)
	}
	
	// 5. Update user
	foundUser.Age = 29
	foundUser.Name = "Alice Smith"
	if err := userService.Update(ctx, user.ID, foundUser); err != nil {
		log.Printf("Update failed: %v", err)
		return
	}
	log.Printf("Updated user: %s (age: %d)", foundUser.Name, foundUser.Age)
	
	// 6. List users with filter
	activeFilter := map[string]any{"is_active": true}
	activeUsers, err := userService.List(ctx, activeFilter)
	if err != nil {
		log.Printf("List failed: %v", err)
		return
	}
	log.Printf("Found %d active users", len(activeUsers))
	
	// 7. Update attributes
	ageFilter := map[string]any{"age": map[string]any{"$lt": 30}}
	updates := map[string]any{"category": "young_adult"}
	if err := userService.UpdateAttributes(ctx, ageFilter, updates); err != nil {
		log.Printf("UpdateAttributes failed: %v", err)
		return
	}
	log.Println("Updated attributes for young users")
	
	// 8. Create index
	emailIndex := map[string]int{"email": 1}
	indexName, err := userService.CreateIndex(ctx, emailIndex, true)
	if err != nil {
		log.Printf("Index creation failed (might already exist): %v", err)
	} else {
		log.Printf("Created index: %s", indexName)
	}
	
	// 9. List all users
	allUsers, err := userService.ListAll(ctx)
	if err != nil {
		log.Printf("ListAll failed: %v", err)
		return
	}
	log.Printf("Total users in database: %d", len(allUsers))
	
	// 10. Clean up - delete the test user
	if err := userService.Delete(ctx, user.ID); err != nil {
		log.Printf("Delete failed: %v", err)
		return
	}
	log.Printf("Deleted test user: %s", user.Name)
	
	log.Println("All CRUD examples completed successfully!")
}
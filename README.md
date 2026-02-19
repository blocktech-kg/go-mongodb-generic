# Go MongoDB Generic CRUD

A generic MongoDB CRUD library for Go that provides type-safe database operations using Go generics.

## Features

- **Type-safe**: Uses Go generics for compile-time type safety
- **Complete CRUD operations**: Create, Read, Update, Delete with various query options
- **Flexible querying**: Support for complex filters and attribute updates
- **Index management**: Easy index creation with unique constraints
- **Automatic timestamps**: Handles `CreatedAt` and `UpdatedAt` fields automatically

## Installation

```bash
go get github.com/blocktech-kg/go-mongodb-generic
```

## Quick Start

### 1. Define Your Model

```go
package main

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
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
```

### 2. Connect to Database

```go
package main

import (
    "context"
    "log"
    "github.com/blocktech-kg/go-mongodb-generic/mongodb"
)

func main() {
    ctx := context.Background()
    
    // Connect to MongoDB
    db, err := mongodb.Connect(ctx, "mongodb://localhost:27017", "myapp")
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    // Get collection
    userCollection := db.Collection("users")
    
    // Create CRUD service
    userService := mongodb.NewGenericObjectDBCtrl[User](userCollection)
    
    // Now you can use userService for CRUD operations
}
```

## Usage Examples

### Create Operations

#### Create a Single Item

```go
func createUser(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    user := &User{
        ID:       primitive.NewObjectID(),
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        IsActive: true,
    }
    
    err := userService.Create(ctx, user)
    if err != nil {
        log.Printf("Error creating user: %v", err)
        return
    }
    
    log.Printf("User created successfully with ID: %s", user.ID.Hex())
}
```

### Read Operations

#### Get by ID

```go
func getUserByID(ctx context.Context, userService mongodb.CRUDDBService[User], id primitive.ObjectID) {
    user, err := userService.Get(ctx, id)
    if err != nil {
        log.Printf("Error getting user: %v", err)
        return
    }
    
    log.Printf("Found user: %+v", user)
}
```

#### Find One by Filter

```go
func findUserByEmail(ctx context.Context, userService mongodb.CRUDDBService[User], email string) {
    filter := map[string]any{
        "email": email,
    }
    
    user, err := userService.Find(ctx, filter)
    if err != nil {
        log.Printf("Error finding user: %v", err)
        return
    }
    
    log.Printf("Found user: %+v", user)
}
```

#### Check if Item Exists

```go
func checkUserExists(ctx context.Context, userService mongodb.CRUDDBService[User], email string) {
    filter := map[string]any{
        "email": email,
    }
    
    user, exists, err := userService.Exists(ctx, filter)
    if err != nil {
        log.Printf("Error checking user existence: %v", err)
        return
    }
    
    if exists {
        log.Printf("User exists: %+v", user)
    } else {
        log.Println("User does not exist")
    }
}
```

#### List All Items

```go
func listAllUsers(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    users, err := userService.ListAll(ctx)
    if err != nil {
        log.Printf("Error listing users: %v", err)
        return
    }
    
    log.Printf("Found %d users:", len(users))
    for _, user := range users {
        log.Printf("- %s (%s)", user.Name, user.Email)
    }
}
```

#### List with Filter

```go
func listActiveUsers(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    filter := map[string]any{
        "is_active": true,
        "age":       map[string]any{"$gte": 18}, // MongoDB query operators
    }
    
    users, err := userService.List(ctx, filter)
    if err != nil {
        log.Printf("Error listing active users: %v", err)
        return
    }
    
    log.Printf("Found %d active users", len(users))
}
```

### Update Operations

#### Update Entire Document

```go
func updateUser(ctx context.Context, userService mongodb.CRUDDBService[User], id primitive.ObjectID) {
    // First get the existing user
    user, err := userService.Get(ctx, id)
    if err != nil {
        log.Printf("Error getting user: %v", err)
        return
    }
    
    // Modify fields
    user.Name = "Jane Doe"
    user.Age = 25
    
    // Update in database
    err = userService.Update(ctx, id, user)
    if err != nil {
        log.Printf("Error updating user: %v", err)
        return
    }
    
    log.Println("User updated successfully")
}
```

#### Update Specific Attributes

```go
func updateUserAttributes(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    // Filter: users older than 30
    filter := map[string]any{
        "age": map[string]any{"$gt": 30},
    }
    
    // Attributes to update
    updates := map[string]any{
        "is_active": false,
        "status":    "senior",
    }
    
    err := userService.UpdateAttributes(ctx, filter, updates)
    if err != nil {
        log.Printf("Error updating user attributes: %v", err)
        return
    }
    
    log.Println("User attributes updated successfully")
}
```

### Delete Operations

#### Delete by ID

```go
func deleteUser(ctx context.Context, userService mongodb.CRUDDBService[User], id primitive.ObjectID) {
    err := userService.Delete(ctx, id)
    if err != nil {
        log.Printf("Error deleting user: %v", err)
        return
    }
    
    log.Println("User deleted successfully")
}
```

#### Delete Multiple Items

```go
func deleteInactiveUsers(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    filter := map[string]any{
        "is_active": false,
    }
    
    err := userService.DeleteRange(ctx, filter)
    if err != nil {
        log.Printf("Error deleting inactive users: %v", err)
        return
    }
    
    log.Println("Inactive users deleted successfully")
}
```

### Index Management

#### Create Simple Index

```go
func createEmailIndex(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    indexKeys := map[string]int{
        "email": 1, // 1 for ascending, -1 for descending
    }
    
    indexName, err := userService.CreateIndex(ctx, indexKeys, true) // unique: true
    if err != nil {
        log.Printf("Error creating index: %v", err)
        return
    }
    
    log.Printf("Index created: %s", indexName)
}
```

#### Create Compound Index

```go
func createCompoundIndex(ctx context.Context, userService mongodb.CRUDDBService[User]) {
    indexKeys := map[string]int{
        "is_active": 1,
        "age":       -1, // descending
        "created_at": 1,
    }
    
    indexName, err := userService.CreateIndex(ctx, indexKeys, false) // unique: false
    if err != nil {
        log.Printf("Error creating compound index: %v", err)
        return
    }
    
    log.Printf("Compound index created: %s", indexName)
}
```

## Complete Example

```go
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
    db, err := mongodb.Connect(ctx, "mongodb://localhost:27017", "myapp")
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    
    // Create service
    userCollection := db.Collection("users")
    userService := mongodb.NewGenericObjectDBCtrl[User](userCollection)
    
    // Create user
    user := &User{
        ID:       primitive.NewObjectID(),
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        IsActive: true,
    }
    
    if err := userService.Create(ctx, user); err != nil {
        log.Fatal("Failed to create user:", err)
    }
    log.Println("✅ User created")
    
    // Get user
    foundUser, err := userService.Get(ctx, user.ID)
    if err != nil {
        log.Fatal("Failed to get user:", err)
    }
    log.Printf("✅ User found: %s", foundUser.Name)
    
    // Update user
    foundUser.Age = 31
    if err := userService.Update(ctx, user.ID, foundUser); err != nil {
        log.Fatal("Failed to update user:", err)
    }
    log.Println("✅ User updated")
    
    // List all users
    users, err := userService.ListAll(ctx)
    if err != nil {
        log.Fatal("Failed to list users:", err)
    }
    log.Printf("✅ Found %d users", len(users))
    
    // Create index
    indexKeys := map[string]int{"email": 1}
    indexName, err := userService.CreateIndex(ctx, indexKeys, true)
    if err != nil {
        log.Fatal("Failed to create index:", err)
    }
    log.Printf("✅ Index created: %s", indexName)
    
    log.Println("🎉 All operations completed successfully!")
}
```

## Advanced Usage

### Using MongoDB Query Operators

```go
// Find users aged between 18 and 65
filter := map[string]any{
    "age": map[string]any{
        "$gte": 18,
        "$lte": 65,
    },
    "is_active": true,
}

users, err := userService.List(ctx, filter)
```

### Working with Embedded Documents

```go
type Address struct {
    Street  string `bson:"street"`
    City    string `bson:"city"`
    Country string `bson:"country"`
}

type UserWithAddress struct {
    ID      primitive.ObjectID `bson:"_id,omitempty"`
    Name    string            `bson:"name"`
    Email   string            `bson:"email"`
    Address Address           `bson:"address"`
    CreatedAt time.Time       `bson:"created_at"`
    UpdatedAt time.Time       `bson:"updated_at"`
}

// Query by nested field
filter := map[string]any{
    "address.city": "New York",
}
```

## Requirements

- Go 1.22+
- MongoDB 4.0+

## Dependencies

- `go.mongodb.org/mongo-driver` - Official MongoDB driver
- `github.com/labstack/gommon` - Logging utilities
- `github.com/pkg/errors` - Enhanced error handling

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
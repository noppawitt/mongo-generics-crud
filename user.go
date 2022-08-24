package genericscrud

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name  string `bson:"name,omitempty" json:"name"`
	Email string `bson:"email,omitempty" json:"email"`
}

type UserService struct {
	MongoService[User]
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		MongoService: MongoService[User]{
			db:             db,
			collectionName: "users",
		},
	}
}

func (s *UserService) Update(ctx context.Context, id string, update *User) (*Model[User], error) {
	// User email can't be update
	update.Email = ""
	return s.MongoService.Update(ctx, id, update)
}

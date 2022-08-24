package main

import (
	"context"
	"encoding/json"
	"fmt"
	"genericscrud"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	opts := options.Client().
		ApplyURI("mongodb://admin:password@127.0.0.1:27017/admin")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	db := client.Database("generics")

	userSvc := genericscrud.NewUserService(db)

	create := &genericscrud.User{
		Name:  "Noppawit",
		Email: "noppawit@gmail.com",
	}

	user, err := userSvc.Create(ctx, create)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("created user:", jsonStr(user))

	id := user.ID.Hex()

	update := &genericscrud.User{
		Name:  "Peeja",
		Email: "peeja@gmail.com",
	}

	user, err = userSvc.Update(ctx, id, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("updated user:", jsonStr(user))

	user, err = userSvc.Find(ctx, id)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("found user:", jsonStr(user))

	if err = userSvc.Delete(ctx, id); err != nil {
		log.Fatal(err)
	}

	fmt.Println("deleted user id: ", id)

	_, err = userSvc.Find(ctx, id)
	if err != nil {
		fmt.Println("user not found")
	}
}

func jsonStr(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(data)
}

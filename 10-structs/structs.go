package main

import (
	"fmt"
	"time"
)

type user struct {
	firstName string
	lastName  string
	createdAt time.Time
}

//* struct embedding

type Admin struct{
	email string
	password string
	user
}

func outputUserDetails(u user) {
	fmt.Println(u.firstName, u.lastName)
}

func main() {

	var user1 user
	user1 = user{
		firstName: "John",
		lastName:  "Doe",
	}
	outputUserDetails((user1))
	// fmt.Println(user1)
}

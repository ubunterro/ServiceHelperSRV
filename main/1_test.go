package main

import (
	"fmt"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	hash := Hash{}
	generated, _ := hash.Generate("admin")
	fmt.Println(generated)
	fmt.Println(hash.Compare(generated, "admin"))
}

func TestAuth(t *testing.T) {
	authenticateUser("admin", "admin")
}

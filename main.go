package main

import (
	"github.com/gophercises/ex1"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	ex1.Main()
}

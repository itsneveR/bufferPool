package main

import (
	"fmt"
	"sync"
)

var pool *sync.Pool

type Person struct {
	Name string
}

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Pool was empty: creating a new object")
			return new(Person)
		},
	}
}

func main() {
	objectFromPool := pool.Get().(*Person)

	fmt.Println("1.get an object from Pool", objectFromPool)

	objectFromPool.Name = "mamad"

	pool.Put(objectFromPool)

	fmt.Println("2.Get Pool Object:  ", pool.Get().(*Person))
	fmt.Println("3.Get Pool Object:  ", pool.Get().(*Person))
}

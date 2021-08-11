package main

import (
	"fmt"
	"log"

	"github.com/brittlesoft/go-uuidv7/pkg/uuidv7"
)

func main() {
	us, err := uuidv7.NewUuidv7Source()
	if err != nil {
		log.Fatal(err)
	}
	u, err := us.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.String())
	fmt.Println(u)
	fmt.Printf("original ts: %f\n", u.Ts())
	du := uuidv7.NewDecodedUuidv7(u.B)
	fmt.Printf("%v\n", du)
	fmt.Printf("Ts: %f\n", du.Ts)
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/didchain/didCard-go/card"
)

func sample1() {
	card, err := card.NewSimpleCard("123")
	if err != nil {
		panic(err)
	}
	fmt.Print(card)

	byts, err := json.Marshal(card)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(byts))
}

func main() {
	sample1()
}

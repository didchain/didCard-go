package ios

import (
	"encoding/json"
	"github.com/didchain/didCard-go/card"
)

func NewCard(auth string) []byte {
	card, err := card.NewSimpleCard(auth)
	if err != nil {
		return nil
	}

	encodedFile, err := json.MarshalIndent(card, "", "\t")
	if err != nil {
		return nil
	}

	return encodedFile
}

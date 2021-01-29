package android

import (
	"encoding/json"
	"github.com/didchain/didCard-go/card"
)

var _cardInst card.DIDCard = nil

func NewCard(auth string) []byte {
	card, err := card.NewSimpleCard(auth)
	if err != nil {
		return nil
	}

	encodedFile, err := json.MarshalIndent(card, "", "\t")
	if err != nil {
		return nil
	}

	_cardInst = card
	return encodedFile
}

func LoadCard(jsonStr string) bool {
	card, err := card.ParseCard(jsonStr)
	if err != nil {
		return false
	}
	_cardInst = card
	return true
}

func LoadCardByPath(fullPath string) bool {
	card, err := card.LoadCard(fullPath)
	if err != nil {
		return false
	}
	_cardInst = card
	return true
}

func Import(auth, jsonStr string) []byte {
	card, err := card.ParseCard(jsonStr)
	if err != nil {
		return nil
	}
	if !card.Open(auth) {
		return nil
	}
	_cardInst = card
	return []byte(card.PublicKey())
}

func Open(auth string) bool {
	if _cardInst == nil {
		return false
	}

	return _cardInst.Open(auth)
}

func IsOpen() bool {
	if _cardInst == nil {
		return false
	}

	return _cardInst.IsOpen()
}

func SignByPassword(msg, auth string) []byte {
	if _cardInst == nil {
		return nil
	}
	_cardInst.Open(auth)
	return _cardInst.Sign([]byte(msg))
}

func Sign(msg string) []byte {
	return _cardInst.Sign([]byte(msg))
}

func Verify(pub, msg, sig []byte) bool {
	return card.Verify(pub, sig, msg)
}

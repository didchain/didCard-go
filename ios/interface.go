package iosLib

import (
	"github.com/didchain/didCard-go/account"
)

var _cardInst account.Wallet = nil

func NewCard(auth string) []byte {
	card, err := account.NewWallet(auth)
	if err != nil {
		return nil
	}

	encodedFile := card.Bytes()

	_cardInst = card
	return encodedFile
}

func LoadCard(jsonStr string) bool {
	card, err := account.LoadWalletByData(jsonStr)
	if err != nil {
		return false
	}
	_cardInst = card
	return true
}

func LoadCardByPath(fullPath string) bool {
	card, err := account.LoadWallet(fullPath)
	if err != nil {
		return false
	}
	_cardInst = card
	return true
}

func Import(auth, jsonStr string) []byte {
	card, err := account.LoadWalletByData(jsonStr)
	if err != nil {
		return nil
	}
	if err:=card.Open(auth);err!=nil {
		return nil
	}
	_cardInst = card
	return []byte(card.Did())
}

func Open(auth string) bool {
	if _cardInst == nil {
		return false
	}

	if err:= _cardInst.Open(auth);err!=nil{
		return false
	}

	return true
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
	return account.VerifySig(account.ConvertToID2(pub), sig, msg)
}

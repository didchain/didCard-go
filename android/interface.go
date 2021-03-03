package androidgolib

import (
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/didchain/didCard-go/account"
)

var _cardInst account.Wallet = nil

func NewCard(auth string) ([]byte, error) {
	card, err := account.NewWallet(auth)
	if err != nil {
		return nil, err
	}

	encodedFile := card.Bytes()

	_cardInst = card
	return encodedFile, nil
}

func LoadCard(jsonStr string) (bool, error) {
	card, err := account.LoadWalletByData(jsonStr)
	if err != nil {
		return false, err
	}
	_cardInst = card
	return true, nil
}

func LoadCardByPath(fullPath string) (bool, error) {
	card, err := account.LoadWallet(fullPath)
	if err != nil {
		return false, err
	}
	_cardInst = card
	return true, nil
}

func Import(auth, jsonStr string) ([]byte, error) {
	card, err := account.LoadWalletByData(jsonStr)
	if err != nil {
		return nil, err
	}
	if err := card.Open(auth); err != nil {
		return nil, err
	}
	_cardInst = card
	return []byte(card.Did()), nil
}

func Open(auth string) error {
	if _cardInst == nil {
		return errors.New("no card instance")
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

func Sign(msg string) string {
	sig:= _cardInst.Sign([]byte(msg))
	return base58.Encode(sig)
}

func Verify(pub, msg, sig []byte) bool {
	return account.VerifySig(account.ConvertToID2(pub), sig, msg)
}

func Close() {
	if _cardInst == nil {
		return
	}
	_cardInst.Close()
}

package iosLib

import (
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
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
	if err := card.Open(auth); err != nil {
		return nil
	}
	_cardInst = card
	return []byte(card.Did())
}

func Open(auth string) bool {
	if _cardInst == nil {
		return false
	}

	if err := _cardInst.Open(auth); err != nil {
		return false
	}

	return true
}

type DeriveKey struct {
	AesKey    string `json:"aes_key"`
	CardBytes []byte `json:"card_bytes"`
}

//AES Key is generate by a new salt, need to save it
func DeriveAesKey(auth string) string {
	if _cardInst == nil {
		return ""
	}

	aesKey, err := _cardInst.DriveAESKey(auth)
	if err != nil {
		return ""
	}

	return aesKey
}

func OpenWithAesKey(aesKey string) string {
	if _cardInst == nil {
		return "no card instance"
	}

	err := _cardInst.OpenWithAesKey(aesKey)
	if err != nil {
		return "open failed"
	}
	return ""
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
	sig := _cardInst.Sign([]byte(msg))
	return base58.Encode(sig)
}

func Verify(pub []byte, msg interface{}, sig string) bool {
	sigbytes := base58.Decode(sig)
	return account.VerifySig(account.ConvertToID2(pub), sigbytes, msg)
}

func SignMessage(did string, latitude, longitude float64, timestamp int64) string {
	msg := struct {
		DID       string  `json:"did"` ///public key in string
		TimeStamp int64   `json:"time_stamp"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}{}

	msg.DID = did
	msg.TimeStamp = timestamp
	msg.Latitude = latitude
	msg.Longitude = longitude

	j, _ := json.Marshal(msg)

	return string(j)

}

func SignUserLoginMessage(did string, randomBytes string, loginUrl string) string{
	msg:= struct {
		DID string `json:"did"`
		RandomBytes string `json:"random_bytes"`
		LoginUrl string `json:"login_url"`
	}{}

	msg.DID = did
	msg.RandomBytes = randomBytes
	msg.LoginUrl = loginUrl

	j,_:=json.Marshal(msg)

	return string(j)

}

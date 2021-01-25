package card

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	bls12 "github.com/herumi/bls-eth-go-binary/bls"
	ksv4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	"io/ioutil"
)

func init() {
	_ = bls12.Init(bls12.BLS12_381)
	_ = bls12.SetETHmode(bls12.EthModeDraft07)
}

type DIDCard interface {
	PublicKey() string
	Sign(msg []byte) []byte
	Open(auth string) bool
	IsOpen() bool
}

type SimpleCard struct {
	Crypto  map[string]interface{} `json:"crypto"`
	ID      string                 `json:"uuid"`
	PubKey  string                 `json:"pubkey"`
	Version uint                   `json:"version"`
	priKey  *bls12.SecretKey
}

func (sc *SimpleCard) PublicKey() string {
	return sc.PubKey
}

func (sc *SimpleCard) Sign(msg []byte) []byte {
	if sc.priKey == nil {
		return nil
	}

	signature := sc.priKey.SignByte(msg)
	return signature.Serialize()
}

func (sc *SimpleCard) Open(auth string) bool {

	if sc.priKey != nil {
		return true
	}
	decrypt := ksv4.New()
	priKey, err := decrypt.Decrypt(sc.Crypto, auth)
	if err != nil {
		return false
	}
	secKey := &bls12.SecretKey{}
	if err := secKey.Deserialize(priKey); err != nil {
		return false
	}

	sc.priKey = secKey
	return true
}

func (sc *SimpleCard) IsOpen() bool {
	return sc.priKey != nil
}

func NewSimpleCard(auth string) (DIDCard, error) {
	secKey := &bls12.SecretKey{}
	secKey.SetByCSPRNG()
	if secKey.IsZero() {
		return nil, errors.New("generated a zero secret key")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	pubKey := string(secKey.GetPublicKey().Serialize())

	encryptor := ksv4.New()
	cryptoFields, err := encryptor.Encrypt(secKey.Serialize(), auth)
	card := &SimpleCard{
		Crypto:  cryptoFields,
		ID:      id.String(),
		Version: encryptor.Version(),
		PubKey:  pubKey,
		priKey:  secKey,
	}

	return card, nil
}

func ParseCard(jsonStr string) (DIDCard, error) {
	var card = &SimpleCard{}
	if err := json.Unmarshal([]byte(jsonStr), &card); err != nil {
		return nil, err
	}
	return card, nil
}

func LoadCard(fullPath string) (DIDCard, error) {
	jsonStr, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	var card = &SimpleCard{}
	if err := json.Unmarshal([]byte(jsonStr), &card); err != nil {
		return nil, err
	}
	return card, nil
}

func Verify(pub, sig, msg []byte) bool {
	var pubKey bls12.PublicKey
	if err := pubKey.Deserialize(pub); err != nil {
		return false
	}

	var signature bls12.Sign
	if err := signature.Deserialize(sig); err != nil {
		return false
	}
	if pubKey.IsZero() {
		return false
	}

	return signature.VerifyByte(&pubKey, msg)
}

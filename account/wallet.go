package account

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
)

const (
	WalletVersion = 1
)

type Wallet interface {
	PrivKey() ed25519.PrivateKey
	Did() ID

	SignJson(v interface{}) []byte
	Sign(v []byte) []byte
	VerifySig(message, signature []byte) bool
	VerifySigObj(obj interface{}, signature []byte) bool

	Open(auth string) error

	DriveAESKey(auth string) (string,error)
	OpenWithAesKey(key string) error

	IsOpen() bool
	SaveToPath(wPath string) error
	String() string
	Bytes() []byte
	Close()
}

type WalletKey struct {
	PriKey  ed25519.PrivateKey
}

type PWallet struct {
	Version   int                 `json:"version"`
	DidAddr   ID 				  `json:"did"`
	Salt      string              `json:"salt,omitempty"`
	CipherTxt string 			  `json:"cipher_txt"`
	key       *WalletKey          `json:"-"`
}

func NewWallet(auth string) (Wallet, error) {

	pub, pri, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		logger.Errorf("Error generate network key: %v", err)
		return nil, err
	}
	cipherTxt,_, err := encryptPriKey(pri, pub[:KP.S], auth)
	if err != nil {
		logger.Errorf("encrypt wallet err:%f", err)
		return nil, err
	}

	obj := &PWallet{
		Version:   WalletVersion,
		DidAddr:		ConvertToID2(pub),
		CipherTxt: cipherTxt,
		key: &WalletKey{
			PriKey:  pri,
		},
	}
	return obj, nil
}

func encryptPriKey(priKey ed25519.PrivateKey, salt []byte, auth string) (string,string, error) {
	aesKey, err := AESKey(salt, auth)
	if err != nil {
		logger.Warning("error to generate aes key:->", err)
		return "", "",err
	}
	cipher, err := Encrypt(aesKey, priKey[:])
	if err != nil {
		logger.Warning("error to encrypt the raw private key:->", err)
		return "","", err
	}
	return base58.Encode(cipher),base58.Encode(aesKey), nil
}

func decryptPriKey(salt []byte, cpTxt, auth string) (ed25519.PrivateKey, error) {
	aesKey, err := AESKey(salt, auth)
	if err != nil {
		return nil, err
	}
	//fmt.Println("aes key == >: ",hex.EncodeToString(aesKey))
	cipherByte := base58.Decode(cpTxt)
	//fmt.Println("cipher base16 == >: ",hex.EncodeToString(cipherByte))
	privBytes := make([]byte, len(cipherByte))
	copy(privBytes, cipherByte)
	return Decrypt(aesKey, privBytes)
}

func aesKey(salt []byte,auth string) (string,error)  {
	 aesk,err:= AESKey(salt, auth)
	 if err!=nil{
	 	return "", err
	 }

	 return base58.Encode(aesk),nil
}

func decryptPrivKeyByAesKey(aesKey []byte,cpTxt string) (ed25519.PrivateKey,error)  {
	cipherByte := base58.Decode(cpTxt)
	//fmt.Println("cipher base16 == >: ",hex.EncodeToString(cipherByte))
	privBytes := make([]byte, len(cipherByte))
	copy(privBytes, cipherByte)
	return Decrypt(aesKey, privBytes)
}


func LoadWallet(wPath string) (Wallet, error) {
	jsonStr, err := ioutil.ReadFile(wPath)
	if err != nil {
		return nil, err
	}

	w := new(PWallet)
	if err := json.Unmarshal(jsonStr, w); err != nil {
		logger.Errorf("error to parse wallet :%s", err)
		return nil, err	}
	return w, nil
}

func LoadWalletByData(jsonStr string) (Wallet, error) {
	w := new(PWallet)
	if err := json.Unmarshal([]byte(jsonStr), w); err != nil {
		logger.Errorf("error to parse wallet :%s", err)
		return nil, err
	}
	return w, nil
}

func VerifySig(didAddr ID, sig []byte, v interface{}) bool {
	data, err := json.Marshal(v)
	if err != nil {
		return false
	}

	return ed25519.Verify(didAddr.ToPubKey(), data, sig)
}

package account

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
)



func (pw *PWallet) PrivKey() ed25519.PrivateKey {
	return pw.key.PriKey
}

func (pw *PWallet)Did() ID {
	return pw.DidAddr
}

func (pw *PWallet) SignJson(v interface{}) []byte {
	rawBytes, _ := json.Marshal(v)
	return ed25519.Sign(pw.key.PriKey, rawBytes)
}

func (pw *PWallet) Sign(v []byte) []byte {
	return ed25519.Sign(pw.key.PriKey, v)
}

func (pw *PWallet) IsOpen() bool {
	return pw.key != nil
}

func (pw *PWallet) SaveToPath(wPath string) error {
	bytes, err := json.MarshalIndent(pw, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(wPath, bytes, 0644)
}

func (pw *PWallet) Open(auth string) error {
	var salt []byte
	if pw.Salt != ""{
		salt=base58.Decode(pw.Salt)
	}else{
		pk := pw.Did().ToPubKey()
		salt = pk[:KP.S]
	}

	privkey, err := decryptPriKey(salt, pw.CipherTxt, auth)
	if err != nil {
		return err
	}

	pubk:=privkey.Public()
	id := ConvertToID2(pubk.(ed25519.PublicKey))
	if pw.DidAddr.String() != id.String(){
		return errors.New("open failed, may be password is not correct")
	}

	key := &WalletKey{
		PriKey: privkey,
	}
	pw.key = key
	return nil
}

func (pw *PWallet)OpenWithAesKey(aeskey string) error{
	privkey, err := decryptPrivKeyByAesKey(base58.Decode(aeskey), pw.CipherTxt)
	if err != nil {
		return err
	}

	pubk:=privkey.Public()
	id := ConvertToID2(pubk.(ed25519.PublicKey))
	if pw.DidAddr.String() != id.String(){
		return errors.New("open failed, may be password is not correct")
	}

	key := &WalletKey{
		PriKey: privkey,
	}
	pw.key = key
	return nil
}

func (pw *PWallet)DriveAESKey(auth string) (string,error) {
	if pw.key != nil{
		err:=pw.Open(auth)
		if err!=nil{
			return "", err
		}
	}

	//salt:=make([]byte,KP.S)
	//if _, err := io.ReadFull(rand.Reader, salt); err != nil {
	//	return "", err
	//}

	//salt:=pw.DidAddr.ToPubKey()
	//
	//cipherTxt,aesk, err:=encryptPriKey(pw.key.PriKey,salt[:KP.S],auth)
	//if err!=nil{
	//	return "", err
	//}
	//
	//pw.CipherTxt = cipherTxt
	////pw.Salt = base58.Encode(salt)
	salt:=pw.DidAddr.ToPubKey()
	aesk,err:=aesKey(salt[:KP.S],auth)
	if err!=nil{
		return "", err
	}

	return aesk,nil

}


func (pw *PWallet) Close() {
	pw.key = nil
}

func (pw *PWallet) String() string {
	return string(pw.Bytes())
}

func (pw *PWallet)Bytes() []byte  {
	b, e := json.MarshalIndent(pw," ","\t")
	if e != nil {
		return nil
	}



	return b
}


func (pw *PWallet)VerifySig(message, signature []byte) bool  {
	return ed25519.Verify(pw.DidAddr.ToPubKey(), message, signature)
}

func (pw *PWallet)VerifySigObj(obj interface{}, signature []byte) bool  {
	return VerifySig(pw.DidAddr,signature,obj)
}


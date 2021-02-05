package account

import (
	"crypto/ed25519"
	"encoding/json"
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

	privkey, err := decryptPriKey(pw.DidAddr, pw.CipherTxt, auth)
	if err != nil {
		return err
	}
	key := &WalletKey{
		PriKey: privkey,
	}
	pw.key = key
	return nil
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


package card

import "encoding/hex"

const PREFIX = "DID"

type DID string

func Parse(rawPub []byte) DID {
	return DID(hex.EncodeToString(rawPub))
}

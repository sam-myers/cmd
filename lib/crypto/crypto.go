package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"

	"github.com/gliderlabs/cmd/lib/daemon"
	"github.com/gliderlabs/comlab/pkg/com"
)

const LocalDevKey = "localdev"

func init() {
	com.Register("crypto", &Component{},
		com.Option("secret_key", LocalDevKey, "key used to encrypt and decrypt secrets"))
}

type Component struct {
}

func (Component) AppPreStart() error {
	key := com.GetString("secret_key")
	if key == "" {
		return errors.New("crypto: secret_key missing")
	}
	if key == LocalDevKey && !daemon.LocalMode() {
		return errors.New("crypto: secret_key must not be default")
	}
	copy(secretKey[:], []byte(key))

	return nil
}

var secretKey [32]byte

func Encrypt(msg string) (string, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	box := secretbox.Seal(nonce[:], []byte(msg), &nonce, &secretKey)
	return base64.StdEncoding.EncodeToString(box), nil
}

func Decrypt(box string) string {
	enc, err := base64.StdEncoding.DecodeString(box)
	if err != nil {
		return ""
	}
	if len(enc) < 25 {
		return ""
	}
	var nonce [24]byte
	copy(nonce[:], enc[:24])
	decrypted, ok := secretbox.Open([]byte{}, enc[24:], &nonce, &secretKey)
	if !ok {
		return ""
	}
	return string(decrypted)
}

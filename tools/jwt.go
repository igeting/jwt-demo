package tools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"strings"
	"time"
)

const (
	ALG = "HS256"
	JWT = "JWT"
)

var (
	secret                = "secret"
	expires time.Duration = time.Minute * 30
)

type Header struct {
	Alg string `json:"alg,omitempty"`
	Typ string `json:"typ,omitempty"`
}

type Payload struct {
	Id        string `json:"jti,omitempty"`
	Subject   string `json:"sub,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
}

type Signature struct {
	header        Header
	payload       Payload
	headerString  string
	payloadString string
	signString    string
}

func (s *Signature) Sign(secret string) {
	s.headerString = getBaseString(s.header)
	s.payloadString = getBaseString(s.payload)

	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(secret))

	h.Write([]byte(s.headerString + "." + s.payloadString))
	res := h.Sum(nil)

	s.signString = base64.URLEncoding.EncodeToString(res)
}

func (s *Signature) SignString() string {
	return fmt.Sprintf("%s.%s.%s", s.headerString, s.payloadString, s.signString)
}

func getBaseString(o interface{}) string {
	js, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(js)
}

func Generate(id string) string {
	s := Signature{
		header: Header{
			Alg: ALG,
			Typ: JWT,
		},
		payload: Payload{
			Id:        id,
			ExpiresAt: time.Now().Add(expires).Unix(),
		},
	}
	s.Sign(secret)
	return s.SignString()
}

func Verify(s string) bool {
	sp := strings.Split(s, ".")
	if len(sp) == 3 {
		h := sp[0]
		p := sp[1]
		s := sp[2]

		pp, err := base64.URLEncoding.DecodeString(p)
		if err != nil {
			panic(err)
		}
		pl := new(Payload)
		err = json.Unmarshal(pp, pl)
		if err != nil {
			panic(err)
		}
		expires := time.Now().Before(time.Unix(pl.ExpiresAt, 0))

		if expires {
			sh := hmac.New(func() hash.Hash {
				return sha256.New()
			}, []byte(secret))

			sh.Write([]byte(h + "." + p))
			ss := sh.Sum(nil)

			return base64.URLEncoding.EncodeToString(ss) == s
		}
	}
	return false
}

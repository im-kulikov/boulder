// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package wfe

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cactus/go-statsd-client/statsd"

	jose "github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/square/go-jose"
	"github.com/letsencrypt/boulder/core"

	"github.com/letsencrypt/boulder/ra"
	"github.com/letsencrypt/boulder/test"
)

const (
	test1KeyPublicJSON = `
	{
		"kty":"RSA",
		"n":"yNWVhtYEKJR21y9xsHV-PD_bYwbXSeNuFal46xYxVfRL5mqha7vttvjB_vc7Xg2RvgCxHPCqoxgMPTzHrZT75LjCwIW2K_klBYN8oYvTwwmeSkAz6ut7ZxPv-nZaT5TJhGk0NT2kh_zSpdriEJ_3vW-mqxYbbBmpvHqsa1_zx9fSuHYctAZJWzxzUZXykbWMWQZpEiE0J4ajj51fInEzVn7VxV-mzfMyboQjujPh7aNJxAWSq4oQEJJDgWwSh9leyoJoPpONHxh5nEE5AjE01FkGICSxjpZsF-w8hOTI3XXohUdu29Se26k2B0PolDSuj0GIQU6-W9TdLXSjBb2SpQ",
		"e":"AAEAAQ"
	}`

	test1KeyPrivatePEM = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyNWVhtYEKJR21y9xsHV+PD/bYwbXSeNuFal46xYxVfRL5mqh
a7vttvjB/vc7Xg2RvgCxHPCqoxgMPTzHrZT75LjCwIW2K/klBYN8oYvTwwmeSkAz
6ut7ZxPv+nZaT5TJhGk0NT2kh/zSpdriEJ/3vW+mqxYbbBmpvHqsa1/zx9fSuHYc
tAZJWzxzUZXykbWMWQZpEiE0J4ajj51fInEzVn7VxV+mzfMyboQjujPh7aNJxAWS
q4oQEJJDgWwSh9leyoJoPpONHxh5nEE5AjE01FkGICSxjpZsF+w8hOTI3XXohUdu
29Se26k2B0PolDSuj0GIQU6+W9TdLXSjBb2SpQIDAQABAoIBAHw58SXYV/Yp72Cn
jjFSW+U0sqWMY7rmnP91NsBjl9zNIe3C41pagm39bTIjB2vkBNR8ZRG7pDEB/QAc
Cn9Keo094+lmTArjL407ien7Ld+koW7YS8TyKADYikZo0vAK3qOy14JfQNiFAF9r
Bw61hG5/E58cK5YwQZe+YcyBK6/erM8fLrJEyw4CV49wWdq/QqmNYU1dx4OExAkl
KMfvYXpjzpvyyTnZuS4RONfHsO8+JTyJVm+lUv2x+bTce6R4W++UhQY38HakJ0x3
XRfXooRv1Bletu5OFlpXfTSGz/5gqsfemLSr5UHncsCcFMgoFBsk2t/5BVukBgC7
PnHrAjkCgYEA887PRr7zu3OnaXKxylW5U5t4LzdMQLpslVW7cLPD4Y08Rye6fF5s
O/jK1DNFXIoUB7iS30qR7HtaOnveW6H8/kTmMv/YAhLO7PAbRPCKxxcKtniEmP1x
ADH0tF2g5uHB/zeZhCo9qJiF0QaJynvSyvSyJFmY6lLvYZsAW+C+PesCgYEA0uCi
Q8rXLzLpfH2NKlLwlJTi5JjE+xjbabgja0YySwsKzSlmvYJqdnE2Xk+FHj7TCnSK
KUzQKR7+rEk5flwEAf+aCCNh3W4+Hp9MmrdAcCn8ZsKmEW/o7oDzwiAkRCmLw/ck
RSFJZpvFoxEg15riT37EjOJ4LBZ6SwedsoGA/a8CgYEA2Ve4sdGSR73/NOKZGc23
q4/B4R2DrYRDPhEySnMGoPCeFrSU6z/lbsUIU4jtQWSaHJPu4n2AfncsZUx9WeSb
OzTCnh4zOw33R4N4W8mvfXHODAJ9+kCc1tax1YRN5uTEYzb2dLqPQtfNGxygA1DF
BkaC9CKnTeTnH3TlKgK8tUcCgYB7J1lcgh+9ntwhKinBKAL8ox8HJfkUM+YgDbwR
sEM69E3wl1c7IekPFvsLhSFXEpWpq3nsuMFw4nsVHwaGtzJYAHByhEdpTDLXK21P
heoKF1sioFbgJB1C/Ohe3OqRLDpFzhXOkawOUrbPjvdBM2Erz/r11GUeSlpNazs7
vsoYXQKBgFwFM1IHmqOf8a2wEFa/a++2y/WT7ZG9nNw1W36S3P04K4lGRNRS2Y/S
snYiqxD9nL7pVqQP2Qbqbn0yD6d3G5/7r86F7Wu2pihM8g6oyMZ3qZvvRIBvKfWo
eROL1ve1vmQF3kjrMPhhK2kr6qdWnTE5XlPllVSZFQenSTzj98AO
-----END RSA PRIVATE KEY-----
`

	test2KeyPublicJSON = `
	{
		"kty":"RSA",
		"n":"m5Cpx3vZ0CjATirDpbILvq78fm3Dv5RBkO1VLWFmJj5Mb54vc9oYZWc1V1k-LJoESuuPHhaNO2Eu8T9tslQWcZSzr5NImxAwMk970gVQa-Hqv-Jr6xstrBpq7TKpXHTx2FnfA2wQrfIQSlBXu0t4jdUOr3oJh-QXvma8nLITdtjpC0AZNtqd0QkRJX_90SaNrl18Rr_0JrBH9ZmUSFcf3mo_BtL0Gx0jE3n-iwCI8rQtfyVP__9-n__r4IhalKLzaeio6o-qrdemh0EZgjKGCS1_RpTIeArkO8uia1KgOq-z-GfemKEm4s07WO_a0_9dLqbvpnyyZvUi405m3vGDfQ",
		"e":"AAEAAQ"
	}`
)

var RandomCert []byte = []byte{48, 130, 1, 186, 48, 130, 1, 102, 160, 3, 2, 1, 2, 2, 2, 6, 117, 48, 11, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 11, 48, 52, 49, 16, 48, 14, 6, 3, 85, 4, 6, 19, 7, 65, 117, 115, 116, 114, 105, 97, 49, 17, 48, 15, 6, 3, 85, 4, 10, 19, 8, 72, 101, 97, 112, 108, 111, 99, 107, 49, 13, 48, 11, 6, 3, 85, 4, 11, 19, 4, 84, 72, 79, 82, 48, 30, 23, 13, 48, 57, 49, 49, 49, 48, 50, 51, 48, 48, 48, 48, 90, 23, 13, 49, 57, 49, 49, 49, 48, 50, 51, 48, 48, 48, 48, 90, 48, 52, 49, 16, 48, 14, 6, 3, 85, 4, 6, 19, 7, 65, 117, 115, 116, 114, 105, 97, 49, 17, 48, 15, 6, 3, 85, 4, 10, 19, 8, 72, 101, 97, 112, 108, 111, 99, 107, 49, 13, 48, 11, 6, 3, 85, 4, 11, 19, 4, 84, 72, 79, 82, 48, 92, 48, 13, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 1, 5, 0, 3, 75, 0, 48, 72, 2, 65, 0, 164, 165, 0, 127, 188, 218, 118, 45, 136, 247, 42, 198, 96, 21, 113, 217, 137, 179, 151, 167, 194, 251, 194, 96, 36, 124, 119, 62, 11, 149, 192, 157, 16, 11, 67, 23, 15, 224, 108, 74, 100, 149, 41, 7, 139, 144, 249, 230, 17, 168, 156, 164, 112, 179, 81, 63, 73, 251, 180, 46, 85, 92, 108, 217, 2, 3, 1, 0, 1, 163, 100, 48, 98, 48, 14, 6, 3, 85, 29, 15, 1, 1, 255, 4, 4, 3, 2, 0, 132, 48, 29, 6, 3, 85, 29, 37, 4, 22, 48, 20, 6, 8, 43, 6, 1, 5, 5, 7, 3, 2, 6, 8, 43, 6, 1, 5, 5, 7, 3, 1, 48, 15, 6, 3, 85, 29, 19, 1, 1, 255, 4, 5, 48, 3, 1, 1, 255, 48, 14, 6, 3, 85, 29, 14, 4, 7, 4, 5, 1, 2, 3, 4, 5, 48, 16, 6, 3, 85, 29, 35, 4, 9, 48, 7, 128, 5, 1, 2, 3, 4, 5, 48, 11, 6, 9, 42, 134, 72, 134, 247, 13, 1, 1, 11, 3, 65, 0, 33, 47, 156, 173, 234, 114, 157, 142, 179, 109, 157, 252, 3, 95, 47, 22, 189, 86, 26, 11, 48, 223, 165, 248, 56, 192, 204, 113, 40, 158, 229, 250, 238, 214, 13, 92, 62, 177, 168, 122, 236, 172, 130, 120, 85, 220, 194, 251, 137, 221, 191, 224, 94, 188, 220, 102, 125, 15, 46, 49, 204, 241, 217, 13}

type MockSA struct {
	// empty
}

func (sa *MockSA) GetRegistration(id int64) (core.Registration, error) {
	if id == 100 {
		// Tag meaning "Missing"
		return core.Registration{}, errors.New("missing")
	}
	if id == 101 {
		// Tag meaning "Malformed"
		return core.Registration{}, nil
	}

	keyJSON := []byte(test1KeyPublicJSON)
	var parsedKey jose.JsonWebKey
	parsedKey.UnmarshalJSON(keyJSON)

	return core.Registration{ID: id, Key: parsedKey, Agreement: "yup"}, nil
}

func (sa *MockSA) GetRegistrationByKey(jwk jose.JsonWebKey) (core.Registration, error) {
	var test1KeyPublic jose.JsonWebKey
	var test2KeyPublic jose.JsonWebKey
	test1KeyPublic.UnmarshalJSON([]byte(test1KeyPublicJSON))
	test2KeyPublic.UnmarshalJSON([]byte(test2KeyPublicJSON))

	if core.KeyDigestEquals(jwk, test1KeyPublic) {
		return core.Registration{ID: 1, Key: jwk, Agreement: "yup"}, nil
	}

	if core.KeyDigestEquals(jwk, test2KeyPublic) {
		// No key found
		return core.Registration{ID: 2}, sql.ErrNoRows
	}

	// Return a fake registration
	return core.Registration{ID: 1, Agreement: "yup"}, nil
}

func (sa *MockSA) GetAuthorization(id string) (core.Authorization, error) {
	if id == "valid" {
		return core.Authorization{Status: core.StatusValid, RegistrationID: 1, Expires: time.Now().AddDate(0, 0, 1), Identifier: core.AcmeIdentifier{Type: "dns", Value: "meep.com"}}, nil
	}
	return core.Authorization{}, nil
}

func (sa *MockSA) GetCertificate(string) ([]byte, error) {
	return []byte{}, nil
}

func (sa *MockSA) GetCertificateByShortSerial(string) ([]byte, error) {
	return []byte{}, nil
}

func (sa *MockSA) GetCertificateStatus(string) (core.CertificateStatus, error) {
	return core.CertificateStatus{}, nil
}

func (sa *MockSA) AlreadyDeniedCSR([]string) (bool, error) {
	return false, nil
}

func (sa *MockSA) AddCertificate(certDER []byte, regID int64) (digest string, err error) {
	return
}

func (sa *MockSA) FinalizeAuthorization(authz core.Authorization) (err error) {
	return
}

func (sa *MockSA) MarkCertificateRevoked(serial string, ocspResponse []byte, reasonCode int) (err error) {
	return
}

func (sa *MockSA) NewPendingAuthorization(authz core.Authorization) (output core.Authorization, err error) {
	return
}

func (sa *MockSA) NewRegistration(reg core.Registration) (regR core.Registration, err error) {
	return
}

func (sa *MockSA) UpdatePendingAuthorization(authz core.Authorization) (err error) {
	return
}

func (sa *MockSA) UpdateRegistration(reg core.Registration) (err error) {
	return
}

type MockRegistrationAuthority struct{}

func (ra *MockRegistrationAuthority) NewRegistration(reg core.Registration) (core.Registration, error) {
	return reg, nil
}

func (ra *MockRegistrationAuthority) NewAuthorization(authz core.Authorization, regID int64) (core.Authorization, error) {
	authz.RegistrationID = regID
	return authz, nil
}

func (ra *MockRegistrationAuthority) NewCertificate(req core.CertificateRequest, regID int64) (core.Certificate, error) {
	return core.Certificate{}, nil
}

func (ra *MockRegistrationAuthority) UpdateRegistration(reg core.Registration, updated core.Registration) (core.Registration, error) {
	return reg, nil
}

func (ra *MockRegistrationAuthority) UpdateAuthorization(authz core.Authorization, foo int, challenge core.Challenge) (core.Authorization, error) {
	return authz, nil
}

func (ra *MockRegistrationAuthority) RevokeCertificate(cert x509.Certificate) error {
	return nil
}

func (ra *MockRegistrationAuthority) OnValidationUpdate(authz core.Authorization) error {
	return nil
}

type MockCA struct {}

func (ca *MockCA) IssueCertificate(csr x509.CertificateRequest, regID int64) (cert core.Certificate, err error) {
	// Just a random cert
	cert.DER = RandomCert
	return
}

func (ca *MockCA) GenerateOCSP(xferObj core.OCSPSigningRequest) (ocsp []byte, err error) {
	return
}

func (ca *MockCA) RevokeCertificate(serial string, reasonCode int) (err error) {
	return
}

func makeBody(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

func signRequest(t *testing.T, req string) string {
	accountKeyJSON := []byte(`{"kty":"RSA","n":"z2NsNdHeqAiGdPP8KuxfQXat_uatOK9y12SyGpfKw1sfkizBIsNxERjNDke6Wp9MugN9srN3sr2TDkmQ-gK8lfWo0v1uG_QgzJb1vBdf_hH7aejgETRGLNJZOdaKDsyFnWq1WGJq36zsHcd0qhggTk6zVwqczSxdiWIAZzEakIUZ13KxXvoepYLY0Q-rEEQiuX71e4hvhfeJ4l7m_B-awn22UUVvo3kCqmaRlZT-36vmQhDGoBsoUo1KBEU44jfeK5PbNRk7vDJuH0B7qinr_jczHcvyD-2TtPzKaCioMtNh_VZbPNDaG67sYkQlC15-Ff3HPzKKJW2XvkVG91qMvQ","e":"AAEAAQ","d":"BhAmDbzBAbCeHbU0Xhzi_Ar4M0eTMOEQPnPXMSfW6bc0SRW938JO_-z1scEvFY8qsxV_C0Zr7XHVZsmHz4dc9BVmhiSan36XpuOS85jLWaY073e7dUVN9-l-ak53Ys9f6KZB_v-BmGB51rUKGB70ctWiMJ1C0EzHv0h6Moog-LCd_zo03uuZD5F5wtnPrAB3SEM3vRKeZHzm5eiGxNUsaCEzGDApMYgt6YkQuUlkJwD8Ky2CkAE6lLQSPwddAfPDhsCug-12SkSIKw1EepSHz86ZVfJEnvY-h9jHIdI57mR1v7NTCDcWqy6c6qIzxwh8n2X94QTbtWT3vGQ6HXM5AQ","p":"2uhvZwNS5i-PzeI9vGx89XbdsVmeNjVxjH08V3aRBVY0dzUzwVDYk3z7sqBIj6de53Lx6W1hjmhPIqAwqQgjIKH5Z3uUCinGguKkfGDL3KgLCzYL2UIvZMvTzr9NWLc0AHMZdee5utxWKCGnZBOqy1Rd4V-6QrqjEDBvanoqA60","q":"8odNkMEiriaDKmvwDv-vOOu3LaWbu03yB7VhABu-hK5Xx74bHcvDP2HuCwDGGJY2H-xKdMdUPs0HPwbfHMUicD2vIEUDj6uyrMMZHtbcZ3moh3-WESg3TaEaJ6vhwcWXWG7Wc46G-HbCChkuVenFYYkoi68BAAjloqEUl1JBT1E"}`)
	var accountKey jose.JsonWebKey
	err := json.Unmarshal(accountKeyJSON, &accountKey)
	test.AssertNotError(t, err, "Failed to unmarshal key")
	signer, err := jose.NewSigner("RS256", &accountKey)
	test.AssertNotError(t, err, "Failed to make signer")
	result, err := signer.Sign([]byte(req))
	test.AssertNotError(t, err, "Failed to sign req")
	ret := result.FullSerialize()
	return ret
}

func TestIndex(t *testing.T) {
	wfe := NewWebFrontEndImpl()
	// panic: http: multiple registrations for / [recovered]
	//wfe.HandlePaths()
	wfe.NewReg = "/acme/new-reg"
	responseWriter := httptest.NewRecorder()

	url, _ := url.Parse("/")
	wfe.Index(responseWriter, &http.Request{
		URL: url,
	})
	test.AssertEquals(t, responseWriter.Code, http.StatusOK)
	test.AssertNotEquals(t, responseWriter.Body.String(), "404 page not found\n")
	test.Assert(t, strings.Contains(responseWriter.Body.String(), wfe.NewReg),
		"new-reg not found")

	responseWriter.Body.Reset()
	url, _ = url.Parse("/foo")
	wfe.Index(responseWriter, &http.Request{
		URL: url,
	})
	//test.AssertEquals(t, responseWriter.Code, http.StatusNotFound)
	test.AssertEquals(t, responseWriter.Body.String(), "404 page not found\n")
}

// TODO: Write additional test cases for:
//  - RA returns with a cert success
//  - RA returns with a failure
func TestIssueCertificate(t *testing.T) {
	// TODO: Use a mock RA so we can test various conditions of authorized, not authorized, etc.
	ra := ra.NewRegistrationAuthorityImpl()
	ra.SA = &MockSA{}
	ra.CA = &MockCA{}
	wfe := NewWebFrontEndImpl()
	wfe.SA = &MockSA{}
	wfe.RA = &ra
	wfe.Stats, _ = statsd.NewNoopClient()
	responseWriter := httptest.NewRecorder()

	// GET instead of POST should be rejected
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "GET",
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Method not allowed\"}")

	// POST, but no body.
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, but body that isn't valid JWS
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody("hi"),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, Properly JWS-signed, but payload is "foo", not base64-encoded JSON.
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
{
    "header": {
        "alg": "RS256",
        "jwk": {
            "e": "AQAB",
            "kty": "RSA",
            "n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
        }
    },
    "payload": "Zm9vCg",
    "signature": "hRt2eYqBd_MyMRNIh8PEIACoFtmBi7BHTLBaAhpSU6zyDAFdEBaX7us4VB9Vo1afOL03Q8iuoRA0AT4akdV_mQTAQ_jhTcVOAeXPr0tB8b8Q11UPQ0tXJYmU4spAW2SapJIvO50ntUaqU05kZd0qw8-noH1Lja-aNnU-tQII4iYVvlTiRJ5g8_CADsvJqOk6FcHuo2mG643TRnhkAxUtazvHyIHeXMxydMMSrpwUwzMtln4ZJYBNx4QGEq6OhpAD_VSp-w8Lq5HOwGQoNs0bPxH1SGrArt67LFQBfjlVr94E1sn26p4vigXm83nJdNhWAMHHE9iV67xN-r29LT-FjA"
}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error unmarshaling certificate request\"}")

	// Same signed body, but payload modified by one byte, breaking signature.
	// should fail JWS verification.
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "vd7rZIoTLEe-z1_8G1FcXSw9CQFEJgV4g9V277sER7yx5Qjz_Pkf2YVth6wwwFJEmzc0hoKY-MMYFNwBE4hQHw"
							}
					},
					"payload": "xm9vCg",
					"signature": "RjUQ679fxJgeAJlxqgvDP_sfGZnJ-1RgWF2qmcbnBWljs6h1qp63pLnJOl13u81bP_bCSjaWkelGG8Ymx_X-aQ"
			}
    `),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// Valid, signed JWS body, payload is '{}'
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
							}
					},
					"payload": "e30K",
					"signature": "JXYA_pin91Bc5oz5I6dqCNNWDrBaYTB31EnWorrj4JEFRaidafC9mpLDLLA9jR9kX_Vy2bA5b6pPpXVKm0w146a0L551OdL8JrrLka9q6LypQdDLLQa76XD03hSBOFcC-Oo5FLPa3WRWS1fQ37hYAoLxtS3isWXMIq_4Onx5bq8bwKyu-3E3fRb_lzIZ8hTIWwcblCTOfufUe6AoK4m6MfBjz0NGhyyk4lEZZw6Sttm2VuZo3xmWoRTJEyJG5AOJ6fkNJ9iQQ1kVhMr0ZZ7NVCaOZAnxrwv2sCjY6R3f4HuEVe1yzT75Mq2IuXq-tadGyFujvUxF6BWHCulbEnss7g"
			}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error unmarshaling certificate request\"}")

	// Valid, signed JWS body, payload has a invalid signature on CSR and no authorizations:
	// {
	//   "csr": "MIICU...",
	//   "authorizations: []
	// }
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
							}
					},
					"payload": "ICAgIHsKICAgICAgImNzciI6ICJNSUlDVXpDQ0FUc0NBUUF3RGpFTU1Bb0dBMVVFQXd3RFptOXZNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTNVV2NlMlBZOXk4bjRCN2pPazNEWFpudTJwVWdMcXM3YTVEelJCeG5QcUw3YXhpczZ0aGpTQkkyRk83dzVDVWpPLW04WGpELUdZV2dmWGViWjNhUVZsQmlZcWR4WjNVRzZSRHdFYkJDZUtvN3Y4Vy1VVWZFU05OQ1hGODc0ZGRoSm1FdzBSRjBZV1NBRWN0QVlIRUdvUEZ6NjlnQ3FsNnhYRFBZMU9scE1BcmtJSWxxOUVaV3dUMDgxZWt5SnYwR1lSZlFpZ0NNSzRiMWdrRnZLc0hqYTktUTV1MWIwQVp5QS1tUFR1Nno1RVdrQjJvbmhBWHdXWFg5MHNmVWU4RFNldDlyOUd4TWxuM2xnWldUMXpoM1JNWklMcDBVaGgzTmJYbkE4SkludWtoYTNIUE84V2dtRGQ0SzZ1QnpXc28wQTZmcDVOcFgyOFpwS0F3TTVpUWx0UUlEQVFBQm9BQXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBRkdKVjNPY2doSkVadk9faEd0SWRhUm5zdTZlWDNDZXFTMGJZY0VFemE4dml6bGo0eDA5bnRNSDNRb29xUE9qOHN1dWwwdkQ3NUhaVHB6NkZIRTdTeUxlTktRQkdOR3AxUE1XbVhzRnFENnhVUkN5TUh2Q1pvSHlucENyN0Q1SHR6SXZ1OWZBVjdYUks3cUJLWGZSeGJ2MjFxMHlzTVduZndrYlMyd3JzMXdBelBQZzRpR0pxOHVWSXRybGNGTDhidUpMenh2S2EzbHVfT2p4TlhqemRFdDNWVmtvLUFLUzFzd2tZRWhzR3dLZDhaek5icEYySVEtb2tYZ1JfWmVjeVc4dDgzcFYtdzMzR2hETDl3NlJMUk1nU001YW9qeThyaTdZSW9JdmMzLTlrbGJ3Mmt3WTVvTTJsbWhvSU9HVTEwVGtFeW4xOG15eV81R1VFR2hOelBBPSIsCiAgICAgICJhdXRob3JpemF0aW9ucyI6IFtdCiAgICB9Cg",
					"signature": "PxtFtDXR74ZDgZUWsNaMFpFAhJrYtCYpl3-vr9SCwuWIxB9hZCnLWB5JFwNuC9CtTSYXqDJhzPs4-Bzh345HdwO-ifu1EIVxmc3bAszYS-cxA0lDzr8wJ0ldX0WvADshRWaeFYWJja7ggW03k5JZiNa9AigKIvkGBS2YWpEpCo954cdCEmIL3UOdVjN9aXRT7zzC9wczv4-hYDR-6uP_8J6ATUXJ-UJaTnMi3R0cwtHIcTBZgtgGspoCbtgv-3KaAGNkm5AY062xO5_GbefWwuD2hd8AjKyoTLdfQtwadu6Q3Zl6ZzW_eAfQVDnoblgSt19Gtm4HP4Rf_GosGjRMog"
			}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		// TODO: I think this is wrong. The CSR in the payload above was created by openssl and should be valid. (But the signature isn't)
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error creating new cert\"}")

	// Valid, signed JWS body, payload has a CSR with no DNS names
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
				"payload":"eyJhdXRob3JpemF0aW9ucyI6W10sImNzciI6Ik1JSUJCVENCc2dJQkFEQk5NUW93Q0FZRFZRUUdFd0ZqTVFvd0NBWURWUVFLRXdGdk1Rc3dDUVlEVlFRTEV3SnZkVEVLTUFnR0ExVUVCeE1CYkRFS01BZ0dBMVVFQ0JNQmN6RU9NQXdHQTFVRUF4TUZUMmdnYUdrd1hEQU5CZ2txaGtpRzl3MEJBUUVGQUFOTEFEQklBa0VBc3I3NlprVTJSVHFpNDFlSGZtcEU1aHREdmtyMjAyeWpSUzh4Mk01eXpUNTJvb1QyV0VWdG5TdWltMFlmT0V3NmYtZkhtYnFzYXNxS21xbHNKZGd6MlFJREFRQUJvQUF3Q3dZSktvWklodmNOQVFFRkEwRUFIa0N2NGtWUEphNTNsdE9HcmhwZEgwbVQwNHFIVXFpVGxsSlBQanhYeG42aXdpVllMOG5RdWhzNFEyNzU4RU5vT0RCdU0yRjhnSDE5VElvWGxjbTNMUT09In0",
				"protected":"eyJhbGciOiJSUzI1NiIsImp3ayI6eyJrdHkiOiJSU0EiLCJuIjoieU5XVmh0WUVLSlIyMXk5eHNIVi1QRF9iWXdiWFNlTnVGYWw0NnhZeFZmUkw1bXFoYTd2dHR2akJfdmM3WGcyUnZnQ3hIUENxb3hnTVBUekhyWlQ3NUxqQ3dJVzJLX2tsQllOOG9ZdlR3d21lU2tBejZ1dDdaeFB2LW5aYVQ1VEpoR2swTlQya2hfelNwZHJpRUpfM3ZXLW1xeFliYkJtcHZIcXNhMV96eDlmU3VIWWN0QVpKV3p4elVaWHlrYldNV1FacEVpRTBKNGFqajUxZkluRXpWbjdWeFYtbXpmTXlib1FqdWpQaDdhTkp4QVdTcTRvUUVKSkRnV3dTaDlsZXlvSm9QcE9OSHhoNW5FRTVBakUwMUZrR0lDU3hqcFpzRi13OGhPVEkzWFhvaFVkdTI5U2UyNmsyQjBQb2xEU3VqMEdJUVU2LVc5VGRMWFNqQmIyU3BRIiwiZSI6IkFBRUFBUSJ9fQ",
				"signature":"LslpZp6wLYQo0LAgMl9_jyTFhKVnvFWcD455-v2b3q3wXJX5Ksvp4sxyczM63j2RGwTUc_Tfu3WEWa2xQ-D74H69XGMnWCikmChwVPDcWwaDwydOEFXff5cGY4Trkxl7xnsO2g3BslxuZ_7uud5IkHIy1-8xa4mHpNHb3XHTAhX5E3tXA1VqC4pVWzD5W74bg4GxuRd8IM2p3toMjgInbzp9vhY7dnPYogwnA8B1uYduF99azdKkb5VbHNBJi5SpTz7nyjvbvh7KTLhaJ1epkSFnd74a-fhyzo8t1Nju9UPT1nc8kF6G3CpOAyWYX27YyA9T0UyM3CVz_hFpvubZjg"
			}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error creating new cert\"}")

	// Valid, signed JWS body, payload has a valid CSR but no authorizations:
	// {
	//   "csr": "MIIBK...",
	//   "authorizations: []
	// }
	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
				"payload":"eyJhdXRob3JpemF0aW9ucyI6W10sImNzciI6Ik1JSUJLekNCMkFJQkFEQk5NUW93Q0FZRFZRUUdFd0ZqTVFvd0NBWURWUVFLRXdGdk1Rc3dDUVlEVlFRTEV3SnZkVEVLTUFnR0ExVUVCeE1CYkRFS01BZ0dBMVVFQ0JNQmN6RU9NQXdHQTFVRUF4TUZUMmdnYUdrd1hEQU5CZ2txaGtpRzl3MEJBUUVGQUFOTEFEQklBa0VBcXZGRUdCTnJqQW90UGJjZFRTeURweHNFU04wLWVZbDRUcVMwWkxZd0xUVi1GdVBIVFBqRmlxMm9IMUJFZ21SempiOFlpUFZYRk1uYU9lSEU3enV1WFFJREFRQUJvQ1l3SkFZSktvWklodmNOQVFrT01SY3dGVEFUQmdOVkhSRUVEREFLZ2dodFpXVndMbU52YlRBTEJna3Foa2lHOXcwQkFRVURRUUJTRWNFcS1sTVVuenYxRE84akswaEpSOFlLYzB5Vjh6dVdWZkFXTjBfZHNQZzVOeS1PSGh0SmNPVElyVXJMVGJfeENVN2NqaUt4VThpM2oxa2FULXJ0In0",
				"protected":"eyJhbGciOiJSUzI1NiIsImp3ayI6eyJrdHkiOiJSU0EiLCJuIjoieU5XVmh0WUVLSlIyMXk5eHNIVi1QRF9iWXdiWFNlTnVGYWw0NnhZeFZmUkw1bXFoYTd2dHR2akJfdmM3WGcyUnZnQ3hIUENxb3hnTVBUekhyWlQ3NUxqQ3dJVzJLX2tsQllOOG9ZdlR3d21lU2tBejZ1dDdaeFB2LW5aYVQ1VEpoR2swTlQya2hfelNwZHJpRUpfM3ZXLW1xeFliYkJtcHZIcXNhMV96eDlmU3VIWWN0QVpKV3p4elVaWHlrYldNV1FacEVpRTBKNGFqajUxZkluRXpWbjdWeFYtbXpmTXlib1FqdWpQaDdhTkp4QVdTcTRvUUVKSkRnV3dTaDlsZXlvSm9QcE9OSHhoNW5FRTVBakUwMUZrR0lDU3hqcFpzRi13OGhPVEkzWFhvaFVkdTI5U2UyNmsyQjBQb2xEU3VqMEdJUVU2LVc5VGRMWFNqQmIyU3BRIiwiZSI6IkFBRUFBUSJ9fQ",
				"signature":"kdu5tXk-Jz9umbi6RH-BACBj5ObJlVPA4qGsLsdbqPfn9W9CDw66Q9E1QQxt9Fxpe-fqdSDiVSfmuXhO7u068xdYptgFxWNDJXM1MH3iCs0EJz5KQ9SfGiJXrhkji_FbOdYwcxSvbThOF_qyztFmBCgZPfKKHbcGJKV3nvFDLHb6P7hIr6iqMutsFykTToYBUv3czzc87iYFpR_ukAnISLJ0hQucbMBqlinvq8TOmzi47o_uv2Fy8MF0V_C9ZYJmhZGjihhVvVlr00OaFE5bNM2uLMfr_02oG83HTjNJOGaxsqi-tuu41m3Dr5M8Ubh2oPA0OrKIuMisMMZ3aCzRtA"
			}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error creating new cert\"}")

	responseWriter.Body.Reset()
	wfe.NewCertificate(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
				"payload":"eyJhdXRob3JpemF0aW9ucyI6WyJ2YWxpZCJdLCJjc3IiOiJNSUlCR3pDQnlBSUJBREE5TVFvd0NBWURWUVFHRXdGak1Rb3dDQVlEVlFRS0V3RnZNUXN3Q1FZRFZRUUxFd0p2ZFRFS01BZ0dBMVVFQnhNQmJERUtNQWdHQTFVRUNCTUJjekJjTUEwR0NTcUdTSWIzRFFFQkFRVUFBMHNBTUVnQ1FRRFpEa3V2TllMZi1sdW1KajZyYzdoaWZiMmpxYVhxZVBYY1FhQVdGOU1YWVhDWkwwYjJRSXdRQm9ZZkp6dzA0RnNPSlFuSlpJXzBQV2MtTFMyb2FKaFRBZ01CQUFHZ0pqQWtCZ2txaGtpRzl3MEJDUTR4RnpBVk1CTUdBMVVkRVFRTU1BcUNDRzFsWlhBdVkyOXRNQXNHQ1NxR1NJYjNEUUVCQlFOQkFLVTkwcHpvMU1TbFFYZnNLWlVQTXJDcXFHWHpXWkxmY29MaHZ0Y3NfWnpDQkpscm9xYzJoUGpMQzBISXBBODNZdEFwaTlsbjIyRmlxdXdNNWVwUHZfST0ifQ",
				"protected":"eyJhbGciOiJSUzI1NiIsImp3ayI6eyJrdHkiOiJSU0EiLCJuIjoieU5XVmh0WUVLSlIyMXk5eHNIVi1QRF9iWXdiWFNlTnVGYWw0NnhZeFZmUkw1bXFoYTd2dHR2akJfdmM3WGcyUnZnQ3hIUENxb3hnTVBUekhyWlQ3NUxqQ3dJVzJLX2tsQllOOG9ZdlR3d21lU2tBejZ1dDdaeFB2LW5aYVQ1VEpoR2swTlQya2hfelNwZHJpRUpfM3ZXLW1xeFliYkJtcHZIcXNhMV96eDlmU3VIWWN0QVpKV3p4elVaWHlrYldNV1FacEVpRTBKNGFqajUxZkluRXpWbjdWeFYtbXpmTXlib1FqdWpQaDdhTkp4QVdTcTRvUUVKSkRnV3dTaDlsZXlvSm9QcE9OSHhoNW5FRTVBakUwMUZrR0lDU3hqcFpzRi13OGhPVEkzWFhvaFVkdTI5U2UyNmsyQjBQb2xEU3VqMEdJUVU2LVc5VGRMWFNqQmIyU3BRIiwiZSI6IkFBRUFBUSJ9fQ",
				"signature":"n7JIc5QeR0PK1JyIOIdJ_KWwNwtz2Ke89FGbnAb3LRbpSbkMnZNx1TNeg5TNxSLGv9FSoCcmNvMlmK5DR1sRJddsnryI1F36qZl6lvv-ZrVjKlJ6__HiQCXdVsXGo5hMNLc834cyp1RwadmjQr46cv5vv0cOSuPfSsjdNC9Aqw7AqVtqOT6ONmZQdwQFKVxK8I34wx67-pc4Rwia5Bbo6yZFebgKQkDc3EzJsqeUHY5aTbCK6aD9wKx4PA5M3bwVy5xTPNsUAsSYHdxvnTQ9WpzY20f_7ha8Zdw8Xf2y6pZ1F8mjGHAifqTzEMZpHTjW8-6syiq9nSWUM4KiIODkdA"
			}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		string(RandomCert))
	test.AssertEquals(
		t, responseWriter.Header().Get("Location"),
		"/acme/authz/asdf?challenge=foo")
}

func TestChallenge(t *testing.T) {
	wfe := NewWebFrontEndImpl()
	wfe.RA = &MockRegistrationAuthority{}
	wfe.SA = &MockSA{}
	wfe.HandlePaths()
	responseWriter := httptest.NewRecorder()

	var key jose.JsonWebKey
	err := json.Unmarshal([]byte(`{
						"e": "AQAB",
						"kty": "RSA",
						"n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
							}`), &key)
	test.AssertNotError(t, err, "Could not unmarshal testing key")

	challengeURL, _ := url.Parse("/acme/authz/asdf?challenge=foo")
	authz := core.Authorization{
		ID: "asdf",
		Identifier: core.AcmeIdentifier{
			Type:  "dns",
			Value: "letsencrypt.org",
		},
		Challenges: []core.Challenge{
			core.Challenge{
				Type: "dns",
				URI:  core.AcmeURL(*challengeURL),
			},
		},
		RegistrationID: 1,
	}

	wfe.Challenge(authz, responseWriter, &http.Request{
		Method: "POST",
		URL:    challengeURL,
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
							}
					},
					"payload": "e30K",
					"signature": "JXYA_pin91Bc5oz5I6dqCNNWDrBaYTB31EnWorrj4JEFRaidafC9mpLDLLA9jR9kX_Vy2bA5b6pPpXVKm0w146a0L551OdL8JrrLka9q6LypQdDLLQa76XD03hSBOFcC-Oo5FLPa3WRWS1fQ37hYAoLxtS3isWXMIq_4Onx5bq8bwKyu-3E3fRb_lzIZ8hTIWwcblCTOfufUe6AoK4m6MfBjz0NGhyyk4lEZZw6Sttm2VuZo3xmWoRTJEyJG5AOJ6fkNJ9iQQ1kVhMr0ZZ7NVCaOZAnxrwv2sCjY6R3f4HuEVe1yzT75Mq2IuXq-tadGyFujvUxF6BWHCulbEnss7g"
			}
		`),
	})

	test.AssertEquals(
		t, responseWriter.Header().Get("Location"),
		"/acme/authz/asdf?challenge=foo")
	test.AssertEquals(
		t, responseWriter.Header().Get("Link"),
		"</acme/authz/asdf>;rel=\"up\"")
	test.AssertEquals(
		t, responseWriter.Body.String(),
		"{\"type\":\"dns\",\"uri\":\"/acme/authz/asdf?challenge=foo\"}")
}

func TestNewRegistration(t *testing.T) {
	wfe := NewWebFrontEndImpl()
	wfe.RA = &MockRegistrationAuthority{}
	wfe.SA = &MockSA{}
	wfe.Stats, _ = statsd.NewNoopClient()
	wfe.SubscriberAgreementURL = "https://letsencrypt.org/be-good"
	responseWriter := httptest.NewRecorder()

	// GET instead of POST should be rejected
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "GET",
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Method not allowed\"}")

	// POST, but no body.
	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, but body that isn't valid JWS
	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody("hi"),
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, Properly JWS-signed, but payload is "foo", not base64-encoded JSON.
	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
{
    "header": {
        "alg": "RS256",
        "jwk": {
            "e": "AQAB",
            "kty": "RSA",
            "n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
        }
    },
    "payload": "Zm9vCg",
    "signature": "hRt2eYqBd_MyMRNIh8PEIACoFtmBi7BHTLBaAhpSU6zyDAFdEBaX7us4VB9Vo1afOL03Q8iuoRA0AT4akdV_mQTAQ_jhTcVOAeXPr0tB8b8Q11UPQ0tXJYmU4spAW2SapJIvO50ntUaqU05kZd0qw8-noH1Lja-aNnU-tQII4iYVvlTiRJ5g8_CADsvJqOk6FcHuo2mG643TRnhkAxUtazvHyIHeXMxydMMSrpwUwzMtln4ZJYBNx4QGEq6OhpAD_VSp-w8Lq5HOwGQoNs0bPxH1SGrArt67LFQBfjlVr94E1sn26p4vigXm83nJdNhWAMHHE9iV67xN-r29LT-FjA"
}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error unmarshaling JSON\"}")

	// Same signed body, but payload modified by one byte, breaking signature.
	// should fail JWS verification.
	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "vd7rZIoTLEe-z1_8G1FcXSw9CQFEJgV4g9V277sER7yx5Qjz_Pkf2YVth6wwwFJEmzc0hoKY-MMYFNwBE4hQHw"
							}
					},
					"payload": "xm9vCg",
					"signature": "RjUQ679fxJgeAJlxqgvDP_sfGZnJ-1RgWF2qmcbnBWljs6h1qp63pLnJOl13u81bP_bCSjaWkelGG8Ymx_X-aQ"
			}
    	`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody(signRequest(t, "{\"contact\":[\"tel:123456789\"],\"agreement\":\"https://letsencrypt.org/im-bad\"}")),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Provided agreement URL [https://letsencrypt.org/im-bad] does not match current agreement URL [https://letsencrypt.org/be-good]\"}")

	responseWriter.Body.Reset()
	wfe.NewRegistration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody(signRequest(t, "{\"contact\":[\"tel:123456789\"],\"agreement\":\"https://letsencrypt.org/be-good\"}")),
	})

	test.AssertEquals(t, responseWriter.Body.String(), "{\"id\":0,\"key\":{\"kty\":\"RSA\",\"n\":\"z2NsNdHeqAiGdPP8KuxfQXat_uatOK9y12SyGpfKw1sfkizBIsNxERjNDke6Wp9MugN9srN3sr2TDkmQ-gK8lfWo0v1uG_QgzJb1vBdf_hH7aejgETRGLNJZOdaKDsyFnWq1WGJq36zsHcd0qhggTk6zVwqczSxdiWIAZzEakIUZ13KxXvoepYLY0Q-rEEQiuX71e4hvhfeJ4l7m_B-awn22UUVvo3kCqmaRlZT-36vmQhDGoBsoUo1KBEU44jfeK5PbNRk7vDJuH0B7qinr_jczHcvyD-2TtPzKaCioMtNh_VZbPNDaG67sYkQlC15-Ff3HPzKKJW2XvkVG91qMvQ\",\"e\":\"AAEAAQ\"},\"recoveryToken\":\"\",\"contact\":[\"tel:123456789\"],\"agreement\":\"https://letsencrypt.org/be-good\",\"thumbprint\":\"\"}")
	var reg core.Registration
	err := json.Unmarshal([]byte(responseWriter.Body.String()), &reg)
	test.AssertNotError(t, err, "Couldn't unmarshal returned registration object")
	uu := url.URL(reg.Contact[0])
	test.AssertEquals(t, uu.String(), "tel:123456789")
}

func TestAuthorization(t *testing.T) {
	wfe := NewWebFrontEndImpl()
	wfe.RA = &MockRegistrationAuthority{}
	wfe.SA = &MockSA{}
	wfe.Stats, _ = statsd.NewNoopClient()
	responseWriter := httptest.NewRecorder()

	// GET instead of POST should be rejected
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "GET",
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Method not allowed\"}")

	// POST, but no body.
	responseWriter.Body.Reset()
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "POST",
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, but body that isn't valid JWS
	responseWriter.Body.Reset()
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody("hi"),
	})
	test.AssertEquals(t, responseWriter.Body.String(), "{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	// POST, Properly JWS-signed, but payload is "foo", not base64-encoded JSON.
	responseWriter.Body.Reset()
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
{
    "header": {
        "alg": "RS256",
        "jwk": {
            "e": "AQAB",
            "kty": "RSA",
            "n": "tSwgy3ORGvc7YJI9B2qqkelZRUC6F1S5NwXFvM4w5-M0TsxbFsH5UH6adigV0jzsDJ5imAechcSoOhAh9POceCbPN1sTNwLpNbOLiQQ7RD5mY_pSUHWXNmS9R4NZ3t2fQAzPeW7jOfF0LKuJRGkekx6tXP1uSnNibgpJULNc4208dgBaCHo3mvaE2HV2GmVl1yxwWX5QZZkGQGjNDZYnjFfa2DKVvFs0QbAk21ROm594kAxlRlMMrvqlf24Eq4ERO0ptzpZgm_3j_e4hGRD39gJS7kAzK-j2cacFQ5Qi2Y6wZI2p-FCq_wiYsfEAIkATPBiLKl_6d_Jfcvs_impcXQ"
        }
    },
    "payload": "Zm9vCg",
    "signature": "hRt2eYqBd_MyMRNIh8PEIACoFtmBi7BHTLBaAhpSU6zyDAFdEBaX7us4VB9Vo1afOL03Q8iuoRA0AT4akdV_mQTAQ_jhTcVOAeXPr0tB8b8Q11UPQ0tXJYmU4spAW2SapJIvO50ntUaqU05kZd0qw8-noH1Lja-aNnU-tQII4iYVvlTiRJ5g8_CADsvJqOk6FcHuo2mG643TRnhkAxUtazvHyIHeXMxydMMSrpwUwzMtln4ZJYBNx4QGEq6OhpAD_VSp-w8Lq5HOwGQoNs0bPxH1SGrArt67LFQBfjlVr94E1sn26p4vigXm83nJdNhWAMHHE9iV67xN-r29LT-FjA"
}
		`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Error unmarshaling JSON\"}")

	// Same signed body, but payload modified by one byte, breaking signature.
	// should fail JWS verification.
	responseWriter.Body.Reset()
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`
			{
					"header": {
							"alg": "RS256",
							"jwk": {
									"e": "AQAB",
									"kty": "RSA",
									"n": "vd7rZIoTLEe-z1_8G1FcXSw9CQFEJgV4g9V277sER7yx5Qjz_Pkf2YVth6wwwFJEmzc0hoKY-MMYFNwBE4hQHw"
							}
					},
					"payload": "xm9vCg",
					"signature": "RjUQ679fxJgeAJlxqgvDP_sfGZnJ-1RgWF2qmcbnBWljs6h1qp63pLnJOl13u81bP_bCSjaWkelGG8Ymx_X-aQ"
			}
    	`),
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")

	responseWriter.Body.Reset()
	wfe.NewAuthorization(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody(signRequest(t, "{\"identifier\":{\"type\":\"dns\",\"value\":\"test.com\"}}")),
	})

	test.AssertEquals(t, responseWriter.Body.String(), "{\"identifier\":{\"type\":\"dns\",\"value\":\"test.com\"},\"expires\":\"0001-01-01T00:00:00Z\"}")

	var authz core.Authorization
	err := json.Unmarshal([]byte(responseWriter.Body.String()), &authz)
	test.AssertNotError(t, err, "Couldn't unmarshal returned authorization object")
}

func TestRegistration(t *testing.T) {
	wfe := NewWebFrontEndImpl()
	wfe.RA = &MockRegistrationAuthority{}
	wfe.SA = &MockSA{}
	wfe.Stats, _ = statsd.NewNoopClient()
	wfe.SubscriberAgreementURL = "https://letsencrypt.org/be-good"
	responseWriter := httptest.NewRecorder()

	// Test invalid method
	path, _ := url.Parse("/1")
	wfe.Registration(responseWriter, &http.Request{
		Method: "MAKE-COFFEE",
		Body:   makeBody("invalid"),
		URL:    path,
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Method not allowed\"}")
	responseWriter.Body.Reset()

	// Test GET proper entry returns 405
	path, _ = url.Parse("/1")
	wfe.Registration(responseWriter, &http.Request{
		Method: "GET",
		URL:    path,
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Method not allowed\"}")
	responseWriter.Body.Reset()

	// Test POST invalid JSON
	path, _ = url.Parse("/2")
	wfe.Registration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody("invalid"),
		URL:    path,
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Unable to read/verify body\"}")
	responseWriter.Body.Reset()

	// Test POST valid JSON but key is not registered
	path, _ = url.Parse("/2")
	wfe.Registration(responseWriter, &http.Request{
		Method: "POST",
		Body: makeBody(`{
		   "payload" : "ewogICJjb250YWN0IjogWwogICAgIm1haWx0bzpjZXJ0LWFkbWluQGV4YW1wbGUuY28ubnoiLAogICAgInRlbDorMjQ5NTU1MTIxMiIKICBdLAogICJhZ3JlZW1lbnQiOiAieWVzIgp9Cg",
		   "protected" : "eyJhbGciOiJQUzI1NiIsImp3ayI6eyJrdHkiOiJSU0EiLCJuIjoibTVDcHgzdlowQ2pBVGlyRHBiSUx2cTc4Zm0zRHY1UkJrTzFWTFdGbUpqNU1iNTR2YzlvWVpXYzFWMWstTEpvRVN1dVBIaGFOTzJFdThUOXRzbFFXY1pTenI1TklteEF3TWs5NzBnVlFhLUhxdi1KcjZ4c3RyQnBxN1RLcFhIVHgyRm5mQTJ3UXJmSVFTbEJYdTB0NGpkVU9yM29KaC1RWHZtYThuTElUZHRqcEMwQVpOdHFkMFFrUkpYXzkwU2FOcmwxOFJyXzBKckJIOVptVVNGY2YzbW9fQnRMMEd4MGpFM24taXdDSThyUXRmeVZQX185LW5fX3I0SWhhbEtMemFlaW82by1xcmRlbWgwRVpnaktHQ1MxX1JwVEllQXJrTzh1aWExS2dPcS16LUdmZW1LRW00czA3V09fYTBfOWRMcWJ2cG55eVp2VWk0MDVtM3ZHRGZRIiwiZSI6IkFBRUFBUSJ9fQ",
		   "signature" : "exg0HJRHk-oSDiaOlgtTkT_COqDRyIAJr4g9fDAJh5GF5evXAfT0Hbkfy4TYzqvF6oOldIaCylYhXjYtve4JLXEMdAj1DaR7kGVALskLg-XbiZ0-IaFBiDDaT6mwyLBTfstX4DD2OL7x0vyuTK16bHEIF0hncwHYVSoX5eFOBQLVu_gjxc7J5OZK4ugSJxZEilTVta0A9EdXdUxth0qqbZg_hJDmGOyNge03C71GbhMs-DF-rujlhe7L4VhcV3U0Wj8kSuAGn_DIHBJ1zM0H46PRgyz_9DgkJ6XnE5W8ZA3kF0VPFSp4ofqBhkFUXLXPPJJUEurAQxBJMaU31ef8bg"
		}`),
		URL: path,
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:unauthorized\",\"detail\":\"No registration exists matching provided key\"}")
	responseWriter.Body.Reset()

	key, err := jose.LoadPrivateKey([]byte(test1KeyPrivatePEM))
	test.AssertNotError(t, err, "Failed to load key")
	rsaKey, ok := key.(*rsa.PrivateKey)
	test.Assert(t, ok, "Couldn't load RSA key")
	signer, err := jose.NewSigner("RS256", rsaKey)
	test.AssertNotError(t, err, "Failed to make signer")

	path, _ = url.Parse("/2")

	// Test POST valid JSON with registration up in the mock (with incorrect agreement URL)
	result, err := signer.Sign([]byte("{\"agreement\":\"https://letsencrypt.org/im-bad\"}"))

	// Test POST valid JSON with registration up in the mock
	path, _ = url.Parse("/1")
	wfe.Registration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody(result.FullSerialize()),
		URL:    path,
	})
	test.AssertEquals(t,
		responseWriter.Body.String(),
		"{\"type\":\"urn:acme:error:malformed\",\"detail\":\"Provided agreement URL [https://letsencrypt.org/im-bad] does not match current agreement URL [https://letsencrypt.org/be-good]\"}")
	responseWriter.Body.Reset()

	// Test POST valid JSON with registration up in the mock (with correct agreement URL)
	result, err = signer.Sign([]byte("{\"agreement\":\"https://letsencrypt.org/be-good\"}"))
	wfe.Registration(responseWriter, &http.Request{
		Method: "POST",
		Body:   makeBody(result.FullSerialize()),
		URL:    path,
	})
	test.AssertNotContains(t, responseWriter.Body.String(), "urn:acme:error")
	responseWriter.Body.Reset()
}

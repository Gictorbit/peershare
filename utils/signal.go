package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"github.com/pion/randutil"
	"io"
)

// Allows compressing offer/answer to bypass terminal input limits.
const (
	compress  = true
	runeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

func zip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		return nil, err
	}
	err = gz.Flush()
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return nil, err
	}
	res, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RandSeq generates a random string to serve as dummy data
//
// It returns a deterministic sequence of values each time a program is run.
// Use rand.Seed() function in your real applications.
func RandSeq(n int) (string, error) {
	val, err := randutil.GenerateCryptoRandomString(n, runeChars)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(b []byte) (string, error) {
	if compress {
		var err error
		b, err = zip(b)
		if err != nil {
			return "", err
		}
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}
	if compress {
		b, err = unzip(b)
		if err != nil {
			return "", err
		}
	}
	return string(b), nil
}

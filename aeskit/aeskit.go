package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/xiaojiaoyu100/lizard/stringkit"
)

/*
	The test case verifies that the average time of
	encrypting or decrypting 1 million data is 1.43s
*/

// AesEncrypt encrypt text by aes
func AesEncrypt(key []byte, content string) (string, error) {
	encryptMap, err := AesEncryptMap(key, []string{content})
	if err != nil {
		return "", err
	}
	return encryptMap[content], nil
}

// AesEncryptMap decrypt text in batches through AES
func AesEncryptMap(key []byte, contents []string) (map[string]string, error) {
	result := make(map[string]string)
	contents = stringkit.UniqueAndFilterStr(contents)
	if len(contents) == 0 {
		return result, nil
	}

	for i := range contents {
		content := contents[i]
		encodeContent := base64.StdEncoding.EncodeToString([]byte(content))
		b, err := base64.StdEncoding.DecodeString(encodeContent)
		if err != nil {
			return nil, fmt.Errorf("[DecodeString] %w", err)
		}
		ciphertext, err := aesEncrypt(b, key)
		if err != nil {
			return nil, fmt.Errorf("[aesEncrypt] %w", err)
		}
		result[content] = base64.StdEncoding.EncodeToString(ciphertext)
	}

	return result, nil
}

// AesDecrypt decrypt text by AES
func AesDecrypt(key []byte, ciphertext string) (string, error) {
	decryptMap, err := AesDecryptMap(key, []string{ciphertext})
	if err != nil {
		return "", err
	}
	return decryptMap[ciphertext], nil
}

// AesDecryptMap decrypt text in batches through AES
func AesDecryptMap(key []byte, ciphertexts []string) (map[string]string, error) {
	result := make(map[string]string)
	ciphertexts = stringkit.UniqueAndFilterStr(ciphertexts)
	if len(ciphertexts) == 0 {
		return result, nil
	}

	for i := range ciphertexts {

		ciphertext := ciphertexts[i]
		b, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return nil, fmt.Errorf("[DecodeString] %w", err)
		}
		content, err := aesDecrypt(b, key)
		if err != nil {
			return nil, fmt.Errorf("[aesDecrypt] %w", err)
		}
		result[ciphertext] = string(content[:])

	}

	return result, nil
}

// aesEncrypt ECB PKCS5 encrypt
func aesEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ecb := newECBEncrypter(block)
	content := src
	content = pkcs5Padding(content, block.BlockSize())
	des := make([]byte, len(content))
	ecb.CryptBlocks(des, content)
	return des, nil
}

// aesDecrypt ECB PKCS5 decrypt
func aesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := newECBDecrypter(block)
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

// pkcs5Padding ...
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// pkcs5UnPadding ...
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

// newECB ...
func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// newECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// newECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func newECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}
func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

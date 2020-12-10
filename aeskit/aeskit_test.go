package utils

import (
	"fmt"
	"testing"
	"time"
)

func init() {
}

func TestAesDecryptMap(t *testing.T) {
	i := 0
	var content []string
	for i < 1000000 {
		content = append(content, fmt.Sprintf("phone%d", i))
		i++
	}
	m, err := AesEncryptMap([]byte("r4DMt6Exz4GU9pYrwh8Fq52NJceIg49d"), content)
	if err != nil {
		t.Fatal(err)
	}
	var list []string
	for _, v := range m {
		list = append(list, v)
	}

	tt := time.Now()
	result, err := AesDecryptMap([]byte("r4DMt6Exz4GU9pYrwh8Fq52NJceIg49d"), list)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Since(tt))
	t.Log("finish, ", len(result))
}

func TestAesEncryptMap(t *testing.T) {
	i := 0
	var content []string
	for i < 1000000 {
		content = append(content, fmt.Sprintf("phone%d", i))
		i++
	}
	tt := time.Now()
	result, err := AesEncryptMap([]byte("r4DMt6Exz4GU9pYrwh8Fq52NJceIg49d"), content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time.Since(tt))
	t.Log("finish, ", len(result))
}

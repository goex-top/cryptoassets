package main

import (
	"testing"
)

func TestAES(t *testing.T) {
	x := []byte("世界上最邪恶最专制的现代奴隶制国家--朝鲜")
	key := []byte("hgfedcba87654321")
	x1 := AESECBEncrypt(x, key)
	t.Log("密文:", len(x1))
	t.Log("密文:", string(x1))
	x2 := AESECBDecrypt(x1, key)
	t.Log("解密:", string(x2))

}

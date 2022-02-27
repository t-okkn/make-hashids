package main

import (
	"encoding/binary"
	"encoding/hex"
	"crypto/sha256"
	"strings"

	hashids "github.com/speps/go-hashids"

	"unsafe"
)

const (
	// Hashidsに使用する文字列
	ALPHABET string = "abcdefghkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXY0123456789"

	// Hashidsの最低文字列長
	MIN_LENGTH int = 4
)


// <summary>: 文字列をHashidsにします
func EncodeToHashids(input string) string {
	hid, err := getHashidsObject(input)
	if err != nil {
		return ""
	}

	if res, err := hid.EncodeInt64(stringToInt64(input)); err != nil {
		return ""
	} else {
		return res
	}
}

// <summary>: 不可逆なHashidsを取得します
func GetShortHashids(input string) string {
	hid, err := getHashidsObject(input)
	if err != nil {
		return ""
	}

	var num int64 = 0
	for _, i := range stringToInt64(input) {
		num += i
	}

	if res, err := hid.EncodeInt64([]int64{num}); err != nil {
		return ""
	} else {
		return res
	}
}

// <summary>: Hashidsを文字列に戻します
func DecodeToString(input string) string {
	hid, err := hashids.New()
	if err != nil {
		return ""
	}

	arr, err := hid.DecodeInt64WithError(input)
	if err != nil {
		return ""
	}

	return int64ToString(arr)
}

// <summary>: 文字列を数字のスライスに変換します
func stringToInt64(input string) []int64 {
	res := make([]int64, len([]rune(input)))
	count := 0

	for _, val := range input {
		// 各文字をbyteスライスに変換
		chr := make([]byte, 8)
		for i, b := range []byte(string(val)) {
			chr[i] = b
		}

		// byteスライスをUint64に変換
		istr := binary.LittleEndian.Uint64(chr)

		// Uint64をInt64にキャスト
		// UTF-8はUint32内に収まるので、エラーはありえない
		res[count] = int64(istr)
		count += 1
	}

	return res
}

// <summary>: 数字のスライスを文字列に変換します
func int64ToString(input []int64) string {
	var sb strings.Builder
	sb.Grow(len(input))

	for _, val := range input {
		// 各数字をbyteスライスに変換
		chr := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(chr, uint64(val))

		// byteスライスから不要な0を切り詰める
		pos := 0
		for i, v := range chr {
			if v == byte(0) {
				pos = i
				break
			}
		}

		sb.Write(chr[:pos])
	}

	return sb.String()
}

// <summary>: HashIDオブジェクトを取得します
func getHashidsObject(input string) (*hashids.HashID, error) {
	// hash := sha256.Sum256([]byte(str))
	hash := sha256.Sum256(*(*[]byte)(unsafe.Pointer(&input)))

	d := hashids.NewData()
	d.Alphabet = ALPHABET
	d.MinLength = MIN_LENGTH
	d.Salt = hex.EncodeToString(hash[:])

	return hashids.NewWithData(d)
}


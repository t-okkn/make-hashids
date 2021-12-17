package main

import (
	"encoding/binary"
	"encoding/hex"
	"crypto/sha256"
	"strings"

	"github.com/speps/go-hashids"

	"unsafe"
)

const (
	// Hashidsに使用する文字列
	ALPHABET string = "abcdefghkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXY0123456789"

	// Hashidsの最低文字列長
	MIN_LENGTH int = 4
)


// summary => 文字列をHashidsにします
// param::input => Hashidsを作成する値
// return::string => Hashids
/////////////////////////////////////////
func EncodeToHashids(input string) string {
	// hash := sha256.Sum256([]byte(str))
	hash := sha256.Sum256(*(*[]byte)(unsafe.Pointer(&input)))

	d := hashids.NewData()
	d.Alphabet = ALPHABET
	d.MinLength = MIN_LENGTH
	d.Salt = hex.EncodeToString(hash[:])

	hid, err := hashids.NewWithData(d)
	if err != nil {
		return ""
	}

	if res, err := hid.EncodeInt64(stringToInt64(input)); err != nil {
		return ""
	} else {
		return res
	}
}

// summary => Hashidsを文字列に戻します
// param::input => Hashids
// return::string => 文字列
/////////////////////////////////////////
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

// summary => 文字列を数字のスライスに変換します
// param::input => 変換する入力値
// return::[]int64 => inputに対する変換後のスライス
/////////////////////////////////////////
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

// summary => 数字のスライスを文字列に変換します
// param::input => 変換する入力値
// return::string => inputに対する変換後の文字列
/////////////////////////////////////////
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


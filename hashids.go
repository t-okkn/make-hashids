package main

import (
	"strconv"
	"time"
	"encoding/binary"
	"encoding/hex"
	"crypto/sha256"
	"github.com/speps/go-hashids"

	"unsafe"
)

// summary => 文字列を一文字ずつ数値に変換します
// param::input => 変換する入力値
// return::[]uint32 => inputに対する変換後の数値
/////////////////////////////////////////
func Str2Uints(input string) []uint32 {
	res := make([]uint32, len([]rune(input)))
	count := 0

	for _, val := range input {
		chr := []byte(string(val))
		length := len(chr)

		if length < 4 {
			for i := 0; i < 4-length; i++ {
				chr = append(chr, 0)
			}
		}

		res[count] = binary.LittleEndian.Uint32(chr)
		count += 1
	}

	return res
}

// summary => Hashidsを作成します
// param::input => Hashidsを作成する値
// return::string => Hashids
/////////////////////////////////////////
func CreateHashids(input []uint32) string {
	ut := time.Now().Unix()

	d := hashids.NewData()
	d.Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	d.MinLength = 8
	d.Salt = getSaltInt64(&ut)

	num := make([]int64, 1)
	num[0] = *sumUint32(input) + ut

	hid, _ := hashids.NewWithData(d)
	res, _ := hid.EncodeInt64(num)

	return res
}

// summary => [Private] uint32のスライスのすべての値を足し合わせます
// param::input => 足し合わせるスライス
// return::*int64 => [p] 合計値
/////////////////////////////////////////
func sumUint32(input []uint32) *int64 {
	var total int64 = 0

	for _, val := range input {
		total += int64(val)
	}

	return &total
}

// summary => [Private] Int64の入力値に対するSha256を文字列で導出します
// param::input => [p] Int64の入力値
// return::string => Sha256文字列
// remark => unsafe使用
/////////////////////////////////////////
func getSaltInt64(input *int64) string {
	str := strconv.FormatInt(*input, 10)
	// hash := sha256.Sum256([]byte(str))
	hash := sha256.Sum256(*(*[]byte)(unsafe.Pointer(&str)))

	return hex.EncodeToString(hash[:])
}


package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// POST時の要求コンテンツの上限値
	LIMIT_CONTENTS int = 1000

	// 変換元文字列長の最大値の上限値
	LIMIT_MAX int = 4096

	// デフォルトの変換元文字列長の最大値
	DEFAULT_MAX int = 256
)

// summary => Response内容について定義した構造体
// param::Error => エラー内容
// param::Hashids => Hashids
// param::Source => 変換元の文字列
/////////////////////////////////////////
type HashSet struct {
	Error   string `json:"error"`
	Hashids string `json:"hashids"`
	Source  string `json:"source"`
}

// summary => 待ち受けるサーバのルーターを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/hashids", getSingleHashids)
	router.GET("/hashids/max/:max", getSingleHashids)
	router.POST("/hashids", getMultiHashids)
	router.POST("/hashids/max/:max", getMultiHashids)

	return router
}

// summary => 単一の要求に対してレスポンスを返します
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getSingleHashids(c *gin.Context) {
	maxprm := c.Param("max")
	max := getMaxValue(maxprm)

	res := make([]HashSet, 1)
	str := c.DefaultQuery("source", "")

	// 単一要求時はエラーがあれば400とする
	if str == "" {
		res[0].Error = "文字列が読み取れませんでした。"
		c.JSON(http.StatusBadRequest, res)

		res = nil
		c.Abort()
		return
	}

	res[0] = getResponse(max, str)
	if res[0].Error == "" {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, res)
	}

	res = nil
}

// summary => 複数の要求に対してレスポンスを返します
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getMultiHashids(c *gin.Context) {
	maxprm := c.Param("max")
	max := getMaxValue(maxprm)

	strs := c.PostFormArray("sources[]")
	if len(strs) < 1 || len(strs) > LIMIT_CONTENTS {
		t := make([]HashSet, 1)

		if len(strs) < 1 {
			t[0].Error = "コンテンツが存在しません。"

		} else if len(strs) > LIMIT_CONTENTS {
			s := "要求コンテンツが%d件を超えています。"
			t[0].Error = fmt.Sprintf(s, LIMIT_CONTENTS)
		}

		c.JSON(http.StatusBadRequest, t)

		t = nil
		c.Abort()
		return
	}

	res := make([]HashSet, len(strs))
	for i, str := range strs {
		res[i] = getResponse(max, str)
	}

	// 複数要求時はエラーがあっても200で返す
	c.JSON(http.StatusOK, res)
	res = nil
}

// summary => 変換元文字列長の最大値を導出します
// param::maxStr => Paramからの流入値
// return::int => 最大値
/////////////////////////////////////////
func getMaxValue(maxStr string) int {
	max, err := strconv.Atoi(maxStr)

	// 変換に失敗 OR :maxのパラメータが存在しない場合、デフォルト値
	if err != nil {
		max = DEFAULT_MAX
		return max
	}

	if max < 1 {
		// 0以下の数値が入っている場合、デフォルト値
		max = DEFAULT_MAX

	} else if max > LIMIT_MAX {
		// 上限値より大きい値が入っている場合、強制的に上限値
		max = LIMIT_MAX
	}

	// 基本的に変換元文字列長の最大値の指定自由度は低い
	return max
}

// summary => レスポンスとして返す構造体を生成します
// param::max => 最大値
// param::input => 変換元の文字列
// return::HashSet => HashSet構造体
/////////////////////////////////////////
func getResponse(max int, input string) HashSet {
	res := HashSet{
		Source:  input,
		Hashids: "",
		Error:   "",
	}

	// max「以下」か（境界値バグテストしっかり）
	if len([]rune(input)) <= max {
		// Hashids生成
		h := GetShortHashids(input)

		if h == "" {
			res.Error = "文字列からHashidsへの変換に失敗しました"
		} else {
			res.Hashids = h
		}

	} else {
		t := "%d文字を超えた文字列が指定されています。"
		res.Error = fmt.Sprintf(t, max)
	}

	return res
}


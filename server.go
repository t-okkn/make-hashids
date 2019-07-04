package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
)

const (
	LIMIT_CONTENTS int = 1000
	LIMIT_MAX int = 4096
	DEFAULT_MAX int = 256
)

type HashSet struct {
	Error   string `json:"Error"`
	Hashids string `json:"Hashids"`
	Source  string `json:"Source"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/hashids", getSingleHashids)
	router.GET("/hashids/max/:max", getSingleHashids)
	router.POST("/hashids", getMultiHashids)
	router.POST("/hashids/max/:max", getMultiHashids)

	 return router
}

func getSingleHashids(c *gin.Context) {
	maxprm := c.Param("max")
	max := *getMaxValue(maxprm)

	res := make([]HashSet, 1)
	str := c.DefaultQuery("string", "")

	if str == "" {
		res[0].Error = "文字列が読み取れませんでした。"
		c.JSON(http.StatusBadRequest, res)

		res = nil
		return
	}

	res[0] = getResponse(&max, str)
	if res[0].Error == "" {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, res)
	}

	res = nil
}

func getMultiHashids(c *gin.Context) {
	maxprm := c.Param("max")
	max := *getMaxValue(maxprm)

	strs := c.PostFormArray("strings[]")
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
		return
	}

	res := make([]HashSet, len(strs))
	for i, str := range strs {
		res[i] = getResponse(&max, str)
	}

	c.JSON(http.StatusOK, res)
	res = nil
}

func getMaxValue(maxStr string) *int {
	max, err := strconv.Atoi(maxStr)

	if err != nil {
		max = DEFAULT_MAX
		return &max
	}

	if max < 1 {
		max = DEFAULT_MAX
	} else if max > LIMIT_MAX {
		max = LIMIT_MAX
	}

	return &max
}

func getResponse(max *int, input string) HashSet {
	res := HashSet{}
	res.Source = input

	if len([]rune(input)) <= *max {
		t := Str2Uints(input)
		res.Hashids = CreateHashids(t)

	} else {
		tmp := "%d文字以上の文字列が指定されています。"
		res.Error = fmt.Sprintf(tmp, *max)
	}

	return res
}


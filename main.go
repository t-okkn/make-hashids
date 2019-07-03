package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	s := gin.Default()

	s.GET("", getSingleHashids)
	s.GET("/max/:max", getSingleHashids)
	s.POST("", getMultiHashids)

	 // 0.0.0.0:8080 でサーバーを立てます。
	s.Run()
}

func getSingleHashids(c *gin.Context) {
	maxprm := c.Param("max")
	fmt.Println(maxprm)

	max, err := strconv.Atoi(maxprm)
	if err != nil {
		max = 256
	}

	str := c.DefaultQuery("string", "")
	if str == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"string": "",
			"hashids": "",
			"salt": "",
			"error": "文字列が読み取れませんでした。",
		})

		return
	}

	fmt.Printf("文字列長:%d\n", len([]rune(str)))
	fmt.Printf("最大値:%d\n", max)
	if max <= 4096 && len([]rune(str)) <= max {
		t := Str2Uints(str)
		hash, salt := CreateHashids(t)

		c.JSON(http.StatusOK, gin.H{
			"string": str,
			"hashids": hash,
			"salt": salt,
			"error": "",
		})

	} else {
		if max > 4096 {
			c.JSON(http.StatusBadRequest, gin.H{
				"string": str,
				"hashids": "",
				"salt": "",
				"error": "文字列長の最大値は4096以下で指定してください。",
			})

		} else {
			errstr := fmt.Sprintf("%d文字以上の文字列が指定されています。", max)

			c.JSON(http.StatusBadRequest, gin.H{
				"string": str,
				"hashids": "",
				"salt": "",
				"error": errstr,
			})
		}
	}
}

func getMultiHashids(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Post",
	})
}


package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func RemoveRepeate(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

func FindIndexInArray(val string, arr []string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

func ContextError(c *gin.Context, message string, err error) {
	fmt.Println(err)
	c.JSON(200, gin.H{"err": message})
}

type PageResponse struct {
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

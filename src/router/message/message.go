package message

import (
	"fmt"
	"github.com/easy-say-api/src/model"
	"github.com/easy-say-api/src/viewModel"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Get(c *gin.Context) {
	msg := model.Message{Id: c.Param("id")}
	err := msg.Get()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"data": msg})
}

func GetList(c *gin.Context) {
	msg := model.Message{}
	list, err := msg.GetList(50, 0)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"data": list})
}

func Add(c *gin.Context) {
	var form viewModel.MessageForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"err": "参数错误"})
		return
	}

	uid, _ := c.Cookie("openid")
	msg := model.SimpleMessage{
		Id:      uuid.New().String(),
		Content: form.Content,
		UserId:  uid,
	}
	err := msg.Add()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"data": msg})
}

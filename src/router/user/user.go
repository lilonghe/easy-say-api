package user

import (
	"github.com/easy-say-api/src/model"
	"github.com/easy-say-api/src/utils"
	"github.com/easy-say-api/src/viewModel"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
)

func Login(c *gin.Context) {
	var form viewModel.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"err": "参数错误"})
		return
	}

	hashPass := utils.GetMD5Hash(utils.GetMD5Hash(form.Password))
	user := model.User{Username: form.Username, Password: hashPass}
	if err := user.Login(); err != nil {
		c.JSON(200, gin.H{"err": "登录失败"})
		return
	}

	store := ginsession.FromContext(c)
	store.Set("openid", user.Id)
	err := store.Save()
	if err != nil {
		c.JSON(200, gin.H{"err": "登录失败,请稍后再试"})
		return
	}
	c.JSON(200, gin.H{"data": user})
}

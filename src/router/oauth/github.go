package oauth

import (
	"fmt"
	"github.com/easy-say-api/src/config"
	"github.com/easy-say-api/src/model"
	"github.com/easy-say-api/src/utils"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUserResponse struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
}

func GithubAuth(c *gin.Context) {
	code := c.Query("code")

	tokenResp := GithubTokenResponse{}
	{
		// get token
		req := gorequest.New().Post("https://github.com/login/oauth/access_token?" +
			fmt.Sprintf(`client_id=%s&`, config.Global.GithubAppId) +
			fmt.Sprintf(`client_secret=%s&`, config.Global.GithubAppSecret) +
			fmt.Sprintf(`code=%s`, code))
		req.Header.Set("Accept", "application/json")
		_, _, err := req.EndStruct(&tokenResp)
		if err != nil && len(err) > 0 {
			utils.ContextError(c, "请求失败", err[0])
			return
		}
	}
	userResp := GithubUserResponse{}
	{
		// get user
		req := gorequest.New().Get("https://api.github.com/user")
		req.Header.Set("Authorization", "token "+tokenResp.AccessToken)
		_, _, err := req.EndStruct(&userResp)
		if err != nil && len(err) > 0 {
			utils.ContextError(c, "请求失败", err[0])
			return
		}
	}

	{
		user := model.User{GithubId: userResp.Id}
		err := user.GetByGithubId()
		if err != nil {
			utils.ContextError(c, "请求失败", err)
			return
		}
		if user.Id != "" {
			// 走登录
			store := ginsession.FromContext(c)
			store.Set("openid", user.Id)
			err := store.Save()
			if err != nil {
				utils.ContextError(c, "登录失败,请稍后再试", err)
				return
			}
			c.JSON(200, gin.H{"data": user})
			return
		} else {
			// 新建用户
			user.Avatar = userResp.AvatarUrl
			user.Nickname = userResp.Name
			user.Username = userResp.Login
			user.Email = userResp.Email
			user.Id = uuid.New().String()
			user.BriefIntro = userResp.Bio
			// 暂时忽略 username 和 email 会重复的问题
			err := user.Add()
			if err != nil {
				utils.ContextError(c, "登录失败,请稍后再试", err)
				return
			}
			store := ginsession.FromContext(c)
			store.Set("openid", user.Id)
			err = store.Save()
			if err != nil {
				utils.ContextError(c, "登录失败,请稍后再试", err)
				return
			}
			c.JSON(200, gin.H{"data": user})
			return
		}
	}

	c.JSON(200, gin.H{})
}

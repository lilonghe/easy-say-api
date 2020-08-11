package message

import (
	"fmt"
	"github.com/easy-say-api/src/model"
	"github.com/easy-say-api/src/utils"
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

	// 填充当前用户点赞的列表
	if err == nil {
		ids := make([]string, 0)
		for _, v := range list {
			ids = append(ids, v.Id)
		}
		uid, _ := c.Get("openid")
		ids, err = model.UserFavoriteMessage{}.CheckUserInList(ids, uid.(string))
		for i, v := range list {
			index := utils.FindIndexInArray(v.Id, ids)
			if index != -1 {
				list[i].Liked = true
			}
		}
	}

	c.JSON(200, gin.H{"data": list})
}

func Add(c *gin.Context) {
	var form viewModel.MessageForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"err": "参数错误"})
		return
	}

	uid, _ := c.Get("openid")
	msg := model.SimpleMessage{
		Id:      uuid.New().String(),
		Content: form.Content,
		UserId:  uid.(string),
	}
	err := msg.Add()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"data": msg})
}

func Like(c *gin.Context) {
	var form viewModel.LikeMessageForm
	if err := c.ShouldBindJSON(&form); err != nil {
		utils.ContextError(c, "参数错误", err)
		return
	}

	uid, _ := c.Get("openid")
	// 取消点赞
	if *form.Unlike {
		obj := model.UserFavoriteMessage{
			UserId:    uid.(string),
			MessageId: *form.MessageId,
			IsEnable:  false,
		}
		err := obj.Del()
		if err != nil {
			c.JSON(200, gin.H{"err": "操作失败"})
			return
		}
		c.JSON(200, gin.H{})
		return
	}

	obj := model.UserFavoriteMessage{
		UserId:    uid.(string),
		MessageId: *form.MessageId,
		IsEnable:  true,
	}

	// 检测是否已存在
	err := obj.GetByUserAndMessage()
	if err != nil {
		c.JSON(200, gin.H{"err": "操作失败"})
		return
	}
	if obj.Id != 0 {
		c.JSON(200, gin.H{"err": "已点赞"})
		return
	}

	// 执行添加
	err = obj.Add()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}
	c.JSON(200, gin.H{})
}

func GetCommentListByMessage(c *gin.Context) {
	messageId := c.Query("message_id")
	comments, err := model.Comment{}.GetListByMessage(messageId)
	if err != nil {
		c.JSON(200, gin.H{"err": "获取失败"})
		return
	}

	c.JSON(200, gin.H{"data": utils.PageResponse{
		Total: len(comments),
		List:  comments,
	}})
}

func AddComment(c *gin.Context) {
	var form viewModel.MessageCommentForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"err": "参数错误"})
		return
	}
	uid, _ := c.Get("openid")
	obj := model.SimpleComment{
		Id:        uuid.New().String(),
		UserId:    uid.(string),
		MessageId: form.MessageId,
		Content:   form.Content,
		IsEnable:  true,
	}

	if len(form.Content) > 140 {
		utils.ContextError(c, "最多140个文字", nil)
		return
	}

	// 检查 message 是否有效
	msg := model.Message{Id: form.MessageId}
	err := msg.Get()
	if err != nil {
		utils.ContextError(c, "执行失败", err)
		return
	}

	err = obj.Add()
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{"err": "添加失败"})
		return
	}
	c.JSON(200, gin.H{"data": obj})
}

func DelComment(c *gin.Context) {
	var form viewModel.DelMessageCommentForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(200, gin.H{"err": "参数错误"})
		return
	}
	uid, _ := c.Get("openid")
	obj := model.Comment{
		Id:     form.CommentId,
		UserId: uid.(string),
	}
	err := obj.Get()
	if err != nil {
		c.JSON(200, gin.H{"err": "执行失败"})
		return
	}
	if obj.MessageId == "" {
		c.JSON(200, gin.H{"err": "无效数据"})
		return
	}

	obj.IsEnable = false
	err = obj.Del()
	if err != nil {
		c.JSON(200, gin.H{"err": "删除失败"})
		return
	}
	c.JSON(200, gin.H{})
}

package app

import (
	"strconv"

	"github.com/Linxhhh/webook/internal/service"
	"github.com/Linxhhh/webook/pkg/jwts"
	"github.com/Linxhhh/webook/pkg/res"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowHandler struct {
	svc *service.FollowService
}

func NewFollowHandler(svc *service.FollowService) *FollowHandler {
	return &FollowHandler{
		svc: svc,
	}
}

func (hdl *FollowHandler) RegistryRouter(router *gin.Engine) {
	ur := router.Group("userRelation")
	ur.POST("follow", hdl.Follow)
	ur.GET("follow", hdl.FollowData)
}

func (hdl *FollowHandler) Follow(ctx *gin.Context) {

	// 绑定参数
	type Req struct {
		Id     int64 `json:"id"`     // id 表示被关注的用户 id
		Follow bool  `json:"follow"` // true 表示关注，false 表示取消
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg("参数错误", ctx)
		return
	}

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)

	var err error
	if req.Follow {
		err = hdl.svc.Follow(ctx, claims.UserId, req.Id)
	} else {
		err = hdl.svc.CancelFollow(ctx, claims.UserId, req.Id)
	}
	if err != nil {
		res.FailWithMsg("系统错误", ctx)
		return
	}
	res.OKWithMsg("操作成功", ctx)
}

func (hdl *FollowHandler) FollowData(ctx *gin.Context) {

	// 获取用户 Token
	_claims, _ := ctx.Get("claims")
	claims := _claims.(*jwts.CustomClaims)
	uid := claims.UserId

	// 绑定参数
	id := ctx.Query("id")
	if id != "" {
		uid, _ = strconv.ParseInt(id, 10, 64)
	}

	// 获取用户关系数据
	data, err := hdl.svc.GetFollowData(ctx, uid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			res.OKWithData(data, ctx)
			return
		}
		res.FailWithMsg("系统错误", ctx)
		return
	}

	if uid != claims.UserId {
		data.IsFollowed, err = hdl.svc.GetFollowed(ctx, claims.UserId, uid)
		if err != nil {
			res.FailWithMsg("系统错误", ctx)
			return
		}
	}
  
	// 返回响应
	res.OKWithData(data, ctx)
}
package service

import (
	"context"

	v1 "github.com/2720851545/realworld-golang-gf/api/v1"
	"github.com/2720851545/realworld-golang-gf/internal/service/internal/dao"
	"github.com/2720851545/realworld-golang-gf/utility"
	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type IUserService interface {
	// Register(ctx context.Context, req *v1.UserRegisterReq) (i int64, err error)
	Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error)
	CurrentUser(ctx context.Context, req *v1.CurrentUserReq) (res *v1.CurrentUserRes, err error)
}

type userImpl struct{}

func UserService() IUserService {
	return IUserService(&userImpl{})
}

func (s *userImpl) CurrentUser(ctx context.Context, req *v1.CurrentUserReq) (res *v1.CurrentUserRes, err error) {
	var (
		token string
		mc    jwt.MapClaims
	)
	mc, token, err = authService.GetClaimsFromJWT(ctx)
	if err != nil {
		return
	}
	g.Log().Info(ctx, mc, token, err)
	if id, ok := mc["id"]; ok {
		res = new(v1.CurrentUserRes)
		res.User.Token = token
		dao.User.Ctx(ctx).Where("id = ?", id).Scan(&res.User)
	} else {
		err = gerror.New("jwt token error, id not found")
	}
	return
}

func (s *userImpl) Register(ctx context.Context, req *v1.UserRegisterReq) (res *v1.UserRegisterRes, err error) {
	req.User.Password = utility.EntryPassword(req.User.Password)
	if g.IsEmpty(req.User.Image) {
		req.User.Image = "https://api.realworld.io/images/smiley-cyrus.jpeg"
	}

	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			id, err = dao.User.Ctx(ctx).TX(tx).InsertAndGetId(req.User)
			if err != nil {
				panic(err)
			}
			g.Log().Info(ctx, "新用户的id=", id)
			res = new(v1.UserRegisterRes)
			err = dao.User.Ctx(ctx).Where("id = ?", id).Scan(&res.User)
		})
		return err
	})

	if err == nil {
		r := g.RequestFromCtx(ctx)
		r.SetCtxVar("Model", "Register")
		r.SetCtxVar("User", map[string]interface{}{
			"id":       id,
			"username": res.User.Username,
		})
		res.User.Token, _ = authService.LoginHandler(ctx)
	}
	return res, err
}

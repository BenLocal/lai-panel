package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/benlocal/lai-panel/pkg/crypto"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	// TODO get value from config
	JwtSecret        = "lai-panel-jwt-secret"
	CtxKeyUserClaims = "user_claims"
)

func (h *BaseHandler) LoginHandler(ctx context.Context, c *app.RequestContext) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req loginRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	cp, err := base64.StdEncoding.DecodeString(req.Password)
	if err != nil {
		c.Error(errors.New("invalid password"))
		return
	}

	userRepository := h.UserRepository()
	user, err := userRepository.GetByUsername(req.Username)
	if err != nil || user == nil {
		c.Error(errors.New("user not found"))
		return
	}

	password, err := crypto.Decrypt(user.Password)
	if err != nil {
		c.Error(errors.New("invalid password"))
		return
	}
	if password != string(cp) {
		c.Error(errors.New("invalid password"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		// TODO get value from config
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
		"kid":   uuid.New().String(),
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"iss":   "lai-panel",
		"aud":   "lai-panel",
		"sub":   user.ID,
		"jti":   uuid.New().String(),
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		c.Error(err)
		return
	}

	type loginResponse struct {
		Token    string `json:"token"`
		UserId   int64  `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	c.JSON(http.StatusOK, SuccessResponse(&loginResponse{
		Token:    tokenString,
		UserId:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Name:     user.Name,
	}))
}

func (h *BaseHandler) ChangePasswordHandler(ctx context.Context, c *app.RequestContext) {
	type changePasswordRequest struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var req changePasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		c.Error(errors.New("old password and new password are required"))
		return
	}

	user, err := h.GetUserInfo(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}

	db, err := h.UserRepository().GetById(user.ID)
	if err != nil {
		c.Error(err)
		return
	}
	if db == nil {
		c.Error(errors.New("user not found"))
		return
	}

	oldPassword, err := crypto.Decrypt(db.Password)
	if err != nil {
		c.Error(errors.New("invalid old password"))
		return
	}

	oldPasswordBytes, err := base64.StdEncoding.DecodeString(req.OldPassword)
	if err != nil {
		c.Error(errors.New("invalid old password"))
		return
	}
	if oldPassword != string(oldPasswordBytes) {
		c.Error(errors.New("invalid old password"))
		return
	}

	newPwdBytes, err := base64.StdEncoding.DecodeString(req.NewPassword)
	if err != nil {
		c.Error(errors.New("invalid new password"))
		return
	}

	userRepository := h.UserRepository()
	err = userRepository.UpdatePassword(user.ID, string(newPwdBytes))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func (h *BaseHandler) AuthMiddleware(ctx context.Context, c *app.RequestContext) {
	path := string(c.Path())
	if path == "/api/auth/login" {
		c.Next(ctx)
		return
	}

	token := string(c.GetHeader("Authorization"))
	if token == "" {
		c.Error(errors.New("unauthorized"))
		c.Abort()
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		c.Error(errors.New("unauthorized"))
		c.Abort()
		return
	}

	var mc jwt.MapClaims
	claims, err := jwt.ParseWithClaims(token, &mc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if token.Header["alg"] == "none" {
			return nil, errors.New("token algorithm is not supported")
		}
		return []byte(JwtSecret), nil
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if !claims.Valid {
		c.Error(errors.New("token is invalid"))
		c.Abort()
		return
	}

	exp, err := mc.GetExpirationTime()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	if exp.Before(time.Now()) {
		c.Error(errors.New("token is expired"))
		c.Abort()
		return
	}

	nbf, err := mc.GetNotBefore()
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	if nbf.After(time.Now()) {
		c.Error(errors.New("token is not yet valid"))
		c.Abort()
		return
	}
	c.Set(CtxKeyUserClaims, mc)
	c.Next(ctx)
}

func (h *BaseHandler) GetUserInfo(ctx context.Context, c *app.RequestContext) (*model.User, error) {
	mc := c.Value(CtxKeyUserClaims).(jwt.MapClaims)
	userId, err := getUserIdFromClaims(mc)
	if err != nil {
		return nil, err
	}
	if userId == nil {
		return nil, errors.New("user id is invalid")
	}

	keys := []string{"username", "role", "name", "email"}
	vs := []string{}
	for _, key := range keys {
		v, err := parseString(mc, key)
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, errors.New(fmt.Sprintf("%s is invalid", key))
		}
		vs = append(vs, *v)
	}

	user := &model.User{
		ID:       *userId,
		Username: vs[0],
		Role:     vs[1],
		Name:     vs[2],
		Email:    vs[3],
	}
	return user, nil
}

func parseString(m jwt.MapClaims, key string) (*string, error) {
	var (
		ok  bool
		raw any
		v   string
	)
	raw, ok = m[key]
	if !ok {
		return nil, nil
	}

	v, ok = raw.(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is invalid", key))
	}

	return &v, nil
}

func getUserIdFromClaims(mc jwt.MapClaims) (*int64, error) {
	return parseInt64(mc, "userId")
}

func parseInt64(m jwt.MapClaims, key string) (*int64, error) {
	v, ok := m[key]
	if !ok {
		return nil, nil
	}

	switch exp := v.(type) {
	case float64:
		if exp == 0 {
			return nil, nil
		}
		val := int64(exp)
		return &val, nil
	case json.Number:
		v, _ := exp.Int64()
		return &v, nil
	}

	return nil, errors.New(fmt.Sprintf("%s is invalid", key))
}

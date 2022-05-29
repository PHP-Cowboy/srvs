package handler

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"shop-srvs/user_srv/global"
	"shop-srvs/user_srv/model"
	"shop-srvs/user_srv/proto/proto"
	"strings"
	"time"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToUserInfoResponse(u model.User) (userInfo proto.UserInfoResponse) {
	userInfo = proto.UserInfoResponse{
		Id:       u.Id,
		Mobile:   u.Mobile,
		NickName: u.Nickname,
		PassWord: u.Password,
		Gender:   u.Gender,
		Role:     u.Role,
	}
	if u.Birthday != nil {
		userInfo.BirthDay = uint32(u.Birthday.Unix())
	}
	return
}

func (u *UserServer) GetUserList(c context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	db := global.DB

	var users []model.User

	result := db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	var rsp proto.UserListResponse

	rsp.Total = uint32(result.RowsAffected)

	db.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, u := range users {
		userInfo := ModelToUserInfoResponse(u)
		rsp.Data = append(rsp.Data, &userInfo)
	}

	return &rsp, nil
}

func (u *UserServer) GetUserByMobile(c context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户未找到")
	}

	userInfo := ModelToUserInfoResponse(user)

	return &userInfo, nil
}

func (u *UserServer) GetUserById(c context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户未找到")
	}

	userInfo := ModelToUserInfoResponse(user)

	return &userInfo, nil
}

func (u *UserServer) CreateUser(c context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).Find(&user)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)

	user.Password = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.Mobile = req.Mobile
	user.Nickname = req.NickName

	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	rsp := ModelToUserInfoResponse(user)

	return &rsp, nil
}

func (u *UserServer) UpdateUser(c context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User

	db := global.DB

	result := db.First(&user, req.Id)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户已存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)

	user.Nickname = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = db.Save(user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (u *UserServer) CheckPassWord(c context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}

	pwdSlice := strings.Split(req.EncryptedPassWord, "$")

	check := password.Verify(req.PassWord, pwdSlice[1], pwdSlice[2], options)

	return &proto.CheckResponse{Status: check}, nil
}

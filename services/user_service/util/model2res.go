package util

import (
	"github.com/zhaohaihang/user_service/model"
	"github.com/zhaohaihang/user_service/proto"
)

// ModelToResponse
func ModelToResponse(user model.User) *proto.UserInfoResponse {
	userInfoResponse := proto.UserInfoResponse{
		Id:       int32(user.Model.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResponse.Birthday = uint64(user.Birthday.Unix())
	}
	return &userInfoResponse
}

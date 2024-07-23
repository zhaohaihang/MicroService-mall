package handler

import "github.com/zhaohaihang/user_service/proto"

type UserService struct {
	proto.UnimplementedUserServer
}

const (
	SERVICE_NAME = "[User_Service]"
)

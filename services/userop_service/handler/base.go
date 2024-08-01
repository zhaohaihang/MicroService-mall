package handler

import "github.com/zhaohaihang/userop_service/proto"

type UserOpService struct {
	proto.UnimplementedAddressServer
	proto.UnimplementedUserFavoriteServer
	proto.UnimplementedMessageServer
}

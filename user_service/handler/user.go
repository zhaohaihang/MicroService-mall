package handler

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/zhaohaihang/user_service/global"
	"github.com/zhaohaihang/user_service/model"
	"github.com/zhaohaihang/user_service/proto"
	"github.com/zhaohaihang/user_service/util"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserList 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, pageInfoRequest *proto.PageInfoRequest) (*proto.UserListResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetUserList", "request", pageInfoRequest)

	parentSpan := opentracing.SpanFromContext(ctx)
	userListSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	
	
	users := []model.User{}
	result := global.DB.Find(&users)
	
	pageNum := pageInfoRequest.PageNum
	pageSize := pageInfoRequest.PageSize
	offset := util.Paginate(int(pageNum), int(pageSize))

	pageUsers := []model.User{}
	global.DB.Offset(offset).Limit(int(pageSize)).Find(&pageUsers)

	response := &proto.UserListResponse{}
	response.Total = int32(result.RowsAffected)
	for _, user := range pageUsers {
		response.Data = append(response.Data, util.ModelToResponse(user))
	}
	userListSpan.Finish()
	return response, nil
}

// GetUserByMobile  通过电话号码获取用户信息
func (s *UserService) GetUserByMobile(ctx context.Context, mobileRequest *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetUserByMobile", "request", mobileRequest)
	parentSpan := opentracing.SpanFromContext(ctx)
	response := &proto.UserInfoResponse{}
	getUserByMobileSpan := opentracing.GlobalTracer().StartSpan("GetUserByMobile", opentracing.ChildOf(parentSpan.Context()))
	var user model.User
	mobile := mobileRequest.Mobile
	result := global.DB.Where("mobile=?", mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}
	getUserByMobileSpan.Finish()
	response = util.ModelToResponse(user)
	return response, nil
}

// GetUserById 通过ID获取用户信息
func (s *UserService) GetUserById(ctx context.Context, idRequest *proto.IdRequest) (*proto.UserInfoResponse, error) {

	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetUserById", "request", idRequest)
	parentSpan := opentracing.SpanFromContext(ctx)
	getUserByIdSpan := opentracing.GlobalTracer().StartSpan("GetUserById", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.UserInfoResponse{}

	var user model.User
	id := idRequest.Id
	result := global.DB.Where("id=?", id).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}
	getUserByIdSpan.Finish()
	response = util.ModelToResponse(user)
	return response, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, createUserInfoRequest *proto.CreateUserInfoRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateUser", "request", createUserInfoRequest)
	parentSpan := opentracing.SpanFromContext(ctx)
	response := &proto.UserInfoResponse{}
	mobile := createUserInfoRequest.Mobile
	createUserSpan := opentracing.GlobalTracer().StartSpan("CreateUser", opentracing.ChildOf(parentSpan.Context()))
	var user model.User
	result := global.DB.Where("mobile=?", mobile)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "user has exists")
	}
	user.Mobile = createUserInfoRequest.Mobile
	user.NickName = createUserInfoRequest.NickName

	password := createUserInfoRequest.Password
	encryptPassword := util.EncryptPassword(password)
	user.Password = encryptPassword

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	createUserSpan.Finish()
	response = util.ModelToResponse(user)
	return response, nil
}

// UpdateUser 更新用户
func (s UserService) UpdateUser(ctx context.Context, UpdateUserInfoRequest *proto.UpdateUserInfoRequest) (*proto.UpdateResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateUser", "request", UpdateUserInfoRequest)
	parentSpan := opentracing.SpanFromContext(ctx)
	response := &proto.UpdateResponse{}
	var user model.User
	updateUserSpan := opentracing.GlobalTracer().StartSpan("update_user", opentracing.ChildOf(parentSpan.Context()))
	result := global.DB.First(&user, UpdateUserInfoRequest.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}
	birthDay := time.Unix(int64(UpdateUserInfoRequest.Birthday), 0)
	user.NickName = UpdateUserInfoRequest.NickName
	user.Birthday = &birthDay
	user.Gender = UpdateUserInfoRequest.Gender

	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	updateUserSpan.Finish()
	response.Success = true
	return response, nil
}

// CheckPassword 检查用户密码
func (s UserService) CheckPassword(ctx context.Context, checkPasswordRequest *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CheckPassword", "request", checkPasswordRequest)
	parentSpan := opentracing.SpanFromContext(ctx)
	checkPasswordSpan := opentracing.GlobalTracer().StartSpan("CheckPassword", opentracing.ChildOf(parentSpan.Context()))
	response := &proto.CheckPasswordResponse{}
	password := checkPasswordRequest.Password
	EncryptedPassword := checkPasswordRequest.EncryptedPassword
	response.Success = util.VerifyPassword(EncryptedPassword, password)
	checkPasswordSpan.Finish()
	return response, nil
}

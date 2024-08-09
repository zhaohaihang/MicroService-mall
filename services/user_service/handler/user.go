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
	getUserListSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	
	response := &proto.UserListResponse{}
	// 1. 获取用户总数
	users := []model.User{}
	result := global.DB.Find(&users)
	response.Total = int32(result.RowsAffected)
	// 2. 分页查询
	pageNum := pageInfoRequest.PageNum
	pageSize := pageInfoRequest.PageSize
	offset := util.Paginate(int(pageNum), int(pageSize))
	pageUsers := []model.User{}
	global.DB.Offset(offset).Limit(int(pageSize)).Find(&pageUsers)
	for _, user := range pageUsers {
		response.Data = append(response.Data, util.ModelToResponse(user))
	}

	getUserListSpan.Finish()
	return response, nil
}

// GetUserByMobile  通过电话获取用户信息
func (s *UserService) GetUserByMobile(ctx context.Context, mobileRequest *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetUserByMobile", "request", mobileRequest)
	
	parentSpan := opentracing.SpanFromContext(ctx)	
	getUserByMobileSpan := opentracing.GlobalTracer().StartSpan("GetUserByMobile", opentracing.ChildOf(parentSpan.Context()))
	
	user := model.User{}
	result := global.DB.Where("mobile=?", mobileRequest.Mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}

	response := util.ModelToResponse(user)
	
	getUserByMobileSpan.Finish()
	return response, nil
}

// GetUserById 通过ID获取用户信息
func (s *UserService) GetUserById(ctx context.Context, idRequest *proto.IdRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "GetUserById", "request", idRequest)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	getUserByIdSpan := opentracing.GlobalTracer().StartSpan("GetUserById", opentracing.ChildOf(parentSpan.Context()))

	user := model.User{}
	result := global.DB.Where("id=?", idRequest.Id).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}
	response := util.ModelToResponse(user)
	
	getUserByIdSpan.Finish()
	return response, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, createUserInfoRequest *proto.CreateUserInfoRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CreateUser", "request", createUserInfoRequest)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	createUserSpan := opentracing.GlobalTracer().StartSpan("CreateUser", opentracing.ChildOf(parentSpan.Context()))

	user := model.User{}
	// 1.校验手机号是否占用
	result := global.DB.Where("mobile=?", createUserInfoRequest.Mobile)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "user has exists")
	}
	// 2.创建用户
	user.Mobile = createUserInfoRequest.Mobile
	user.NickName = createUserInfoRequest.NickName
	user.Password = util.EncryptPassword( createUserInfoRequest.Password)
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	response := util.ModelToResponse(user)

	createUserSpan.Finish()
	return response, nil
}

// UpdateUser 更新用户
func (s UserService) UpdateUser(ctx context.Context, UpdateUserInfoRequest *proto.UpdateUserInfoRequest) (*proto.UpdateResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "UpdateUser", "request", UpdateUserInfoRequest)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	updateUserSpan := opentracing.GlobalTracer().StartSpan("update_user", opentracing.ChildOf(parentSpan.Context()))
	// 1.校验用户是否存在
	user :=  model.User{}
	result := global.DB.First(&user, UpdateUserInfoRequest.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user is not exists")
	}
	// 2.更新用户
	birthDay := time.Unix(int64(UpdateUserInfoRequest.Birthday), 0)
	user.NickName = UpdateUserInfoRequest.NickName
	user.Birthday = &birthDay
	user.Gender = UpdateUserInfoRequest.Gender
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	response := &proto.UpdateResponse{
		Success: true,
	}

	updateUserSpan.Finish()
	return response, nil
}

// CheckPassword 检查用户密码
func (s UserService) CheckPassword(ctx context.Context, checkPasswordRequest *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	zap.S().Infow("Info", "service", SERVICE_NAME, "method", "CheckPassword", "request", checkPasswordRequest)
	
	parentSpan := opentracing.SpanFromContext(ctx)
	checkPasswordSpan := opentracing.GlobalTracer().StartSpan("CheckPassword", opentracing.ChildOf(parentSpan.Context()))

	response := &proto.CheckPasswordResponse{
		Success : util.VerifyPassword(checkPasswordRequest.EncryptedPassword,  checkPasswordRequest.Password),
	}

	checkPasswordSpan.Finish()
	return response, nil
}

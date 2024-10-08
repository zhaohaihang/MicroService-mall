package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/zhaohaihang/userop_service/global"
	"github.com/zhaohaihang/userop_service/model"
	"github.com/zhaohaihang/userop_service/proto"
)

func (s *UserOpService) MessageList(ctx context.Context, request *proto.MessageRequest) (*proto.MessageListResponse, error) {
	zap.S().Infow("Info", "method", "MessageList", "request", request)
	parentSpan := opentracing.SpanFromContext(ctx)
	messageListSpan := opentracing.GlobalTracer().StartSpan("MessageList", opentracing.ChildOf(parentSpan.Context()))

	response := &proto.MessageListResponse{}
	var messages []model.LeavingMessages
	var messageList []*proto.MessageResponse

	result := global.DB.Where(&model.LeavingMessages{User: request.UserId}).Find(&messages)
	if result.RowsAffected == 0 {
		zap.S().Warnw("Warning", "message", "查询地址数据为空", "request", request.UserId)
	}
	response.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          int32(message.ID),
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	response.Data = messageList
	messageListSpan.Finish()
	return response, nil
}

func (s *UserOpService) CreateMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message model.LeavingMessages
	parentSpan := opentracing.SpanFromContext(ctx)
	createMessageSpan := opentracing.GlobalTracer().StartSpan("CreateMessage", opentracing.ChildOf(parentSpan.Context()))

	message.User = request.UserId
	message.MessageType = request.MessageType
	message.Subject = request.Subject
	message.Message = request.Message
	message.File = request.File

	result := global.DB.Save(&message)
	if result.Error != nil {
		zap.S().Errorw("Error", "message", "创建地址失败", "err", result.Error)
		return nil, status.Errorf(codes.Internal, "创建地址失败")
	}
	createMessageSpan.Finish()
	return &proto.MessageResponse{Id: int32(message.ID)}, nil
}

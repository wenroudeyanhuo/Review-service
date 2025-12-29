// @program:     review-service
// @file:        review.go
// @author:      16574
// @create:      2025-12-26 17:35
// @description:

package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
)

type ReviewerRepo interface {
	//在data层实现
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewReply(context.Context, int64) (*model.ReviewReplyInfo, error)
	AuditReview(context.Context, *AuditParam) error
	AppealReview(context.Context, *AppealParam) error
	ListReviewByUserID(ctx context.Context, UserID int64, offset, limit int) ([]*model.ReviewInfo, error)
	AuditAppeal(context.Context, *AuditAppealParam) error
}
type ReviewerUsecase struct {
	repo ReviewerRepo
	log  *log.Helper
}

func NewReviewerUsecase(repo ReviewerRepo, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 创建评价
// 实现业务逻辑的地方
// service 层调用这个方法
func (uc *ReviewerUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("[biz]  CreateReview: %v", review)
	//1. 数据校验
	//1.1 参数基础校验：正常来说不应该在这一层，你在上一层或着框架层都应该能拦住（validate 参数校验）
	//1.2 参数业务校验：沾点带业务逻辑的参数校验，比如已经评价过的订单不能再创建评价
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbFalied("查询评价数据库失败")
	}
	if len(reviews) > 0 {
		return nil, v1.ErrorOrderReviewed("已经评价过,订单%d", review.OrderID)
	}
	//2. 生成review id    (使用雪花算法生成)  也可以直接介入公司内部的分布式id生成服务
	review.ReviewID = snowflake.GenID()
	//3. 查询订单和商品快照信息
	//实际业务场景下就需要查询订单服务和商家服务  比如通过RPC调用订单服务和商家服务
	//4. 拼装数据入库
	return uc.repo.SaveReview(ctx, review)
}

// GetReview 根据评价ID获取评价
func (uc *ReviewerUsecase) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Infof("[biz]  GetReview  reviewID : %d", reviewID)
	return uc.repo.GetReview(ctx, reviewID)
}

// CreateReply 创建评价回复
func (uc *ReviewerUsecase) CreateReply(ctx context.Context, param *ReplyParm) (*model.ReviewReplyInfo, error) {
	//调用data层创建一个评价回复
	uc.log.WithContext(ctx).Debugf("[biz]  Create Reply  param : %v", param)
	return uc.repo.SaveReply(ctx, &model.ReviewReplyInfo{
		ReviewID:  param.ReviewID,
		ReplyID:   snowflake.GenID(),
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
}

// 审核评价
func (uc *ReviewerUsecase) AuditReview(ctx context.Context, auditParam *AuditParam) error {
	uc.log.WithContext(ctx).Infof("[biz]  AuditReview  auditParam : %v", auditParam)
	return uc.repo.AuditReview(ctx, auditParam)
}

// 申诉评价
func (uc *ReviewerUsecase) AppealReview(ctx context.Context, appealParam *AppealParam) error {
	uc.log.WithContext(ctx).Infof("[biz] AppealReview param:%#v", appealParam)
	return uc.repo.AppealReview(ctx, appealParam)
}

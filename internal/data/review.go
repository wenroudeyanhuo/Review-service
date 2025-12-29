package data

import (
	"context"
	"errors"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"review-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewerRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo
func NewReviewerRepo(data *Data, logger log.Logger) biz.ReviewerRepo {
	return &reviewerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewerRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}
func (r *reviewerRepo) ListReviewByUserID(ctx context.Context, UserID int64, offset, limit int) ([]*model.ReviewInfo, error) {
	//TODO implement me
	panic("implement me")
}

// 根据订单id  查询评价
func (r *reviewerRepo) GetReviewByOrderID(ctx context.Context, OrderID int64) ([]*model.ReviewInfo, error) {
	//	Find() ([]*model.ReviewInfo, error)
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(OrderID)).Find()
}
func (r *reviewerRepo) GetReviewByID(ctx context.Context, ID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ID.Eq(ID)).First()
}
func (r *reviewerRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

// SaveReply  保存评价回复
func (r *reviewerRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	//业务逻辑实现的地方
	//1.数据校验
	//1.1 数据合法性校验（已回复的不允许再评价）
	//先用评价ID查库，看下是否已经回复
	review, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).First()
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, errors.New("该评价已回复")
	}
	//1.2.水平越权校验  （A 商家只能回复自己的，不能回复B商家的
	//例子：用户A删除订单     使用userID+orderID 当条件去查询
	if reply.StoreID != review.StoreID {
		return nil, errors.New("无权限")
	}
	//2.更新数据库中的数据  评价回复表和评价表要同时更新，涉及事务操作
	err = r.data.query.Transaction(func(tx *query.Query) error {
		//	回复表插入一条数据  save方法是保持有则更新无则创建
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply create reply fail.err:%v\n", err)
			return err
		}
		//	评价表修改hasreply字段
		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(reply.ReplyID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("SaveReply update review fail.err:%v\n", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	//3.返回
	return reply, nil
}

func (r *reviewerRepo) GetReviewReply(ctx context.Context, reviewID int64) (*model.ReviewReplyInfo, error) {
	return r.data.query.ReviewReplyInfo.WithContext(ctx).Where(r.data.query.ReviewReplyInfo.ReviewID.Eq(reviewID)).First()
}
func (r *reviewerRepo) AuditReview(ctx context.Context, param *biz.AuditParam) error {
	return nil
}
func (r *reviewerRepo) AppealReview(ctx context.Context, parms *biz.AppealParam) error {
	return nil
}
func (r *reviewerRepo) AuditAppeal(ctx context.Context, parms *biz.AuditAppealParam) error {
	return nil
}

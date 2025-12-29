// @program:     review-service
// @file:        param.go
// @author:      16574
// @create:      2025-12-29 15:09
// @description:

package biz

// ReplyParam 商家回复评价参数
type ReplyParm struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

// AuditParam 运营审核评价参数
type AuditParam struct {
	ReviewID  int64
	OpUser    string
	OpReason  string
	OpRemarks string
	Status    int32
}

// AppealParam 商家申诉评价的参数
type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}

type AuditAppealParam struct {
	AppealID  int64
	OpUser    string
	OpReason  string
	OpRemarks string
	Status    int32
}

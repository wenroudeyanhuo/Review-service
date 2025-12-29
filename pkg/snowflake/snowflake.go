// @program:     review-service
// @file:        snowflake.go
// @author:      16574
// @create:      2025-12-27 21:31
// @description:

package snowflake

import (
	"errors"
	sf "github.com/bwmarrin/snowflake"
	"time"
)

//利用雪花算法 生成id
/*
在公司中如果有生成唯一id的服务可以直接对接
本质上只是为了获取一个全局唯一id
*/
var (
	InvalidInitParmErr   = errors.New("snowflake 初始化失败，无效的startTime或machineID")
	InvalidTimeFormatErr = errors.New("snowflake 初始化失败，时间格式有误")
)
var node *sf.Node

// 雪花算法初始化配置
func Init(startTime string, machineID int64) (err error) {
	if len(startTime) == 0 || machineID <= 0 {
		return InvalidInitParmErr
	}
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return InvalidTimeFormatErr
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

// 生成一个id
func GenID() int64 {

	return node.Generate().Int64()
}

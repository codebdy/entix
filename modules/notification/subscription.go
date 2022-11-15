package notification

import (
	"github.com/graphql-go/graphql"
)

//1、无限通道存储事件源 ==== 不用无限通道，单独开一个线程发送消息
//2、单独线程分发事件
//3、所有订阅列表

func (m *SubscriptionModule) SubscriptionFields() []*graphql.Field {
	if m.app != nil {
		return []*graphql.Field{
			{
				Name: "unreadNoticationCounts",
				Type: graphql.Int,
				//Resolve: PublishMetaResolveFn(m.app),
			},
		}
	} else {
		return []*graphql.Field{}
	}
}

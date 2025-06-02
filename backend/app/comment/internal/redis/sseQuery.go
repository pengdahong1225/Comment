package redis

import (
	"context"
	"fmt"
)

// 根据roomId查询相关的sse服务
func QuerySseServiceByRoomId(roomId int64) ([]string, error) {
	key := fmt.Sprintf("SseWithRoomID_%d", roomId)
	result := rdb.SMembers(context.Background(), key)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return result.Result()
}

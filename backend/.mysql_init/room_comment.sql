use comment_server;

-- 直播间评论区
create table if not exists room_comment
(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    room_id BIGINT not NULL comment '直播间id',
    user_id BIGINT not NULL comment '评论用户id',
    user_name VARCHAR(64) DEFAULT '' comment '评论时用户名',
    content TEXT not NULL comment '评论内容',
    pub_stamp BIGINT NOT NULL COMMENT '时间戳',
    pub_region VARCHAR(64) DEFAULT '' COMMENT '发布地区',

    PRIMARY KEY(id),
    INDEX idx_room_user_id(room_id, user_id)
)engine = InnoDB charset = utf8mb4;
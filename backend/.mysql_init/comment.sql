use comment_server;

-- 两层结构评论表
create table if not exists comment(
    id BIGINT AUTO_INCREMENT,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delete_at TIMESTAMP,

    obj_id BIGINT NOT NULL COMMENT '对象id（题目id）',
    user_id BIGINT NOT NULL COMMENT '评论用户id',
    user_name VARCHAR(64) DEFAULT '' COMMENT '评论时用户名',
    user_avatar_url VARCHAR(256) DEFAULT '' COMMENT '评论时头像url',
    content TEXT NOT NULL COMMENT '评论内容',
    status tinyint DEFAULT 1 COMMENT '状态 1：正常 0：隐藏',
    reply_count INT DEFAULT 0 COMMENT '回复数量',
    like_count INT DEFAULT 0 COMMENT '点赞数量',
    child_count INT DEFAULT 0 COMMENT '两层结构，root才有子评论',
    pub_stamp BIGINT NOT NULL COMMENT '时间戳',
    pub_region VARCHAR(64) DEFAULT '' COMMENT '发布地区',

    is_root tinyint DEFAULT 1 COMMENT '是否是楼主 1：是 0：不是',
    root_id BIGINT NOT NULL COMMENT '楼主id',
    root_comment_id BIGINT NOT NULL COMMENT '楼主评论id',
    reply_id BIGINT NOT NULL COMMENT '回复用户id',
    reply_comment_id BIGINT NOT NULL COMMENT '回复评论id',
    reply_user_name VARCHAR(64) DEFAULT '' COMMENT '回复用户名',

    PRIMARY KEY(id),
    INDEX idx_obj_id(obj_id),
    INDEX idx_user_id(user_id),
    INDEX idx_root_id(root_id),
    INDEX idx_reply_id(reply_id),
    INDEX idx_like_count(like_count)
)engine = InnoDB charset = utf8mb4;

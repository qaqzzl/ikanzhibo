-- 创建库
create database if not exists `ikanzhibo`;

use `ikanzhibo`;

-- 用户会员
create table if not exists `user_member`(
    `member_id` int unsigned auto_increment primary key,
    `nickname` varchar(50) not null default '' comment '用户昵称',
    `gender` char(5) not null default 'wz' comment 'wz-未知, w-女, m-男, z-中性',
    `birthdate` int not null default 0 comment '出生日期',
    `avatar` varchar(255) not null default '' comment '头像',
    `signature` varchar(64) not null default '' comment '个性签名',
    `city` char(50) not null default '' comment '城市',
    `province` char(50) not null default '' comment '省份',
    `created_at` int not null default 0 comment '添加时间',
    `updated_at` int not null default 0 comment '修改时间',
    `deleted_at` int not null default 0 comment '删除时间',
    UNIQUE KEY `nickname` (`nickname`)
)engine=innodb default charset=utf8mb4 comment '用户会员';

-- 会员授权账号表
create table if not exists `user_auths`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `identity_type` char(20) not null comment '类型,wechat_applet,qq,wb,phone,number,email',
    `identifier` varchar(255) not null default '' comment '微信,QQ,微博opendid | 手机号,邮箱,账号',
    `credential` varchar(255) not null default '' comment '密码凭证（站内的保存密码，站外的不保存或保存access_token）',
    KEY `member_id` (`member_id`),
    UNIQUE KEY `identity_type_identifier` (`identity_type`,`identifier`) USING BTREE
)engine=innodb default charset=utf8mb4 comment '会员授权账号表';

-- 用户授权 token 表 ,这个表用redis比较好 , 也可以使用JWS
create table if not exists `user_auths_token`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `token` varchar(255) not null default '' comment 'token',
    `client` char(20) not null comment 'app,web,wechat_applet',
    `last_time` int not null comment '上次刷新时间',
    `status` tinyint(1) not null default 0 comment '1-其他设备强制下线',
    `created_at` int not null default 0 comment '添加时间',
    UNIQUE KEY `token` (`token`)
)engine=innodb default charset=utf8 comment '用户授权 token 表';


-- 直播用户关注表
create table if not exists `live_user_follow`(
    `id` int unsigned auto_increment primary key,
    `member_id` int not null comment '会员ID',
    `live_id` int not null comment '直播ID',
    `status` tinyint(1) not null default 1 comment '1-关注, 0-取消关注',
	`is_notice` tinyint(1) not null default 1 comment '是否通知',
	`send_notice_time` int not null default 0 comment '上次通知时间',
    `created_at` int not null default 0 comment '添加时间',
    `updated_at` int not null default 0 comment '修改时间',
    UNIQUE KEY `member_id_live_id` (`member_id`,`live_id`) USING BTREE
)engine=innodb default charset=utf8 comment '直播用户关注表';

-- 直播表(主播表)
create table if not exists `live`(
    `live_id` int unsigned AUTO_INCREMENT,
    `live_title` varchar(255) COLLATE utf8mb4_unicode_ci not null DEFAULT '' comment '直播标题',
    `live_anchortv_name` varchar(100) COLLATE utf8mb4_unicode_ci not null DEFAULT '' comment '主播名称',
    `live_anchortv_photo` varchar(255) not null DEFAULT '' comment '主播头像',
    `live_anchortv_sex` tinyint(1) not null DEFAULT 0 comment '主播性别 0-保密 1-女 2-男',
    `live_cover` varchar(255) not null DEFAULT '' comment '直播封面',
    `live_play` varchar(255) not null DEFAULT '' comment '播放地址,json格式?',
    `live_class` varchar(255) not null DEFAULT '' comment '平台直播类型',
    `live_type_id` int not null DEFAULT 0 comment '直播类型ID-自定义',
    `live_type_name` char(50) not null DEFAULT '' comment '直播类型NAME-自定义',
    `live_tag` varchar(50) not null DEFAULT '' comment '直播标签.多个 #号分割',
    `live_introduction` varchar(250) COLLATE utf8mb4_unicode_ci not null DEFAULT '' comment '直播间简介',
    `live_online_user` int not null DEFAULT 0 comment '直播间人数',
    `live_follow` int not null DEFAULT 0 comment '直播间关注人数',
    `live_uri` varchar(255) not null DEFAULT '' comment '直播间地址',
    `live_is_online` char(5) not null DEFAULT 'no' comment '直播间是否在播 ,yes|no',
    `live_platform` char(20) not null DEFAULT '' comment '所属平台',
    `live_status` char(10) not null DEFAULT 'display' comment '状态 ,隐藏:hide 显示:display',
	`live_play_time` int not null default 0 comment '最近开播时间',
	`live_play_end_time` int not null default 0 comment '最近关播时间',
    `created_at` int not null DEFAULT 0 comment '添加时间',
    `updated_at` int not null DEFAULT 0 comment '修改时间',
    `spider_pull_url` varchar(255) not null default '' comment '爬虫拉取URL',
    `platform_room_id` varchar(255) not null default '' comment '平台房间ID',
    key `live_anchortv_name` (live_anchortv_name),
    PRIMARY KEY (`live_id`),
    UNIQUE KEY `platform_room_id_live_platform` (`platform_room_id`,`live_platform`) USING BTREE
)engine=innodb default charset=utf8mb4 comment '直播间表';

alter table live add platform_room_id varchar(255) not null default '' comment '平台房间ID';
--
-- ALTER TABLE `live` MODIFY COLUMN `live_title` VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- ALTER TABLE `live` MODIFY COLUMN `live_anchortv_name` VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- ALTER TABLE `live` MODIFY COLUMN `live_introduction` VARCHAR(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
-- SET NAMES utf8mb4
-- ALTER TABLE `live` CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
--

-- 分类表
create table if not exists `live_type`(
  `type_id` int unsigned auto_increment primary key,
  `name` varchar(20) unique NOT NULL COMMENT '分类名',
  `key` varchar(20) unique COMMENT '分类标识',
  `parent_id` int  DEFAULT 0 comment '父级ID',
  `img` varchar(255) DEFAULT NULL comment '图片',
  `icon` varchar(255) DEFAULT NULL comment '图标',
  `order` tinyint(2) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '0' COMMENT '0-隐藏 1-显示 2-header',
  `title` varchar(50) DEFAULT NULL comment '网站title',
  `keywords` varchar(100) DEFAULT NULL comment '网站关键词',
  `description` varchar(255) DEFAULT NULL comment '网站说明',
  `subset` text DEFAULT NULL comment '分类子分类 关系映射 ##分割'
)engine=innodb default charset=utf8 comment '分类表';

-- 平台表
create table if not exists `live_platform`(
  `platform_id` int unsigned auto_increment primary key,
  `mark`  varchar(20) NOT NULL COMMENT '平台标识',
  `name`  varchar(50) NOT NULL COMMENT '平台名称',
  `domain`  varchar(255) NOT NULL COMMENT '平台域名',
  `pull_url`  varchar(255) NOT NULL COMMENT '抓取初始地址',
  `status` tinyint(1) DEFAULT 0 COMMENT '0-关闭 1-开启'
)engine=innodb default charset=utf8 comment '平台表';

-- 轮播表 , 首页,分类下轮播
create table if not exists `carousel`(
  `carousel_id` int unsigned auto_increment primary key,
  `associate_id` int default 0 comment '关联id',
  `associate_scene` char(50) default '' comment '关联场景',
  `jump_value` varchar(255) default '' comment '跳转值',
  `jump_method` char(50) default 'url' comment '跳转方式',
  `title` varchar(100) COLLATE utf8mb4_unicode_ci default '' comment '轮播标题',
  `imge` varchar(255) default '' comment '轮播图片',
  `display_time` datetime default null comment '显示时间',
  `hide_time` datetime default null comment '隐藏时间',
  `status` tinyint(1) default 1 comment '状态 0-隐藏,1-显示',
  `created_at` int DEFAULT NULL comment '添加时间',
  `updated_at` int DEFAULT NULL comment '修改时间',
  key `associate_id` (associate_id)
)engine=innodb default charset=utf8 comment '轮播表';






-- 记录golang坑 , 请忽略.
修改字符编码必须要修改mysql的配置文件my.cnf,然后重启才能生效

通常需要修改my.cnf的如下几个地方：

【client】 下面, 加上default-character-set=utf8mb4

【mysqld】 下面, 加上character_set_server = utf8mb4

重点:
    连接编码改为 utf8mb4






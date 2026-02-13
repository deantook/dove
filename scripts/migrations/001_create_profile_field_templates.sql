-- 创建系统资料字段模板表（存储系统预设的单个字段类型定义）
CREATE TABLE IF NOT EXISTS `profile_field_templates` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `field_key` VARCHAR(100) NOT NULL COMMENT '字段唯一标识（如：name, education, school）',
    `field_name` VARCHAR(100) NOT NULL COMMENT '字段显示名称（如：姓名、学历、毕业学校）',
    `field_type` VARCHAR(50) NOT NULL COMMENT '字段类型（TEXT_SINGLE, SELECT_SINGLE, NUMBER等）',
    `is_required` TINYINT(1) DEFAULT 0 COMMENT '是否必填',
    `is_searchable` TINYINT(1) DEFAULT 0 COMMENT '是否可搜索',
    `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开（默认解锁状态）',
    `default_value` TEXT COMMENT '默认值（JSON格式）',
    `options` TEXT COMMENT '选项配置（JSON格式，用于SELECT类型）',
    `validation` TEXT COMMENT '验证规则（JSON格式）',
    `display_order` INT DEFAULT 0 COMMENT '显示顺序',
    `icon` VARCHAR(500) COMMENT '图标URL',
    `description` VARCHAR(500) COMMENT '字段描述',
    `default_unlock_rules` TEXT COMMENT '默认解锁规则（JSON格式）',
    `category` VARCHAR(50) COMMENT '字段分类（如：基本信息、教育背景、工作经历等）',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_field_key` (`field_key`),
    INDEX `idx_field_type` (`field_type`),
    INDEX `idx_category` (`category`),
    INDEX `idx_is_active` (`is_active`),
    INDEX `idx_display_order` (`display_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统资料字段模板表';

-- 插入示例字段模板数据
INSERT INTO `profile_field_templates` (
    `field_key`,
    `field_name`,
    `field_type`,
    `is_required`,
    `is_searchable`,
    `is_public`,
    `display_order`,
    `description`,
    `category`,
    `default_unlock_rules`
) VALUES
-- 基本信息类
('nickname', '昵称', 'TEXT_SINGLE', 1, 1, 1, 1, '用户昵称', '基本信息', '{"unlock_type":"PUBLIC"}'),
('real_name', '真实姓名', 'TEXT_SINGLE', 0, 0, 0, 2, '真实姓名', '基本信息', '{"unlock_type":"CHAT","conditions":{"message_count":50}}'),
('gender', '性别', 'SELECT_SINGLE', 1, 1, 1, 3, '性别', '基本信息', '{"unlock_type":"PUBLIC"}'),
('age', '年龄', 'NUMBER', 1, 1, 0, 4, '年龄', '基本信息', '{"unlock_type":"CHAT","conditions":{"message_count":30}}'),
('birthday', '生日', 'DATE', 0, 0, 0, 5, '生日', '基本信息', '{"unlock_type":"TIME","conditions":{"friend_days":30}}'),
('avatar', '头像', 'IMAGE', 1, 0, 1, 6, '用户头像', '基本信息', '{"unlock_type":"PUBLIC"}'),

-- 教育背景类
('education', '学历', 'SELECT_SINGLE', 0, 1, 0, 10, '最高学历', '教育背景', '{"unlock_type":"CHAT","conditions":{"message_count":50}}'),
('school', '毕业学校', 'TEXT_SINGLE', 0, 1, 0, 11, '毕业院校', '教育背景', '{"unlock_type":"CHAT","conditions":{"message_count":80}}'),
('major', '专业', 'TEXT_SINGLE', 0, 1, 0, 12, '所学专业', '教育背景', '{"unlock_type":"CHAT","conditions":{"message_count":80}}'),
('graduation_year', '毕业年份', 'NUMBER', 0, 0, 0, 13, '毕业年份', '教育背景', '{"unlock_type":"TIME","conditions":{"friend_days":60}}'),

-- 工作信息类
('occupation', '职业', 'TEXT_SINGLE', 0, 1, 0, 20, '职业', '工作信息', '{"unlock_type":"CHAT","conditions":{"message_count":50}}'),
('company', '公司', 'TEXT_SINGLE', 0, 1, 0, 21, '所在公司', '工作信息', '{"unlock_type":"CHAT","conditions":{"message_count":100}}'),
('industry', '行业', 'SELECT_SINGLE', 0, 1, 0, 22, '所属行业', '工作信息', '{"unlock_type":"CHAT","conditions":{"message_count":70}}'),

-- 联系方式类
('phone', '手机号', 'TEXT_SINGLE', 0, 0, 0, 30, '手机号码', '联系方式', '{"unlock_type":"PAID","conditions":{"price":9.9}}'),
('wechat', '微信号', 'TEXT_SINGLE', 0, 0, 0, 31, '微信号', '联系方式', '{"unlock_type":"REQUEST","conditions":{"require_reason":true}}'),
('email', '邮箱', 'TEXT_SINGLE', 0, 0, 0, 32, '电子邮箱', '联系方式', '{"unlock_type":"TIME","conditions":{"friend_days":90}}'),

-- 个人介绍类
('bio', '个人简介', 'TEXT_MULTI', 0, 0, 0, 40, '个人简介', '个人介绍', '{"unlock_type":"CHAT","conditions":{"message_count":100}}'),
('hobbies', '兴趣爱好', 'TAG', 0, 1, 0, 41, '兴趣爱好标签', '个人介绍', '{"unlock_type":"TIME","conditions":{"friend_days":7}}'),
('photos', '照片', 'IMAGE', 0, 0, 0, 42, '个人照片', '个人介绍', '{"unlock_type":"COMBINED","conditions":{"logic":"OR","rules":[{"type":"CHAT","message_count":200},{"type":"TIME","friend_days":30}]}}'),
('video', '自我介绍视频', 'VIDEO', 0, 0, 0, 43, '自我介绍视频', '个人介绍', '{"unlock_type":"PAID","conditions":{"price":9.9}}'),

-- 位置信息类
('city', '所在城市', 'LOCATION', 1, 1, 0, 50, '当前所在城市', '位置信息', '{"unlock_type":"CHAT","conditions":{"message_count":50}}'),
('hometown', '家乡', 'LOCATION', 0, 1, 0, 51, '家乡', '位置信息', '{"unlock_type":"CHAT","conditions":{"message_count":100}}')
ON DUPLICATE KEY UPDATE `update_time` = CURRENT_TIMESTAMP;

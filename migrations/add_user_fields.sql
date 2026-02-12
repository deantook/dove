-- 扩展用户表字段
-- 表名: u_user

-- 1. 添加手机号字段（唯一索引）
ALTER TABLE `u_user` 
ADD COLUMN `phone` VARCHAR(20) NULL COMMENT '手机号' AFTER `username`;

-- 为手机号字段添加唯一索引
ALTER TABLE `u_user` 
ADD UNIQUE INDEX `idx_phone` (`phone`);

-- 2. 添加昵称字段
ALTER TABLE `u_user` 
ADD COLUMN `nickname` VARCHAR(100) NULL COMMENT '昵称' AFTER `phone`;

-- 3. 添加头像字段
ALTER TABLE `u_user` 
ADD COLUMN `avatar` VARCHAR(500) NULL COMMENT '头像URL' AFTER `nickname`;

-- 4. 添加状态字段（默认值为1，表示启用）
ALTER TABLE `u_user` 
ADD COLUMN `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：0-禁用，1-启用' AFTER `avatar`;

-- 5. 添加更新时间字段
ALTER TABLE `u_user` 
ADD COLUMN `update_time` DATETIME NULL COMMENT '更新时间' AFTER `create_time`;

-- 为现有数据设置默认更新时间（可选，如果需要的话）
UPDATE `u_user` SET `update_time` = `create_time` WHERE `update_time` IS NULL;

-- 如果需要将手机号字段设置为必填（NOT NULL），可以执行以下语句（注意：需要先确保所有现有记录都有手机号）
-- ALTER TABLE `u_user` MODIFY COLUMN `phone` VARCHAR(20) NOT NULL COMMENT '手机号';

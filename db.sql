DROP DATABASE IF EXISTS records;
CREATE DATABASE records;
USE records;
-- 创建用户信息表
CREATE TABLE users(
    `id` BIGINT AUTO_INCREMENT NOT NULL,
    `name` VARCHAR(10) NOT NULL,
    `id_number` CHAR(18) NOT NULL,
    `password` VARCHAR(16) NOT NULL,
    `type` BOOLEAN NOT NULL,
    `private_key` VARCHAR(500) NOT NULL,
    `public_key` VARCHAR(200) NOT NULL,
    UNIQUE (`id_number`),
    PRIMARY KEY(`id`)
) CHARSET = UTF8MB4;
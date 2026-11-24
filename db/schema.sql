-- GDG 식당 웨이팅 서비스 스키마 (MySQL 8)
SET NAMES utf8mb4;
CREATE TABLE IF NOT EXISTS users (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL COMMENT 'bcrypt hash',
    phone_num  VARCHAR(20)  NOT NULL DEFAULT '',
    name       VARCHAR(100) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS restaurants (
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    address      VARCHAR(255) NOT NULL DEFAULT '',
    description  TEXT         NOT NULL,
    phone_num    VARCHAR(20)  NOT NULL DEFAULT '',
    waiting_time INT          NOT NULL DEFAULT 0 COMMENT '예상 대기 시간(분)',
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 사장님(user) - 식당 연결
CREATE TABLE IF NOT EXISTS user_rest (
    user_id    BIGINT   NOT NULL,
    rest_id    BIGINT   NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, rest_id),
    UNIQUE KEY uq_user_rest_rest (rest_id) COMMENT '식당 하나에 사장님 한 명',
    CONSTRAINT fk_user_rest_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_rest_rest FOREIGN KEY (rest_id) REFERENCES restaurants (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 웨이팅. is_entry: 0=대기, 1=입장, 2=취소
CREATE TABLE IF NOT EXISTS reservation (
    id                 BIGINT AUTO_INCREMENT PRIMARY KEY,
    rest_id            BIGINT      NOT NULL,
    count              INT         NOT NULL,
    customer_phone_num VARCHAR(20) NOT NULL,
    created_at         DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_entry           TINYINT     NOT NULL DEFAULT 0,
    KEY idx_reservation_rest_status (rest_id, is_entry),
    KEY idx_reservation_rest_phone (rest_id, customer_phone_num, is_entry),
    CONSTRAINT fk_reservation_rest FOREIGN KEY (rest_id) REFERENCES restaurants (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

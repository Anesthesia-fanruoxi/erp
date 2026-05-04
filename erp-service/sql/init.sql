-- =====================================================================
-- ERP 系统初始化脚本 (一键执行: 建库 + 建表 + 初始数据)
-- 命名规范: 所有系统表统一以 sys_ 开头
-- MySQL 8.0+
-- =====================================================================
DROP DATABASE erp;
CREATE DATABASE IF NOT EXISTS erp DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE erp;

-- 为方便重复执行,先清理旧表(按外键依赖从子到父)
DROP TABLE IF EXISTS sys_audit_log;
DROP TABLE IF EXISTS sys_role_menu;
DROP TABLE IF EXISTS sys_user_role;
DROP TABLE IF EXISTS sys_device_binding;
DROP TABLE IF EXISTS sys_menu;
DROP TABLE IF EXISTS sys_role;
DROP TABLE IF EXISTS sys_user;

-- ---------------------------------------------------------------------
-- 1. 用户表
-- ---------------------------------------------------------------------
CREATE TABLE sys_user (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE                                    COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL                                           COMMENT 'bcrypt加密密码',
    real_name     VARCHAR(50)  NOT NULL DEFAULT ''                                COMMENT '真实姓名',
    email         VARCHAR(100) NOT NULL DEFAULT ''                                COMMENT '邮箱',
    phone         VARCHAR(20)  NOT NULL DEFAULT ''                                COMMENT '手机号',
    status        TINYINT      NOT NULL DEFAULT 0                                 COMMENT '状态: 0-待审核 1-正常 2-禁用',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP                 COMMENT '创建时间',
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at    DATETIME     NULL     DEFAULT NULL                              COMMENT '软删除时间',
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统用户表';

-- ---------------------------------------------------------------------
-- 2. 角色表
-- ---------------------------------------------------------------------
CREATE TABLE sys_role (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL UNIQUE                     COMMENT '角色名',
    description VARCHAR(255) NOT NULL DEFAULT ''                 COMMENT '角色描述',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统角色表';

-- ---------------------------------------------------------------------
-- 3. 用户-角色关联表
-- ---------------------------------------------------------------------
CREATE TABLE sys_user_role (
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES sys_user(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES sys_role(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- ---------------------------------------------------------------------
-- 4. 菜单表
-- 每个菜单自带三个权限标识:
--   code      → 菜单可见权限 (查看)
--   code:r    → 只读接口权限 (GET)
--   code:w    → 操作接口权限 (POST/PUT/DELETE)
-- type: 1-菜单(有路由) 2-目录(无路由,仅分组)
-- ---------------------------------------------------------------------
CREATE TABLE sys_menu (
    id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code      VARCHAR(100) NOT NULL UNIQUE             COMMENT '权限标识基础码 (如 system:user)',
    name      VARCHAR(100) NOT NULL                    COMMENT '菜单名称',
    type      TINYINT      NOT NULL DEFAULT 1          COMMENT '类型: 1-菜单 2-目录',
    parent_id BIGINT UNSIGNED NOT NULL DEFAULT 0       COMMENT '父级ID (0为顶级)',
    path      VARCHAR(200) NOT NULL DEFAULT ''         COMMENT '前端路由路径',
    icon      VARCHAR(50)  NOT NULL DEFAULT ''         COMMENT '菜单图标',
    sort      INT          NOT NULL DEFAULT 0          COMMENT '排序序号',
    visible   TINYINT      NOT NULL DEFAULT 1          COMMENT '是否可见: 0-隐藏 1-显示',
    INDEX idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统菜单表';

-- ---------------------------------------------------------------------
-- 5. 角色-菜单权限关联表
-- perm_type: 1-查看(菜单可见) 2-只读(GET接口) 3-操作(写接口)
-- ---------------------------------------------------------------------
CREATE TABLE sys_role_menu (
    role_id   BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    menu_id   BIGINT UNSIGNED NOT NULL COMMENT '菜单ID',
    perm_type TINYINT NOT NULL DEFAULT 1 COMMENT '权限类型: 1-查看 2-只读 3-操作',
    PRIMARY KEY (role_id, menu_id, perm_type),
    FOREIGN KEY (role_id) REFERENCES sys_role(id) ON DELETE CASCADE,
    FOREIGN KEY (menu_id) REFERENCES sys_menu(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单权限关联表';

-- ---------------------------------------------------------------------
-- 6. 设备绑定表
-- ---------------------------------------------------------------------
CREATE TABLE sys_device_binding (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id      BIGINT UNSIGNED NOT NULL UNIQUE                    COMMENT '用户ID (一个用户只能绑定一台设备)',
    machine_hash VARCHAR(255) NOT NULL                              COMMENT '机器码加密哈希',
    device_name  VARCHAR(100) NOT NULL DEFAULT ''                   COMMENT '设备名称',
    bound_at     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP    COMMENT '绑定时间',
    FOREIGN KEY (user_id) REFERENCES sys_user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备绑定表';

-- ---------------------------------------------------------------------
-- 7. 操作审计日志表
-- ---------------------------------------------------------------------
CREATE TABLE sys_audit_log (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL DEFAULT 0                   COMMENT '操作人ID',
    user_name   VARCHAR(50)  NOT NULL DEFAULT ''                     COMMENT '操作人用户名',
    action      VARCHAR(100) NOT NULL DEFAULT ''                     COMMENT '操作描述',
    method      VARCHAR(10)  NOT NULL DEFAULT ''                     COMMENT 'HTTP方法',
    path        VARCHAR(200) NOT NULL DEFAULT ''                     COMMENT '请求路径',
    query       TEXT                                                 COMMENT 'URL查询参数',
    body        TEXT                                                 COMMENT '请求体(已脱敏)',
    status_code INT          NOT NULL DEFAULT 0                      COMMENT '业务响应码',
    ip          VARCHAR(50)  NOT NULL DEFAULT ''                     COMMENT '来源IP',
    duration    BIGINT       NOT NULL DEFAULT 0                      COMMENT '耗时(ms)',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP      COMMENT '操作时间',
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作审计日志表';

-- ---------------------------------------------------------------------
-- 8. 合同主表
-- ---------------------------------------------------------------------
CREATE TABLE biz_contract (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    project_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '项目名称',
    order_no     VARCHAR(100) NOT NULL DEFAULT '' COMMENT '订单号',
    order_date   VARCHAR(20)  NOT NULL DEFAULT '' COMMENT '订单日期',
    from_company VARCHAR(200) NOT NULL DEFAULT '' COMMENT '发货方',
    to_company   VARCHAR(200) NOT NULL DEFAULT '' COMMENT '收货方',
    buyer        VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '采购员',
    attn         VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '联系人',
    buyer_email  VARCHAR(100) NOT NULL DEFAULT '' COMMENT '采购员邮箱',
    buyer_tel    VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '采购员电话',
    attn_tel     VARCHAR(50)  NOT NULL DEFAULT '' COMMENT '联系人电话',
    total_amount VARCHAR(30)  NOT NULL DEFAULT '' COMMENT '合计金额',
    delivery_addr VARCHAR(500) NOT NULL DEFAULT '' COMMENT '收货地址',
    remark       VARCHAR(500) NOT NULL DEFAULT '' COMMENT '备注',
    created_by   BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建人ID',
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_order_no (order_no),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合同主表';

-- ---------------------------------------------------------------------
-- 9. 合同明细表
-- ---------------------------------------------------------------------
CREATE TABLE biz_contract_item (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    contract_id BIGINT UNSIGNED NOT NULL COMMENT '合同ID',
    seq         INT          NOT NULL DEFAULT 0   COMMENT '序号',
    name        VARCHAR(200) NOT NULL DEFAULT ''  COMMENT '材料/设备名称',
    spec        VARCHAR(200) NOT NULL DEFAULT ''  COMMENT '型号规格/技术参数',
    brand       VARCHAR(100) NOT NULL DEFAULT ''  COMMENT '品牌/特征',
    qty         VARCHAR(30)  NOT NULL DEFAULT ''  COMMENT '数量',
    unit        VARCHAR(20)  NOT NULL DEFAULT ''  COMMENT '单位',
    unit_price  VARCHAR(30)  NOT NULL DEFAULT ''  COMMENT '含税单价',
    amount      VARCHAR(30)  NOT NULL DEFAULT ''  COMMENT '单项合计',
    operator    VARCHAR(50)  NOT NULL DEFAULT ''  COMMENT '下单人',
    location    VARCHAR(100) NOT NULL DEFAULT ''  COMMENT '安装位置',
    remark      VARCHAR(200) NOT NULL DEFAULT ''  COMMENT '备注',
    INDEX idx_contract_id (contract_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合同明细表';

-- =====================================================================
-- 初始化数据
-- =====================================================================

-- 默认管理员 (密码: password, bcrypt hash)
INSERT INTO sys_user (username, password_hash, real_name, status) VALUES
('admin', '$2a$10$FLxj8iax3DaaP/cV7U8usOm/zAXxPQast9Qnx4MH3IhkX.HmsYQN.', '超级管理员', 1);

-- 默认角色
INSERT INTO sys_role (name, description) VALUES
('admin', '超级管理员'),
('user',  '普通用户');

-- 给 admin 用户分配 admin 角色
INSERT INTO sys_user_role (user_id, role_id) VALUES (1, 1);

-- ---------------------------------------------------------------------
-- 菜单数据
-- 结构: 顶级目录 → 二级菜单(叶子节点,有路由)
-- 权限标识规则:
--   code      → 菜单可见 (查看)
--   code:r    → GET 接口 (只读)
--   code:w    → POST/PUT/DELETE 接口 (操作)
-- ---------------------------------------------------------------------

-- 顶级目录 (type=2, 无路由)
INSERT INTO sys_menu (id, code, name, type, parent_id, path, icon, sort, visible) VALUES
(1, 'system',   '系统管理', 2, 0, '', 'Setting',       1, 1),
(2, 'audit',    '审计管理', 2, 0, '', 'DocumentCopy',  2, 1),
(3, 'business', '业务管理', 2, 0, '', 'ShoppingCart',  3, 1),
(4, 'workflow', '流程管理', 2, 0, '', 'Connection',    4, 1);

-- 系统管理下的菜单 (type=1, 有路由)
INSERT INTO sys_menu (id, code, name, type, parent_id, path, icon, sort, visible) VALUES
(10, 'system:user',  '用户管理', 1, 1, '/system/user',  'User',       1, 1),
(11, 'system:role',  '角色管理', 1, 1, '/system/role',  'UserFilled', 2, 1),
(12, 'system:menu',  '菜单管理', 1, 1, '/system/menu',  'Menu',       3, 1),
(13, 'system:audit', '注册审核', 1, 1, '/system/audit', 'Checked',    4, 1);

-- 审计管理下的菜单 (type=1, 有路由)
INSERT INTO sys_menu (id, code, name, type, parent_id, path, icon, sort, visible) VALUES
(20, 'audit:log', '操作日志', 1, 2, '/audit/log', 'List', 1, 1);

-- 业务管理下的菜单 (type=1, 有路由)
INSERT INTO sys_menu (id, code, name, type, parent_id, path, icon, sort, visible) VALUES
(30, 'business:sale',     '销售管理', 1, 3, '/business/sale',     'TrendCharts', 1, 1),
(31, 'business:purchase', '采购管理', 1, 3, '/business/purchase', 'Box',         2, 1),
(32, 'business:retail',   '零售管理', 1, 3, '/business/retail',   'Goods',       3, 1),
(33, 'business:stock',    '库存管理', 1, 3, '/business/stock',    'Grid',        4, 1),
(34, 'business:contract', '合同管理', 1, 3, '/business/contract', 'Tickets',     5, 1);

-- 流程管理下的菜单 (type=1, 有路由)
INSERT INTO sys_menu (id, code, name, type, parent_id, path, icon, sort, visible) VALUES
(40, 'workflow:purchase', '采购流程', 1, 4, '/workflow/purchase', 'DocumentChecked', 1, 1),
(41, 'workflow:sale',     '销售流程', 1, 4, '/workflow/sale',     'Promotion',       2, 1);

-- ---------------------------------------------------------------------
-- admin 角色拥有所有菜单的全部权限 (查看+只读+操作)
-- perm_type: 1-查看 2-只读 3-操作
-- ---------------------------------------------------------------------
INSERT INTO sys_role_menu (role_id, menu_id, perm_type)
SELECT 1, id, 1 FROM sys_menu  -- 查看
UNION ALL
SELECT 1, id, 2 FROM sys_menu  -- 只读
UNION ALL
SELECT 1, id, 3 FROM sys_menu; -- 操作

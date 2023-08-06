

# 一、表设计

## 1、服务

```bash
CREATE TABLE `service`
(
    `id`          int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '服务ID',
    `name`        varchar(255)     NOT NULL DEFAULT '' COMMENT '服务名',
    `description` varchar(2000)    NOT NULL DEFAULT '' COMMENT '描述',
    `protocol`    varchar(255)     NOT NULL DEFAULT '' COMMENT '协议 http、grpc、double',
    `timeout`     int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '超时时间(毫秒)',
    `discovery`   varchar(255)     NOT NULL DEFAULT '' COMMENT '服务发现',
    `service`     varchar(255)     NOT NULL DEFAULT '' COMMENT '所在注册中心的服务名',
    `nodes`       varchar(255)     NOT NULL DEFAULT '' COMMENT '匿名服务地址，可以填多个。格式：addr1 weight=num1;addr2 weight=num2; 案例：172.17.0.3:80;172.17.0.4:80 weight=100。addr可以填域名或者ip地址。weight可省略，默认为1。',
    `balance`     varchar(255)     NOT NULL DEFAULT '' COMMENT '负载均衡算法',
    `creator`     varchar(62)      NOT NULL DEFAULT '' COMMENT '创建人',
    `operator`    varchar(62)      NOT NULL DEFAULT '' COMMENT '更新人',
    `created_at`  timestamp(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at`  timestamp(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) DEFAULT CHARSET = utf8mb4 COMMENT '服务';
```



## 2、插件

```bash
CREATE TABLE `plugin`
(
    `id`           int(10) unsigned    NOT NULL AUTO_INCREMENT COMMENT '插件ID',
    `name`         varchar(62)         NOT NULL DEFAULT '' COMMENT '名称',
    `description`  varchar(2000)       NOT NULL DEFAULT '' COMMENT '描述',
    `config`       text COMMENT '配置(JSON Schema)',
    `sort`         int(10) UNSIGNED    NOT NULL DEFAULT '0' COMMENT '排序',
    `enabled`      tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否打开 0-关闭 1-打开',
    `creator`      varchar(62)         NOT NULL DEFAULT '' COMMENT '创建人',
    `operator`     varchar(62)         NOT NULL DEFAULT '' COMMENT '更新人',
    `created_at`   timestamp(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at`   timestamp(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`name`)
) DEFAULT CHARSET = utf8mb4 COMMENT '插件';
```





## 3、路由

```bash
CREATE TABLE `http_router`
(
    `id`                 int(10) unsigned    NOT NULL AUTO_INCREMENT COMMENT '路由ID',
    `name`               varchar(62)         NOT NULL DEFAULT '' COMMENT '名称',
    `description`        varchar(100)        NOT NULL DEFAULT '' COMMENT '描述',
    `host`               varchar(62)         NOT NULL DEFAULT '' COMMENT '请求host',
    `method`             varchar(100)        NOT NULL DEFAULT '' COMMENT '请求方式 GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS',
    `location`           varchar(62)         NOT NULL DEFAULT '' COMMENT '客户端请求路径',
    `rules`              varchar(2000)       NOT NULL DEFAULT '' COMMENT '参数类型 header、query、cookie',
    `retry`              tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '超时重试次数',
    `time_out`           tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '超时时间，当为0时不设置超时，单位：ms',
    `plugin_names`       varchar(1000)       NOT NULL DEFAULT '' COMMENT '插件列表',
    `plugin_template_id` int(10) UNSIGNED    NOT NULL DEFAULT '0' COMMENT '插件模版ID',
    `service_id`         int(10) UNSIGNED    NOT NULL DEFAULT '0' COMMENT '服务ID',
    `disable`            tinyint(3) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否禁用路由 0-生效 1-禁用',
    `creator`            varchar(62)         NOT NULL DEFAULT '' COMMENT '创建人',
    `operator`           varchar(62)         NOT NULL DEFAULT '' COMMENT '更新人',
    `created_at`         timestamp(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    `updated_at`         timestamp(3)        NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8mb4 COMMENT 'http路由';
```






package main

import (
	"fmt"
	"log"

	"erp-service/common"
	"erp-service/config"
	"erp-service/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	if err := config.Load("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	common.LogInfo("配置加载成功")

	// 2. 初始化MySQL
	if err := common.InitMySQL(); err != nil {
		log.Fatalf("初始化MySQL失败: %v", err)
	}
	common.LogInfo("MySQL连接成功")

	// 3. 初始化Redis
	if err := common.InitRedis(); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	common.LogInfo("Redis连接成功")

	// 4. 预热所有角色权限到 Redis
	if err := common.WarmUpAllRolePerms(); err != nil {
		log.Fatalf("角色权限预热失败: %v", err)
	}

	// 5. 关闭 Gin debug 输出，使用 release 模式
	gin.SetMode(gin.ReleaseMode)

	// 6. 创建 Engine（不使用 Default，避免引入 Gin 自带的 Logger/Recovery）
	r := gin.New()

	// 7. 注册全局中间件（使用自定义实现）
	r.Use(common.CORSMiddleware())
	r.Use(common.RecoveryMiddleware())
	r.Use(common.RequestLogMiddleware())

	// 8. 注册路由
	router.SetupRouter(r)

	// 9. 启动服务
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	common.LogInfo("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

# webook

## 技术栈
Gin + Gorm + Kafka + Mysql + Redis + MongoDB

## 项目描述

1. 用户社交博客论坛，采用Gin实现用户注册，登录和发帖功能，支持阅读、点赞和收藏
2. 通过使用JWT实现长短Token和OAuth2原理实现了注册手机号密码和微信扫码注册登录，提升安全性和用户体验
3. 手写 RBAC 用户角色权限控制
4. 使用Kafka消息队列改造了阅读计数功能，采用批量消费，提升性能和解藕，一定程度解决了消息积压的问题
   - 同步转异步 批量处理
   - 开启一个事务处理批次
5. 通过Redis的ZSet实现生成热榜和点赞排行榜，再使用分布式任务调度定期调度热榜数据，确保时效性
   - 本地缓存 + redis缓存 + mysql
   - 本地缓存 同步 给其他实例
   - 本地 或 redis 缓存预加载 id -> article
   - BFF(Web) 缓存前置(service 的本地缓存 同步给 BFF 上)
6. 将单体应用拆分为微服务，使用gRPC进行高效通信，采用不停机迁移策略确保高可用性
7. 实现服务注册，发现和负载均衡机制，提升系统的可扩展性和可靠性
8. 采用了Prometheus，Zipkin，Grafana进行监控和报警，提高系统可观测性
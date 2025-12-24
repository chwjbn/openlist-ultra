# openlist-ultra
### 项目说明
- A new AList fork from OpenList, extending a host of enhanced features.
- 增强版AList,代码拷贝自OpenList仓库
- prod目录下是最终docker-compose部署文件,调整好配置后docker-compose up -d 就可以跑起来
### 增强功能
- 新IP访问飞书Webhook通知(配置参考docker-compose.yml的NOTIFY_WEBHOOK_FEISHU)
- 本地存储媒体支持连接到jellyfin,解决部分格式媒体网页端无法播放问题(配置参考docker-compose.yml的JELLYFIN_API_KEY/JELLYFIN_API_HOST/JELLYFIN_API_ROOT_FOLDER)
- 高性能反向代理jellyfin所有api接口

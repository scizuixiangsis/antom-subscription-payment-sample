# Antom 订阅支付演示项目

[English](README.md) | 简体中文

## 项目简介

这是一个演示 Antom 订阅支付的示例项目，包含 Go 后端和 React 前端。

**技术栈：** Go 1.22.5 + React + Vite

## 快速启动

### 1. 配置密钥（必需）

编辑 `server.go` 文件（第 23-41 行）：

```go
const (
    ClientID = "your_client_id"              // 替换为您的客户端 ID
    AntomPublicKey = "antom_public_key"      // 替换为 Antom 公钥
    MerchantPrivateKey = "your_private_key"  // 替换为您的私钥
)
```

**获取密钥：** 访问 [Antom 控制台](https://dashboard.antom.com/global-payments/developers/quickStart)

### 2. 配置回调 URL（可选）

编辑 `server.go`（第 137-138 行），本地测试可保持默认：

```go
subscriptionCreateRequest.PaymentNotificationUrl = "http://www.yourNotifyUrl.com/subscriptions/receivePaymentNotify"
subscriptionCreateRequest.SubscriptionNotificationUrl = "http://www.yourNotifyUrl.com/subscriptions/receiveSubscriptionNotify"
```

💡 本地开发如需接收回调，可使用 [ngrok](https://ngrok.com/) 将服务暴露到公网。

### 3. 启动服务

```bash
# 启动后端（端口 8080）
go run server.go

# 新开终端，启动前端（端口 5173）
cd client
npm run dev
```

### 4. 访问应用

打开浏览器：**http://localhost:5173**

---

## 项目结构

```
antom-subscription-payment-sample/
├── server.go           # Go 后端服务器
├── go.mod             # Go 依赖管理
├── go.sum             # Go 依赖校验
├── README.md          # 英文文档
├── README_CN.md       # 中文文档
└── client/            # React 前端
    ├── src/
    │   ├── App.jsx    # 主应用组件
    │   └── main.jsx   # 入口文件
    ├── package.json
    └── vite.config.js
```

## API 接口

- `POST /subscriptions/create` - 创建订阅
- `POST /subscriptions/receivePaymentNotify` - 接收支付通知
- `POST /subscriptions/receiveSubscriptionNotify` - 接收订阅通知

## 支付流程

1. 用户选择套餐（Standard 或 Premium）
2. 点击支付按钮，前端调用 `/subscriptions/create`
3. 跳转到 Antom 支付页面
4. 完成支付后返回商户页面
5. 后端接收支付和订阅状态通知

## 演示套餐

| 套餐 | 价格 | 入驻时间 | 语言支持 | 响应时间 |
|------|------|----------|----------|----------|
| Standard | 688.80 HKD | 3-5 工作日 | 中文/英文 | < 48 小时 |
| Premium | 1088.80 HKD | 1-2 工作日 | 中文/英文/日文/韩文 | < 24 小时 |

## 回调地址配置参考

| 平台 | 产品 | 集成类型 | 重定向地址 |
|------|------|----------|-----------|
| Web | 订阅支付 | API | subscriptionRedirectUrl: "http://localhost:5173/index.html?subscriptionRequestId=" + subscriptionRequestId |
| Android/iOS | 订阅支付 | API | subscriptionRedirectUrl: "subscription://app?subscriptionRequestId=" + subscriptionRequestId |

更多配置信息请参考 [Antom 文档](https://docs.antom.com/)

## 常见问题

### 没有 Antom 账号可以运行吗？
可以启动查看界面，但支付功能需要配置真实密钥。

### 必须配置哪些参数？
- ✅ 必需：`ClientID`、`AntomPublicKey`、`MerchantPrivateKey`
- ⭕ 可选：`PaymentNotificationUrl`、`SubscriptionNotificationUrl`（本地测试可不改）

### 点击支付按钮报错？
1. 检查密钥是否正确配置
2. 确认后端服务运行在 http://localhost:8080
3. 查看浏览器控制台和后端终端日志

### 如何接收支付回调？
本地开发使用 ngrok：
```bash
ngrok http 8080
# 将生成的公网 URL 配置到 server.go 的回调地址中
```

### 如何修改订阅金额？
- 前端：修改 `client/src/App.jsx` 中的套餐信息
- 后端：修改 `server.go` 的 `handleSubscriptionCreate` 函数（第 192-196 行）

## 注意事项

⚠️ 本项目仅用于沙箱测试  
⚠️ 请妥善保管私钥，不要提交到代码仓库  
⚠️ 生产环境需配置真实的回调 URL  

## 相关链接

- [Antom 官方文档](https://docs.antom.com/)
- [Antom 控制台](https://dashboard.antom.com/)
- [API 密钥配置](https://docs.antom.com/ac/ref/key_config)
- [测试资源](https://global.alipay.com/docs/ac/cashierpay/test)

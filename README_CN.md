# Antom 在线支付演示项目

[English](README.md) | 简体中文

## 项目简介

这是一个演示 Antom 订阅支付如何工作的示例项目。更多详细信息，请访问 [Antom 文档](https://docs.antom.com/)

## 技术栈

**后端：**
- Go 1.22.5 或更高版本
- Alipay Global Open SDK

**前端：**
- React
- Vite

## 环境要求

-   Go 1.22.5 或更高版本
-   Node.js 和 npm

## 配置说明

**<font color="#5A5AAD">⚠️ 重要：请务必替换以下参数为您自己的参数，以确保正确集成。</font>**

### 1. 配置支付密钥

在 `server.go` 文件中替换以下参数：

```go
const (
    ClientID = "your_client_id"              // 替换为您的客户端 ID
    AntomPublicKey = "antom_public_key"      // 替换为 Antom 公钥
    MerchantPrivateKey = "your_private_key"  // 替换为您的私钥
)
```

**参数说明：**
-  `CLIENT_ID`: 客户端 ID
-  `ANTOM_PUBLIC_KEY`: Antom 公钥（用于验证签名）
-  `MERCHANT_PRIVATE_KEY`: 您的私钥（用于签名）
   - ⚠️ Antom 不会保存您的私钥，因此您无法在控制台看到它
   - 请参考 [API 密钥配置文档](https://docs.antom.com/ac/ref/key_config)

**获取密钥：**

您可以在 [Antom 控制台](https://dashboard.antom.com/global-payments/developers/quickStart) 找到 `CLIENT_ID` 和 `ANTOM_PUBLIC_KEY`

### 2. 配置回调地址

按照下表信息在 `server.go` 中配置回调地址信息，以确保程序正常工作。

| 平台        | 产品                 | 集成类型          | 重定向地址                                                                                                      |
|-------------|----------------------|-------------------|------------------------------------------------------------------------------------------------------------|
| Web         | 收银台支付           | API/SDK/CKP       | paymentRedirectUrl: "http://localhost:5173/index.html?paymentRequestId=" + paymentRequestId                |
|             | 自动扣款             | API               | authRedirectUrl: "http://localhost:5173/receiveAuthCode.html"                                              |
|             |                      | SDK               | authRedirectUrl: "http://localhost:5173/receiveAuthCode"                                                   |
|             | 安心付               | SDK               | paymentRedirectUrl: "http://localhost:5173/index.html?paymentRequestId=" + paymentRequestId                |              
|             | 订阅支付             | API               | subscriptionRedirectUrl: "http://localhost:5173/index.html?subscriptionRequestId=" + subscriptionRequestId |
| Android/iOS | 收银台支付           | API               | paymentRedirectUrl: "cashierapi://app?paymentRequestId=" + paymentRequestId                                |
|             |                      | SDK               | paymentRedirectUrl: "cashiersdk://app?paymentRequestId=" + paymentRequestId                                |
|             | 自动扣款             | API               | authRedirectUrl: "autodebitapi://app/receiveAuthCode"                                                      |
|             |                      | SDK               | authRedirectUrl: "autodebitsdk://app/receiveAuthCode"                                                      |
|             | 安心付               | SDK               | paymentRedirectUrl: "easysafepay://app?paymentRequestId=" + paymentRequestId                               |              
|             | 订阅支付             | API               | subscriptionRedirectUrl: "subscription://app?subscriptionRequestId=" + subscriptionRequestId               |

## 安装和运行

### 1. 安装后端依赖

在项目根目录下运行：

```bash
go mod tidy
```

这将下载所有必需的 Go 依赖包。

### 2. 启动后端服务器

```bash
go run server.go
```

✅ 后端服务器将在 **http://localhost:8080** 启动

**后端提供的 API 接口：**
- `POST /subscriptions/create` - 创建订阅
- `POST /subscriptions/receivePaymentNotify` - 接收支付通知
- `POST /subscriptions/receiveSubscriptionNotify` - 接收订阅通知

### 3. 安装前端依赖

在新的终端窗口中，进入 client 目录：

```bash
cd client
npm install
```

### 4. 启动前端开发服务器

```bash
npm run dev
```

✅ 前端将在 **http://localhost:5173** 启动

### 5. 访问应用

打开浏览器访问：**http://localhost:5173**

## 项目结构

```
antom-subscription-payment-sample/
├── server.go           # Go 后端服务器
├── go.mod             # Go 依赖管理
├── go.sum             # Go 依赖校验
├── README.md          # 英文说明文档
├── README_CN.md       # 中文说明文档
└── client/            # React 前端
    ├── src/
    │   ├── App.jsx    # 主应用组件
    │   └── main.jsx   # 入口文件
    ├── index.html
    ├── package.json
    └── vite.config.js
```

## 功能说明

### 订阅支付流程

1. **用户选择套餐**：在前端页面选择 Standard 或 Premium 套餐
2. **创建订阅**：点击支付按钮，前端调用后端 `/subscriptions/create` 接口
3. **跳转支付**：后端返回支付链接，用户跳转到 Antom 支付页面
4. **完成支付**：用户在 Antom 页面完成支付
5. **接收通知**：后端接收支付和订阅状态通知
6. **返回结果**：用户被重定向回商户页面，显示支付结果

### 演示套餐

**Standard 套餐**
- 价格：688.80 HKD
- 入驻时间：3-5 个工作日
- 语言支持：中文/英文
- 咨询响应时间：< 48 小时

**Premium 套餐**
- 价格：1088.80 HKD
- 入驻时间：1-2 个工作日
- 语言支持：中文/英文/日文/韩文
- 咨询响应时间：< 24 小时

## 注意事项

1. **测试环境**：此项目仅用于沙箱测试，更多测试信息请参考 [测试资源](https://global.alipay.com/docs/ac/cashierpay/test)
2. **密钥安全**：请妥善保管您的私钥，不要提交到代码仓库
3. **回调地址**：生产环境需要配置真实的回调 URL
4. **端口占用**：确保 8080 和 5173 端口未被占用

## 常见问题

### Q: Go 依赖下载失败？
A: 尝试设置 Go 代理：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### Q: 前端启动失败？
A: 确保已安装 Node.js，并在 client 目录下运行 `npm install`

### Q: 支付回调收不到？
A: 本地开发环境需要使用内网穿透工具（如 ngrok）将本地服务暴露到公网

## 相关链接

- [Antom 官方文档](https://docs.antom.com/)
- [Antom 控制台](https://dashboard.antom.com/)
- [API 密钥配置](https://docs.antom.com/ac/ref/key_config)
- [测试资源](https://global.alipay.com/docs/ac/cashierpay/test)

## 许可证

本项目仅供学习和测试使用。


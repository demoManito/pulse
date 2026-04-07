package callback

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/demoManito/pulse/internal/biz/doc"
	"github.com/demoManito/pulse/internal/service"
	"github.com/demoManito/pulse/pkg/handler"
	"github.com/demoManito/pulse/pkg/logger"
)

// Handler return callback controller
func Handler() handler.Handler {
	return handler.Handler{
		Name: "callback",
		Actions: map[string]*handler.Action{
			"": handler.NewAction(handler.GET|handler.POST, ActionCallback),
		},
	}
}

// ActionCallback 处理企业微信回调
// GET  — URL 验证，返回解密后的 echostr
// POST — 接收事件消息，解密后分发处理
func ActionCallback(ctx *handler.Context) (handler.ActionResponse, error) {
	c := ctx.C

	msgSignature := c.Query("msg_signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")

	if c.Request.Method == http.MethodGet {
		return handleVerify(c, msgSignature, timestamp, nonce)
	}
	return handleEvent(ctx, c, msgSignature, timestamp, nonce)
}

// handleVerify 处理 GET 请求的 URL 验证
func handleVerify(c *gin.Context, msgSignature, timestamp, nonce string) (handler.ActionResponse, error) {
	echostr := c.Query("echostr")

	plaintext, err := service.CallbackCrypto.VerifyURL(msgSignature, timestamp, nonce, echostr)
	if err != nil {
		logger.Errorf("callback: URL 验证失败: %v", err)
		return nil, handler.NewActionError(http.StatusForbidden, -1, "verify failed")
	}

	logger.Info("callback: URL 验证成功")
	return handler.NewActionResponseRaw(func(c *gin.Context) {
		c.String(http.StatusOK, string(plaintext))
	}), nil
}

// handleEvent 处理 POST 请求的事件消息
func handleEvent(ctx *handler.Context, c *gin.Context, msgSignature, timestamp, nonce string) (handler.ActionResponse, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("callback: 读取请求体失败: %v", err)
		return nil, handler.NewActionError(http.StatusBadRequest, -1, "read body failed")
	}

	plaintext, err := service.CallbackCrypto.DecryptMsg(msgSignature, timestamp, nonce, body)
	if err != nil {
		logger.Errorf("callback: 解密消息失败: %v", err)
		return nil, handler.NewActionError(http.StatusForbidden, -1, "decrypt failed")
	}

	logger.Infof("callback: 收到事件消息: %s", string(plaintext))

	if err := doc.HandleEvent(ctx, plaintext); err != nil {
		logger.Errorf("callback: 处理事件失败: %v", err)
	}

	return handler.NewActionResponseWithMessage(nil, "success"), nil
}

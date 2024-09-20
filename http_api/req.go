package http_api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"time"
)

func SendReq(url string, req, data interface{}) error {
	var resp ApiResp
	resp.Data = &data

	res, _, errs := gorequest.New().Post(url).Retry(3, time.Second*5).
		Timeout(time.Second * 10).SendStruct(&req).EndStruct(&resp)

	if len(errs) > 0 {
		return fmt.Errorf("SendReq errs: %v", errs)
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("SendReq StatusCode: %d", res.StatusCode)
	}

	if resp.ErrNo != ApiCodeSuccess {
		return fmt.Errorf("%d - %s", resp.ErrNo, resp.ErrMsg)
	}
	return nil
}

func SendReqV2(url string, req, data interface{}) (*ApiResp, error) {
	var resp ApiResp
	resp.Data = &data

	res, _, errs := gorequest.New().Post(url).Retry(3, time.Second*5).
		Timeout(time.Second * 10).SendStruct(&req).EndStruct(&resp)

	if len(errs) > 0 {
		return nil, fmt.Errorf("SendReq errs: %v", errs)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SendReq StatusCode: %d", res.StatusCode)
	}

	return &resp, nil
}

func ReqIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		requestID := ctx.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			ctx.Request.Header.Set("X-Request-ID", requestID)
		}
		ctx.Writer.Header().Set("X-Request-ID", requestID)
		//c.Next()

		//requestId := ctx.GetHeader("Request-Id")

		c := context.WithValue(ctx.Request.Context(), "request_id", requestID)
		c1 := context.WithValue(c, "user_ip", ctx.ClientIP())
		c2 := context.WithValue(c1, "user_agent", ctx.GetHeader("User-Agent"))
		ctx.Request = ctx.Request.WithContext(c2)

		ctx.Next()
	}
}

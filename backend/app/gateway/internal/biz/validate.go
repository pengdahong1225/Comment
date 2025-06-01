package biz

import (
	"Comment/app/gateway/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// 表单类型集
type formTyper interface {
	models.AddCommentForm
}

// 通用表单验证
func validate[T formTyper](ctx *gin.Context, form *T) (*T, bool) {
	if err := ctx.ShouldBindJSON(&form); err != nil {
		var errs validator.ValidationErrors
		ok := errors.As(err, &errs)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    models.Failed,
				"message": "表单验证错误",
			})
			return nil, false
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    models.Failed,
			"message": errs.Error(),
		})
		return nil, false
	}
	return form, true
}

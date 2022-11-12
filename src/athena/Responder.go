package athena

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"sync"
)

var ResponderList []Responder
var ResponsePool *sync.Pool

type Responder interface {
	RespondTo() gin.HandlerFunc
}

type Response struct {
	HttpCode int    `json:"code"`
	Message  string `json:"message"`
	Data     any    `json:"data"`
}

func init() {
	ResponderList = []Responder{
		new(StringResponder),
		new(ModelResponder),
		new(ModelsResponder),
		new(JsonResponder),
		new(HttpCodeResponder),
		new(CollectionResponder),
	}

	ResponsePool = &sync.Pool{New: func() any {
		return &Response{
			HttpCode: http.StatusOK,
			Message:  "success",
			Data:     nil,
		}
	}}
}

// GetResponse 从response池中拿出一个对象
func GetResponse() *Response {
	return ResponsePool.Get().(*Response)
}

// PutResponse 放回response对象池
func PutResponse(response *Response) {
	ResponsePool.Put(response)
}

func Convert(handler interface{}) gin.HandlerFunc {
	handlerRef := reflect.ValueOf(handler)
	for _, responder := range ResponderList {
		responderRef := reflect.ValueOf(responder).Elem()
		if handlerRef.Type().ConvertibleTo(responderRef.Type()) {
			responderRef.Set(handlerRef)
			return responderRef.Interface().(Responder).RespondTo()
		}
	}
	return nil
}

// Controller return Json

type Json map[string]any
type JsonResponder func(*gin.Context) *Json

func (this JsonResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := GetResponse()
		defer PutResponse(response)
		// get的对象可能是上一次回收的对象,需要重新赋值
		response.HttpCode = http.StatusOK
		response.Message = "success"
		response.Data = this(context)
		context.JSON(response.HttpCode, response)
	}
}

// Controller return Collection

type CollectionResponder func(ctx *gin.Context) *Collection

func (this CollectionResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := GetResponse()
		defer PutResponse(response)
		response.HttpCode = http.StatusOK
		response.Message = "success"
		response.Data = this(context)
		context.JSON(response.HttpCode, response)
	}
}

// Controller return with http code

type HttpCodeResponder func(ctx *gin.Context) (int, string, *Json)

func (this HttpCodeResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := GetResponse()
		defer PutResponse(response)
		response.HttpCode, response.Message, response.Data = this(context)
		context.JSON(response.HttpCode, response)
	}
}

// Controller return string

type StringResponder func(*gin.Context) string

func (this StringResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := GetResponse()
		defer PutResponse(response)
		response.HttpCode = http.StatusOK
		response.Message = this(context)
		response.Data = nil
		context.JSON(response.HttpCode, response)
	}
}

// Controller return model

type ModelResponder func(*gin.Context) Model

// RespondTo Deprecated
func (this ModelResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		response := GetResponse()
		defer PutResponse(response)
		response.Data = this(context)
		context.JSON(response.HttpCode, response)
	}
}

// Controller return models

type ModelsResponder func(*gin.Context) Models

// RespondTo Deprecated
func (this ModelsResponder) RespondTo() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Content-type", "application/json")
		context.Writer.WriteString(string(this(context)))
	}
}

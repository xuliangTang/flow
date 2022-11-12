package classes

import (
	"flow/src/athena"
	"github.com/gin-gonic/gin"
)

type IndexClass struct {
}

func NewIndexClass() *IndexClass {
	return &IndexClass{}
}

func (this *IndexClass) version(ctx *gin.Context) string {
	if e, _ := ctx.GetQuery("err"); e == "1" {
		panic("test error")
	}
	return "v0.0.1"
}

func (this *IndexClass) Build(athena *athena.Athena) {
	athena.Handle("GET", "/version", this.version)
}

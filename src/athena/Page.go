package athena

import (
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type Page struct {
	CurrentPage int    `json:"current_page"` // 当前页
	PerPage     int    `json:"per_page"`     // 每页条数
	TotalSize   int64  `json:"total_size"`   // 总条数
	TotalPage   int    `json:"total_page"`   // 总页数
	Order       string `json:"-"`
}

func NewPage(currentPage int, perPage int, order string) *Page {
	return &Page{CurrentPage: currentPage, PerPage: perPage, Order: order}
}

func NewPageWithCtx(ctx *gin.Context) *Page {
	getCurrentPage := ctx.Query("page")
	currentPage, err := strconv.Atoi(getCurrentPage)
	if err != nil {
		currentPage = 1
	}

	getPerPage := ctx.Query("per_page")
	perPage, err := strconv.Atoi(getPerPage)
	if err != nil {
		perPage = 20
	}

	getOrder := ctx.Query("order")
	if getOrder == "" {
		getOrder = "id DESC"
	}

	return &Page{CurrentPage: currentPage, PerPage: perPage, Order: getOrder}
}

func (this *Page) IsValid() bool {
	return this.CurrentPage > 0 && this.PerPage > 0
}

func (this *Page) GetOffset() int {
	return (this.CurrentPage - 1) * this.PerPage
}

func (this *Page) SetTotal(totalSize int64) {
	this.TotalSize = totalSize
	this.TotalPage = int(math.Ceil(float64(this.TotalSize) / float64(this.PerPage)))
}

type Collection struct {
	Items any   `json:"items"`
	Page  *Page `json:"page"`
}

func NewCollection(items any, page *Page) *Collection {
	return &Collection{Items: items, Page: page}
}

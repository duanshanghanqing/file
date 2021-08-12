package dto

type PageBean struct {
	Current   int           `json:"current"`   // 第几页
	PageSize  int           `json:"pageSize"`  // 每页条数
	SortField string        `json:"sortField"` // 排序字段
	Sort      string        `json:"sort"`      // 排序方式
	Total     int           `json:"total"`     // 总条数
	TotalPage int           `json:"totalPage"` // 总页码 ,总条数 % 每页条数 == 0 ? 总条数/每页条数 : 总条数/每页条数 + 1
	List      []interface{} `json:"list"`
}

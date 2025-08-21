package models

// Pagination 分页信息结构体
type Pagination struct {
	CurrentPage  int
	TotalPages   int
	TotalRecords int
	PerPage      int
	HasNext      bool
	HasPrevious  bool
	NextPage     int
	PreviousPage int
}

// View 视图模型结构体
type View struct {
	Todos      []Todo
	Pagination Pagination
}

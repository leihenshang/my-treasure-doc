package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
)

type ListDocHistoryRequest struct {
	DocId int64 `json:"docId" binding:""`
	request.ListPagination
	request.ListSort
}

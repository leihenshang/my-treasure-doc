package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
)

type ListDocHistoryRequest struct {
	DocId string `json:"docId" form:"docId" binding:""`
	request.Pagination
	request.Sort
}

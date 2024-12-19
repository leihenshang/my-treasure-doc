package doc

import (
	"fastduck/treasure-doc/service/user/data/request"
	"fastduck/treasure-doc/service/user/gid"
)

type ListDocHistoryRequest struct {
	DocId gid.Gid `json:"docId" form:"docId" binding:""`
	request.Pagination
	request.Sort
}

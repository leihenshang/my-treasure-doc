package response

import "fastduck/treasure-doc/service/user/model"

type DocTree struct {
	*model.DocGroup
	Children []*model.DocGroup
}

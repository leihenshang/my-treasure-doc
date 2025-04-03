package response

import (
	"fastduck/treasure-doc/service/user/data/model"
)

type DocTree struct {
	*model.DocGroup
	Children []*model.DocGroup
}

type DocTrees []*DocTree

package response

import (
	"fastduck/treasure-doc/module/user/data/model"
)

type DocTree struct {
	*model.DocGroup
	Children []*model.DocGroup
}

type DocTrees []*DocTree

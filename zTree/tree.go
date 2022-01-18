package zTree

type Tree struct {
	Data     interface{} `json:"data"`
	Children []Tree      `json:"children"`
}

type INode interface {
	GetId() int
	GetFatherId() int
	GetData() interface{}
	IsRoot() bool
}
type INodes []INode

func GenerateTree(nodes []INode) (trees []Tree) {
	trees = []Tree{}
	var roots, childrenNodes []INode
	for _, v := range nodes {
		if v.IsRoot() {
			roots = append(roots, v)
			continue
		}
		childrenNodes = append(childrenNodes, v)
	}

	for _, v := range roots {
		childTree := &Tree{
			Data: v.GetData(),
		}

		recursiveTree(childTree, childrenNodes)

		trees = append(trees, *childTree)
	}
	return
}

func recursiveTree(tree *Tree, nodes []INode) {
	data := tree.Data.(INode)

	for _, v := range nodes {
		if v.IsRoot() {
			continue
		}
		if data.GetId() == v.GetFatherId() {
			childTree := &Tree{
				Data: v.GetData(),
			}
			recursiveTree(childTree, nodes)

			tree.Children = append(tree.Children, *childTree)
		}
	}
}

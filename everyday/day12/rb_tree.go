package rbtree

type RbValue interface {
	GetID() int
}

type RbTree struct {
	count int
	Root  *RbNode
}

type RbNode struct {
	Parent, Left, Right *RbNode
	IsBlack             bool
	Value               RbValue
}

func (t *RbNode) leftRotate(rbt *RbTree) {
	if t.Right == nil {
		return
	}
	y := t.Right
	if t.Parent != nil {
		if t.Parent.Left != nil && t.Parent.Left.Value.GetID() == t.Value.GetID() {
			t.Parent.Left = y
		} else {
			t.Parent.Right = y
		}
	} else {
		rbt.Root = y
	}
	y.Parent = t.Parent
	t.Parent = y
	t.Right = y.Left
	y.Left = t
	return
}

func (t *RbNode) rightRotate(rbt *RbTree) {
	if t.Left == nil {
		return
	}
	y := t.Left
	if t.Parent != nil {
		if t.Parent.Left != nil && t.Parent.Left.Value.GetID() == t.Value.GetID() {
			t.Parent.Left = y
		} else {
			t.Parent.Right = y
		}
	} else {
		rbt.Root = y
	}
	y.Parent = t.Parent
	t.Parent = y
	t.Left = y.Right
	y.Right = t
	return
}

func (t *RbNode) uncle() (bool, *RbNode) {
	// 父节点是红色，肯定存在祖父节点
	if t.Parent.Parent.Left != nil &&
		t.Parent.Parent.Left.Value.GetID() == t.Parent.Value.GetID() {
		return true, t.Parent.Parent.Right
	} else {
		return false, t.Parent.Parent.Left
	}
}

func (t *RbTree) Count() int {
	return t.count
}

func (t *RbTree) Depth() (depth int) {
	t.depth(t.Root, 1, &depth)
	return depth
}

func (t *RbTree) depth(node *RbNode, count int, depth *int) {
	if node == nil {
		return
	}
	if count > *depth {
		*depth = count
	}
	t.depth(node.Left, count+1, depth)
	t.depth(node.Right, count+1, depth)
	return
}

func (t *RbTree) Search(id int) (result RbValue) {
	node := t.Root
	for {
		if node == nil {
			return
		}
		if id > node.Value.GetID() {
			node = node.Right
		} else if id < node.Value.GetID() {
			node = node.Left
		} else {
			return node.Value
		}
	}
}

func (t *RbTree) List() []RbValue {
	var result = make([]RbValue, 0, t.count)
	t.list(t.Root, &result)
	return result
}

func (t *RbTree) list(node *RbNode, result *[]RbValue) {
	if node == nil {
		return
	}
	t.list(node.Left, result)
	*result = append(*result, node.Value)
	t.list(node.Right, result)
	return
}

func (t *RbTree) Check() bool {
	// 根节点
	if !t.Root.IsBlack {
		return false
	}
	// 颜色
	if !t.check(t.Root, t.Root.IsBlack) {
		return false
	}
	// 查重
	var m = make(map[int]struct{})
	for _, v := range t.List() {
		if _, ok := m[v.GetID()]; ok {
			return false
		}
		m[v.GetID()] = struct{}{}
	}
	return true
}

func (t *RbTree) check(node *RbNode, isBlack bool) bool {
	// 红色是否相连
	if node.Left != nil {
		if !isBlack && !node.Left.IsBlack {
			return false
		}
		if !t.check(node.Left, node.Left.IsBlack) {
			return false
		}
	}
	if node.Right != nil {
		if !isBlack && !node.Right.IsBlack {
			return false
		}
		if !t.check(node.Right, node.Right.IsBlack) {
			return false
		}
	}
	// 左右子节点黑色是否相同
	l, r := 0, 0
	t.checkLeft(node, &l)
	t.checkRight(node, &r)
	if l != r {
		return false
	}
	return true
}

func (t *RbTree) checkLeft(node *RbNode, count *int) {
	if node.Left != nil {
		if node.Left.IsBlack {
			*count++
		}
		t.checkLeft(node.Left, count)
	}
}

func (t *RbTree) checkRight(node *RbNode, count *int) {
	if node.Right != nil {
		if node.Right.IsBlack {
			*count++
		}
		t.checkRight(node.Right, count)
	}
}

func (t *RbTree) Insert(value RbValue) {
	node := &RbNode{
		Value: value,
	}
	if t.Root == nil {
		t.Root = node
	} else {
		parent := t.Root
		for {
			if node.Value.GetID() > parent.Value.GetID() {
				if parent.Right != nil {
					parent = parent.Right
				} else {
					parent.Right = node
					node.Parent = parent
					break
				}
			} else if node.Value.GetID() < parent.Value.GetID() {
				if parent.Left != nil {
					parent = parent.Left
				} else {
					parent.Left = node
					node.Parent = parent
					break
				}
			} else {
				parent.Value = value
				return
			}
		}
		// 循环调整
		// 父节点是红色才需要调整
		for node.Parent != nil && !node.Parent.IsBlack {
			// 获取叔叔节点
			parentIsLeftNode, uncle := node.uncle()
			if uncle != nil && !uncle.IsBlack {
				// 叔叔节点红色
				// 父节点和叔叔节点设置成黑色，祖父节点设置成红色，将祖父节点当做当前节点
				node.Parent.IsBlack = true
				uncle.IsBlack = true
				node.Parent.Parent.IsBlack = false
				node = node.Parent.Parent
			} else {
				// 叔叔节点是黑色
				if parentIsLeftNode {
					// 父节点是左节点
					if node.Parent.Right != nil &&
						node.Value.GetID() == node.Parent.Right.Value.GetID() {
						// 当前节点是右节点
						// 将父节点当作当前节点，然后左旋
						node = node.Parent
						node.leftRotate(t)
					}
					node.Parent.IsBlack = true
					node.Parent.Parent.IsBlack = false
					node.Parent.Parent.rightRotate(t)
				} else {
					// 父节点是右节点
					if node.Parent.Left != nil &&
						node.Value.GetID() == node.Parent.Left.Value.GetID() {
						node = node.Parent
						node.rightRotate(t)
					}
					node.Parent.IsBlack = true
					node.Parent.Parent.IsBlack = false
					node.Parent.Parent.leftRotate(t)
				}
			}
		}
	}
	t.Root.IsBlack = true
	t.count++
}

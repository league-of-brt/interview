package rbtree

import "fmt"

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

func (t *RbTree) leftRotate(x *RbNode) {
	if x.Right == nil {
		return
	}
	y := x.Right
	if x.Parent != nil {
		if x.Parent.Left == x {
			x.Parent.Left = y
		} else {
			x.Parent.Right = y
		}
	} else {
		t.Root = y
	}
	y.Parent = x.Parent
	x.Parent = y
	x.Right = y.Left
	if y.Left != nil {
		y.Left.Parent = x
	}
	y.Left = x
	return
}

func (t *RbTree) rightRotate(x *RbNode) {
	if x.Left == nil {
		return
	}
	y := x.Left
	if x.Parent != nil {
		if x.Parent.Right == x {
			x.Parent.Right = y
		} else {
			x.Parent.Left = y
		}
	} else {
		t.Root = y
	}
	y.Parent = x.Parent
	x.Parent = y
	x.Left = y.Right
	if y.Right != nil {
		y.Right.Parent = x
	}
	y.Right = x
	return
}

func (t *RbNode) uncle() (bool, *RbNode) {
	// 父节点是红色，肯定存在祖父节点
	if t.Parent.Parent.Left == t.Parent {
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
	node := t.search(id)
	if node != nil {
		return node.Value
	}
	return nil
}

func (t *RbTree) search(id int) (n *RbNode) {
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
			return node
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
					if node.Parent.Right == node {
						// 当前节点是右节点
						// 将父节点当作当前节点，然后左旋
						node = node.Parent
						t.leftRotate(node)
					}
					node.Parent.IsBlack = true
					node.Parent.Parent.IsBlack = false
					t.rightRotate(node.Parent.Parent)
				} else {
					// 父节点是右节点
					if node.Parent.Left == node {
						node = node.Parent
						t.rightRotate(node)
					}
					node.Parent.IsBlack = true
					node.Parent.Parent.IsBlack = false
					t.leftRotate(node.Parent.Parent)
				}
			}
		}
	}
	t.Root.IsBlack = true
	t.count++
}

func (t *RbTree) min(x *RbNode) *RbNode {
	if x == nil {
		return nil
	}
	for x.Left != nil {
		x = x.Left
	}
	return x
}

func (t *RbTree) max(x *RbNode) *RbNode {
	if x == nil {
		return nil
	}
	for x.Right != nil {
		x = x.Right
	}
	return x
}

func (t *RbTree) successor(x *RbNode) *RbNode {
	if x == nil {
		return nil
	}
	if x.Right != nil {
		return t.min(x.Right)
	}
	y := x.Parent
	for y != nil && x == y.Right {
		x = y
		y = y.Parent
	}
	return y
}

func (t *RbTree) Remove(id int) {
	z := t.search(id)
	if z == nil {
		return
	}
	var x, y *RbNode
	if z.Left == nil || z.Right == nil {
		// 没有子节点或者有一个子节点
		y = z
	} else {
		// 两个子节点都不为空，需要找到后继节点赋值给y
		// 在"被删除节点"有两个非空子节点的情况下，它的后继节点不可能是双子非空
		y = t.successor(z)
	}
	if y.Left != nil {
		x = y.Left
	} else {
		x = y.Right
	}
	// y没有子节点时，x为nil
	if x != nil {
		x.Parent = y.Parent
	}
	xIsLeft := true
	if y.Parent == nil {
		t.Root = x
	} else if y == y.Parent.Left {
		y.Parent.Left = x
	} else {
		y.Parent.Right = x
		xIsLeft = false
	}
	if y != z {
		z.Value = y.Value
	}
	if y.IsBlack {
		t.deleteFixup(x, y.Parent, xIsLeft)
	}
	t.count--
}

func (t *RbTree) deleteFixup(x, xp *RbNode, xIsLeft bool) {
	// x黑色
	for x != t.Root && (x == nil || x.IsBlack) {
		if xIsLeft {
			// x是左孩子
			w := xp.Right
			if !w.IsBlack {
				w.IsBlack = true
				xp.IsBlack = false
				t.leftRotate(xp)
				w = xp.Right
			}
			if (w.Left == nil || w.Left.IsBlack) &&
				(w.Right == nil || w.Right.IsBlack) {
				w.IsBlack = false
				x = xp
				xp = x.Parent
				if x == x.Parent.Left {
					xIsLeft = true
				} else {
					xIsLeft = false
				}
			} else {
				if w.Right == nil || w.Right.IsBlack {
					w.Left.IsBlack = true
					w.IsBlack = false
					t.rightRotate(w)
					w = xp.Right
				}
				w.IsBlack = xp.IsBlack
				xp.IsBlack = true
				w.Right.IsBlack = true
				t.leftRotate(xp)
				x = t.Root
			}
		} else {
			w := xp.Left
			if w != nil && !w.IsBlack {
				w.IsBlack = true
				xp.IsBlack = false
				t.rightRotate(xp)
				w = xp.Left
			}
			if (w.Left == nil || w.Left.IsBlack) &&
				(w.Right == nil || w.Right.IsBlack) {
				w.IsBlack = false
				x = xp
				xp = x.Parent
				if x == x.Parent.Left {
					xIsLeft = true
				} else {
					xIsLeft = false
				}
			} else {
				if w.Left == nil || w.Left.IsBlack {
					w.Right.IsBlack = true
					w.IsBlack = false
					t.leftRotate(w)
					w = xp.Left
				}
				w.IsBlack = xp.IsBlack
				xp.IsBlack = true
				w.Left.IsBlack = true
				t.rightRotate(xp)
				x = t.Root
			}
		}
	}
	x.IsBlack = true
}

// all-树中所有数据的集合
func Check(t RbTree, all map[int]struct{}) error {
	// 根节点
	if !t.Root.IsBlack {
		return fmt.Errorf("root red")
	}
	// 颜色
	err := check(t.Root)
	if err != nil {
		return err
	}
	var treeAll = make(map[int]struct{})
	for _, v := range t.List() {
		id := v.GetID()
		// 查重
		if _, ok := treeAll[id]; ok {
			return fmt.Errorf("id repeat %d", id)
		}
		treeAll[id] = struct{}{}
		// 多余数据
		if _, ok := all[id]; !ok {
			return fmt.Errorf("id not exit %d", id)
		}
	}
	// 数据丢失
	for id := range all {
		if _, ok := treeAll[id]; !ok {
			return fmt.Errorf("id lost %d", id)
		}
	}
	return nil
}

func check(node *RbNode) error {
	// 红色是否相连
	if node.Left != nil {
		if !node.IsBlack && !node.Left.IsBlack {
			return fmt.Errorf("red connect %d", node.Value.GetID())
		}
		err := check(node.Left)
		if err != nil {
			return err
		}
	}
	if node.Right != nil {
		if !node.IsBlack && !node.Right.IsBlack {
			return fmt.Errorf("red connect %d", node.Value.GetID())
		}
		err := check(node.Right)
		if err != nil {
			return err
		}
	}
	// 左右子节点黑色是否相同
	l, r := 0, 0
	checkLeft(node, &l)
	checkRight(node, &r)
	if l != r {
		return fmt.Errorf("black not balance Left:%d Right：%d", l, r)
	}
	return nil
}

func checkLeft(node *RbNode, count *int) {
	if node.Left != nil {
		if node.Left.IsBlack {
			*count++
		}
		checkLeft(node.Left, count)
	}
}

func checkRight(node *RbNode, count *int) {
	if node.Right != nil {
		if node.Right.IsBlack {
			*count++
		}
		checkRight(node.Right, count)
	}
}

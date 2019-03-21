package rbtree

//RBTree 红黑树本身为一个结构体，可以直接使用
type RBTree struct {
	root *node
}

//node 节点为一个结构体
type node struct {
	color      bool
	leftNode   *node
	rightNode  *node
	fatherNode *node
	value      int
}

//红黑通过bool类型来存储，并设置常量
const (
	RBTRed   = false
	RBTBlack = true
)

func (t *RBTree) leftRotate(n *node) {
	rn := n.rightNode
	//first give n's father to rn's father
	rn.fatherNode = n.fatherNode
	if n.fatherNode != nil {
		if n.fatherNode.leftNode == n {
			n.fatherNode.leftNode = rn
		} else {
			n.fatherNode.rightNode = rn
		}
	} else {
		t.root = rn
	}
	n.fatherNode = rn
	n.rightNode = rn.leftNode
	if n.rightNode != nil {
		n.rightNode.fatherNode = n
	}
	rn.leftNode = n
}
func (t *RBTree) rightRotate(n *node) {

	ln := n.leftNode
	ln.fatherNode = n.fatherNode
	if n.fatherNode != nil {
		if n.fatherNode.leftNode == n {
			n.fatherNode.leftNode = ln
		} else {
			n.fatherNode.rightNode = ln
		}
	} else {
		t.root = ln
	}
	n.fatherNode = ln

	n.leftNode = ln.rightNode
	if n.leftNode != nil {

		n.leftNode.fatherNode = n
	}
	ln.rightNode = n

}

func (t *RBTree) insert(v int) {
	if t.root == nil {
		t.root = &node{value: v, color: RBTBlack}
		return
	}
	n := t.root

	insertNode := &node{value: v, color: RBTRed}
	var nf *node
	for n != nil {
		//fmt.Println("condition", n)
		nf = n
		if v < n.value {
			n = n.leftNode
		} else if v > n.value {
			n = n.rightNode
		} else {
			//TODO fix the condition that replace value
			return
		}

	}
	//设置新插入节点的父节点
	insertNode.fatherNode = nf
	//将新的节点挂到父节点上
	if v < nf.value {
		nf.leftNode = insertNode
	} else {
		nf.rightNode = insertNode
	}
	t.insertFixUp(insertNode)
}
func (t *RBTree) insertFixUp(n *node) {

	for !isBlack(n.fatherNode) {
		//fmt.Printf("%v\t", n)
		//grandpa's color is black
		//case1  uncle's color is red then set grandpa's red color and his child black
		// n -> n's grandpa
		//if n is the same side of its father
		//exchange its father and grandpa by rotate
		//else make its side by rotate
		uncleNode := findBroNode(n.fatherNode)
		if !isBlack(uncleNode) {
			n.fatherNode.color = RBTBlack
			uncleNode.color = RBTBlack
			uncleNode.fatherNode.color = RBTRed
			n = n.fatherNode.fatherNode
			//	fmt.Println("condition1")
		} else if n.fatherNode == n.fatherNode.fatherNode.leftNode {
			//fmt.Println("condition2")
			if n == n.fatherNode.leftNode {

				n.fatherNode.fatherNode.color = RBTRed
				n.fatherNode.color = RBTBlack
				n = n.fatherNode.fatherNode
				t.rightRotate(n)

			} else {
				n = n.fatherNode
				t.leftRotate(n)
			}

		} else {
			//fmt.Println("condition2")

			if n == n.fatherNode.rightNode {
				n.fatherNode.fatherNode.color = RBTRed
				n.fatherNode.color = RBTBlack
				n = n.fatherNode.fatherNode
				t.leftRotate(n)

			} else {
				n = n.fatherNode
				t.rightRotate(n)
			}
		}
		t.root.color = RBTBlack
	}
}
func isBlack(n *node) bool {
	if n == nil {
		return true
	} else {
		return n.color == RBTBlack
	}
}
func setColor(n *node, color bool) {
	if n == nil {
		return
	}
	n.color = color
}
func findBroNode(n *node) (bro *node) {
	if n.fatherNode == nil {
		return nil
	}

	if n.fatherNode.leftNode == n {
		bro = n.fatherNode.rightNode
	} else if n.fatherNode.rightNode == n {
		bro = n.fatherNode.leftNode
	} else {
		if n.fatherNode.leftNode == nil {
			bro = n.fatherNode.rightNode
		} else {
			bro = n.fatherNode.leftNode

		}
	}
	return bro
}

func (t *RBTree) delete(v int) {
	n := t.find(v)
	if n == nil {
		return
	}
	// if n == t.root {
	// 	fmt.Println("delete root")
	// 	t.printGra()
	// }
	//copy color of fixNode
	var fixColor = n.color
	//if fixNode == nil copy node of start fix node
	//set it's father and set color black
	var fixNode = &node{fatherNode: n.fatherNode, color: RBTBlack}

	if n.leftNode == nil {
		t.transplant(n, n.rightNode)
		if n.rightNode != nil {
			fixNode = n.rightNode
		}
	} else if n.rightNode == nil {
		t.transplant(n, n.leftNode)
		if n.leftNode != nil {
			fixNode = n.leftNode
		}
	} else {
		succNode := t.miniNum(n.rightNode)
		fixColor = succNode.color
		if succNode.rightNode == nil {
			if succNode.fatherNode != n {
				fixNode = &node{fatherNode: succNode.fatherNode, color: RBTBlack}
			} else {
				fixNode = &node{fatherNode: succNode, color: RBTBlack}
			}
		} else {
			fixNode = succNode.rightNode
		}
		if succNode.fatherNode != n {
			t.transplant(succNode, succNode.rightNode)
			succNode.rightNode = n.rightNode
			succNode.rightNode.fatherNode = succNode
		} else {

		}
		t.transplant(n, succNode)
		succNode.leftNode = n.leftNode
		succNode.leftNode.fatherNode = succNode
		succNode.color = n.color
	}

	if fixColor == RBTBlack {
		t.deleteFixUp(fixNode)
	}

}
func (t *RBTree) deleteFixUp(n *node) {

	// fmt.Printf("%v,%v\n", n, t.root)
	// fmt.Printf("\n")

	if t.root == nil {
		return
	}
	for n != t.root && isBlack(n) {
		bro := findBroNode(n)

		//fmt.Printf("%v\n", n.fatherNode)
		/**fmt.Printf("%v\n", bro.fatherNode)
		fmt.Printf("\n")
		*/
		if bro != n.fatherNode.leftNode {
			//case 1
			if !isBlack(bro) {

				n.fatherNode.color = RBTRed
				bro.color = RBTBlack
				t.leftRotate(n.fatherNode)
				bro = findBroNode(n)
				// now new bro is black
			}

			//if bro is black its children perhaps be nil
			//if bro's children are black
			// n up
			if isBlack(bro.leftNode) && isBlack(bro.rightNode) {
				//case2
				setColor(bro, RBTRed)
				n = n.fatherNode
				//fmt.Printf("%v\n", n)

			} else {
				//case3
				if !isBlack(bro.rightNode) {

					bro.color = n.fatherNode.color
					bro.rightNode.color = RBTBlack
					n.fatherNode.color = RBTBlack
					t.leftRotate(n.fatherNode)
					n = t.root
				} else {
					//case4

					bro.color = RBTRed
					bro.leftNode.color = RBTBlack
					t.rightRotate(bro)
					bro = findBroNode(n)
				}
			}

		} else {
			//case 1
			if !isBlack(bro) {

				n.fatherNode.color = RBTRed
				bro.color = RBTBlack
				t.rightRotate(n.fatherNode)
				bro = findBroNode(n)
				// now new bro is black
			}

			//if bro is black its children perhaps be nil
			//if bro's children are black
			// n up

			if isBlack(bro.leftNode) && isBlack(bro.rightNode) {
				//case2

				setColor(bro, RBTRed)
				n = n.fatherNode
			} else {
				//case3
				if !isBlack(bro.leftNode) {

					bro.color = n.fatherNode.color
					bro.leftNode.color = RBTBlack
					n.fatherNode.color = RBTBlack
					t.rightRotate(n.fatherNode)
					break

				} else {
					//case4

					bro.color = RBTRed
					bro.rightNode.color = RBTBlack
					t.leftRotate(bro)
				}
			}

		}
	}
	n.color = RBTBlack
}
func (t *RBTree) miniNum(n *node) *node {
	for n.leftNode != nil {
		n = n.leftNode
	}
	return n
}
func (t *RBTree) transplant(u, v *node) {

	if u.fatherNode == nil {
		t.root = v
		if v != nil {
			v.fatherNode = nil
		}

	} else if u == u.fatherNode.leftNode {
		u.fatherNode.leftNode = v

	} else {
		u.fatherNode.rightNode = v

	}
	if v != nil {
		v.fatherNode = u.fatherNode
	}

}

//find 查找
func (t *RBTree) find(v int) *node {
	n := t.root
	for n != nil {
		if v < n.value {
			//小于当前节点的话，往左节点找
			n = n.leftNode
		} else if v > n.value {
			//大于当前节点的话，往右节点找
			n = n.rightNode
		} else {
			//等于的话表示找到，返回
			return n
		}
	}
	//循环结束没找到，返回
	return nil
}

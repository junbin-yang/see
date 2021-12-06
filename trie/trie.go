package trie

import (
	"github.com/junbin-yang/golib/bytesconv"
	"net/url"
	"strings"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

type nodeType uint8

const (
	root nodeType = iota + 1
	param
	catchAll
)

type Node struct {
	//表示当前节点的path
	path string

	// children列表的path的各首字符组成的string
	indices string

	//默认是false，当children是 通配符类型时，wildChild为true
	wildChild bool

	//节点的类型，默认是static类型
	nType nodeType

	//代表了有几条路由会经过此节点，用于在节点进行排序时使用
	priority uint32

	//子节点
	children []*Node

	// 是否是最后一个
	isEnd bool

	// 是从root节点到当前节点的全部path部分
	fullPath string
}

func NewTree() *Node {
	return &Node{fullPath: "/"}
}

// addChild will add a child node, keeping wildcards at the end
func (n *Node) addChild(child *Node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

// Increments priority of the given child and reorders if necessary
func (n *Node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// Build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // Unchanged prefix, might be empty
			n.indices[pos:pos+1] + // The index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // Rest without char at 'pos'
	}

	return newPos
}

func (n *Node) Insert(path string) {
	isEnd := true
	fullPath := path
	n.priority++

	// 如果单前树是空树，直接在当前node插入path
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, isEnd)
		n.nType = root
		return
	}

	//path的共同前缀位置
	parentFullPathIndex := 0

walk:
	for {
		// Find the longest common prefix.
		// This also implies that the common prefix contains no ':' or '*'
		// since the existing key can't contain those chars.
		i := longestCommonPrefix(path, n.path)

		// 如果path与当前的node有部分匹配，需要拆分当前的node
		if i < len(n.path) {
			child := Node{
				path:      n.path[i:],  //新的节点包括 node path 没有匹配上的后半部分
				wildChild: n.wildChild, //设置与当前节点相同
				indices:   n.indices,
				children:  n.children,
				isEnd:     n.isEnd,
				priority:  n.priority - 1,
				fullPath:  n.fullPath,
			}

			//将后半部分设置为孩子节点
			n.children = []*Node{&child}
			n.indices = bytesconv.BytesToString([]byte{n.path[i]})
			//当前节点的path只保持前半部分
			n.path = path[:i]
			n.isEnd = false
			//后半部分节点一定不包含通配符
			n.wildChild = false
			//当前节点的fullPath截取
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		//path没有完成匹配，需要继续向下寻找
		if i < len(path) {
			//path 更新为没有匹配上的后半部分
			path = path[i:]
			//当前节点的孩子节点不是通配符类型，取出第一个字符
			c := path[0]

			//冒号通配符后面的 下划线处理
			if n.nType == param && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				//更新node节点为孩子节点，继续查找
				n = n.children[0]
				n.priority++
				continue walk
			}

			//当前节点的某个孩子与path有相同的前缀
			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					//更新当前节点为对应的孩子节点，继续查找
					n = n.children[i]
					continue walk
				}
			}

			//如果是其他情况，新增一个child节点，并且基于这个child节点，插入剩下的path
			if c != ':' && c != '*' && n.nType != catchAll {
				// []byte for proper unicode char conversion, see #65
				n.indices += bytesconv.BytesToString([]byte{c})
				child := &Node{
					fullPath: fullPath,
				}
				n.addChild(child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			} else if n.wildChild {
				// inserting a wildcard node, need to check if it conflicts with the existing wildcard
				n = n.children[len(n.children)-1]
				n.priority++

				// Check if the wildcard matches
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
					// Adding a child to a catchAll is not possible
					n.nType != catchAll &&
					// Check for longer wildcard, e.g. :name and :names
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				}

				// Wildcard conflict
				pathSeg := path
				if n.nType != catchAll {
					pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				panic("'" + pathSeg +
					"' in new path '" + fullPath +
					"' conflicts with existing wildcard '" + n.path +
					"' in existing prefix '" + prefix +
					"'")
			}

			n.insertChild(path, fullPath, isEnd)
			return
		}

		// Otherwise add handle to current node
		if n.isEnd != false {
			panic("handlers are already registered for path '" + fullPath + "'")
		}
		n.isEnd = isEnd
		n.fullPath = fullPath
		return
	}
}

// Search for a wildcard segment and check the name for invalid characters.
// Returns -1 as index, if no wildcard was found.
func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		//如果没有遇到通配符就继续向后查找
		if c != ':' && c != '*' {
			continue
		}

		//找到通配符设置valid为true，那么通配符在path的起始位置就是start
		valid = true
		//从通配符后面继续查找
		for end, c := range []byte(path[start+1:]) {
			switch c {
			//如果遇到下划线，返回wildCard（不包括下划线）、start、true
			case '/':
				return path[start : start+1+end], start, valid
			//如果遇到通配符，valid设置为false
			case ':', '*':
				valid = false
			}
		}
		//在这个位置返回，遍历完了path，valid为true和false的可能性都有
		return path[start:], start, valid
	}
	//在path里没有找到通配符
	return "", -1, false
}

func (n *Node) insertChild(path string, fullPath string, isEnd bool) {
	for {
		// Find prefix until first wildcard
		wildcard, i, valid := findWildcard(path)
		//path中不包含通配符，直接结束对numParams条件的for循环
		if i < 0 { // No wildcard found
			break
		}

		// valid为false的两种情况是没有找到通配符（之前已经break） 或者 一个path段有多个通配符
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// 如果path段只有通配符没有名字 也会panic；由于wildCard一定是以通配符开头的，通配符后面不能直接加下划线
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		//冒号类型的通配符处理
		if wildcard[0] == ':' { // param
			if i > 0 {
				//设置当前节点的path
				n.path = path[:i]
				//更新path
				path = path[i:]
			}

			child := &Node{
				nType:    param,    //冒号类型的通配符类型
				path:     wildcard, //设置path为wildCard path包含通配符和名字
				fullPath: fullPath,
			}
			n.addChild(child)
			//孩子节点是通配符，当前节点设置为true
			n.wildChild = true
			//n更新为下沉到孩子节点
			n = child
			n.priority++

			//如果wildCard的长度小于path，则说明path中还包含以及path
			if len(wildcard) < len(path) {
				//重新更新path
				path = path[len(wildcard):]

				//new一个child节点
				child := &Node{
					priority: 1,
					fullPath: fullPath,
				}
				n.addChild(child)
				//更新n节点为child节点
				n = child
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			n.isEnd = isEnd
			return
		}

		//星号通配符类型的处理
		//星号通配符必须是path的最后一个通配符 否则会panic
		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			pathSeg := strings.SplitN(n.children[0].path, "/", 2)[0]
			panic("catch-all wildcard '" + path +
				"' in new path '" + fullPath +
				"' conflicts with existing path segment '" + pathSeg +
				"' in existing prefix '" + n.path + pathSeg +
				"'")
		}

		//星号通配符的前一个字符，必须为下划线，否则panic
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		//当前node的path为星号通配符之前的path
		n.path = path[:i]

		//一个path为空的节点
		child := &Node{
			wildChild: true, //空节点的wildCard为true
			nType:     catchAll,
			fullPath:  fullPath,
		}

		n.addChild(child)
		//node节点 indices设置为 下划线
		n.indices = string('/')
		//node节点下沉为path为空的节点
		n = child
		n.priority++

		// second node: node holding the variable
		child = &Node{
			path:     path[i:], //path 为从 下划线开始的包含星号通配符的path
			nType:    catchAll,
			isEnd:    isEnd,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*Node{child} //将包含星号通配符的path节点挂接到空节点下

		return
	}

	// 剩下的path不再包含冒号或者星号
	n.path = path
	n.isEnd = isEnd
	n.fullPath = fullPath
}

func (n *Node) Search(path string, params *[]Param) (fullPath string) {
	var unescape = false
	var globalParamsCount int16

walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			// 有公共前缀
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// Try all the non-wildcard children first by matching the indices
				idxc := path[0]
				for i, c := range []byte(n.indices) {
					if c == idxc {
						n = n.children[i]
						continue walk
					}
				}

				if !n.wildChild {
					return
				}

				// Handle wildcard child, which is always at the end of the array
				n = n.children[len(n.children)-1]
				globalParamsCount++

				switch n.nType {
				case param:
					// fix truncate the parameter
					// tree_test.go  line: 204

					// Find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if params != nil && cap(*params) > 0 {
						if params == nil {
							params = params
						}
						// Expand slice within preallocated capacity
						i := len(*params)
						*params = (*params)[:i+1]
						val := path[:end]
						if unescape {
							if v, err := url.QueryUnescape(val); err == nil {
								val = v
							}
						}
						(*params)[i] = Param{
							Key:   n.path[1:],
							Value: val,
						}
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						return
					}

					if n.isEnd {
						fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						n = n.children[0]
					}
					return

				case catchAll:
					// Save param value
					if params != nil {
						if params == nil {
							params = params
						}
						// Expand slice within preallocated capacity
						i := len(*params)
						*params = (*params)[:i+1]
						val := path
						if unescape {
							if v, err := url.QueryUnescape(path); err == nil {
								val = v
							}
						}
						(*params)[i] = Param{
							Key:   n.path[2:],
							Value: val,
						}
					}

					fullPath = n.fullPath
					return

				default:
					panic("invalid node type")
				}
			}
		}

		if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.isEnd {
				fullPath = n.fullPath
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			if path == "/" && n.wildChild && n.nType != root {
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					return
				}
			}

			return
		}

		return
	}
}

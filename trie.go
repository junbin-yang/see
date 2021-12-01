package see

import (
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
	"strings"
)

type array []string

func (a *array) GetKey(value string) (int32, bool) {
	l := len(a)
	data := make([]int32, l)
	for i := int32(0); i < int32(l); i++ {
		data[i] = i
	}
	codec := encode.I32{}
	st, _ := trie.NewSlimTrie(codec, a, data)
	return st.GetI32(value)
}

// 前缀树实现
type node struct {
	pattern  string  // 是否一个完整的url，不是则为空字符串
	part     string  // URL块值，用/分割的部分，比如/abc/123，abc和123就是2个part
	children []*node // 节点下的子节点
	isWild   bool    // 是否模糊匹配，比如:id*id这样的node就为true
}

/*
 * 对于路由来说，最重要的当然是注册与匹配了。开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler。
 * 因此，前缀树需要支持节点的插入与查询。插入功能很简单，递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个；
 * 有一点需要注意，/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。
 * p和:lang节点的pattern属性皆为空。因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。
 * 例如，/p/python虽能成功匹配到:lang，但:lang的pattern值为空，因此匹配失败。
 * 查询功能，同样也是递归查询每一层的节点，退出规则是，匹配到了*、匹配失败、或者匹配到了第len(parts)层节点。
 */

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 如果已经匹配完了，那么将pattern赋值给该node，表示它是一个完整的url
		// 这是递归的终止条件
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 没有匹配上，那就初始化一个，放到n节点的子列表中
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 接着插入下一个part节点
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 递归终止条件，找到末尾了或者通配符
		if n.pattern == "" {
			// pattern为空字符串表示它不是一个完整的url，匹配失败
			return nil
		}
		return n
	}

	part := parts[height]
	// 获取所有可能的子路径
	children := n.matchChildren(part)

	for _, child := range children {
		// 对于每条路径接着用下一part去查找
		result := child.search(parts, height+1)
		if result != nil {
			// 找到了即返回
			return result
		}
	}

	return nil
}

// 查找匹配的子节点，场景是用在插入时使用，找到1个匹配的就立即返回
func (n *node) matchChild(part string) *node {
	// 遍历n节点的所有子节点
	for _, child := range n.children {
		//if child.isWild {
		//	panic(part + "路由冲突，同级已经有" + child.part)
		//}
		//if child.part == part {
		//	return child
		//}

		// 动态匹配做强校验,防止路由注册覆盖
		if child.part == part || ((part[0] == ':' || part[0] == '*') && child.isWild) {
			return child
		}

	}
	return nil
}

// 查找所有匹配成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	wildNodes := make([]*node, 0)
	for _, child := range n.children {
		// 静态路由节点优先,动态路由节点延后
		if child.part == part {
			nodes = append(nodes, child)
		} else if child.isWild {
			wildNodes = append(wildNodes, child)
		}
	}
	nodes = append(nodes, wildNodes...)
	return nodes
}

package see

import (
	"strings"
)

// 前缀树实现
type node struct {
	pattern  string           // 是否一个完整的url，不是则为空字符串
	part     string           // URL块值，用/分割的部分，比如/abc/123，abc和123就是2个part
	children map[string]*node // 节点下的子节点
	isWild   bool             // 是否模糊匹配，比如:id*id这样的node就为true
}

/*
 * 对于路由来说，最重要的当然是注册与匹配了。开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler。
 * 因此，前缀树需要支持节点的插入与查询。插入功能很简单，递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个；
 * 有一点需要注意，/p/:lang/doc只有在第三层节点，即doc节点，pattern才会设置为/p/:lang/doc。
 * p和:lang节点的pattern属性皆为空。因此，当匹配结束时，我们可以使用n.pattern == ""来判断路由规则是否匹配成功。
 * 例如，/p/python虽能成功匹配到:lang，但:lang的pattern值为空，因此匹配失败。
 * 查询功能，同样也是递归查询每一层的节点，退出规则是，匹配到了*、匹配失败、或者匹配到了第len(parts)层节点。
 */

func (n *node) insert(pattern string, parts []string) {
	this := n
	for _, part := range parts {
		// 没有匹配上，那就初始化一个，放到n节点的子列表中
		if this.children[part] == nil {
			this.children[part] = &node{
				part:     part,
				children: make(map[string]*node),
				isWild:   part[0] == ':' || part[0] == '*',
			}
		}
		this = this.children[part]
	}
	this.pattern = pattern
}

func (n *node) search(parts []string) (*node, []Param) {
	var params []Param
	this := n
	for i, part := range parts {
		var temp string
		// child是否等于part
		for _, child := range this.children {
			// 处理当前节点参数
			if child.part == part || child.isWild {
				if child.part[0] == ':' {
					params = append(params, Param{child.part[1:], part})
				}
				if child.part[0] == '*' {
					params = append(params, Param{child.part[1:], strings.Join(parts[i:], "/")})
				}
				temp = child.part
			}
		}
		// 遇到通配符*，直接返回
		//if temp[0] == '*' {
		//    return this.children[temp], params
		//}
		this = this.children[temp]
	}
	return this, params
}

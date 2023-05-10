package web

import (
	"fmt"
	"strings"
)

// node 路由树的节点.
type node struct {
	pattern  string        // 完整的匹配路由，例如 /p/:param
	part     string        // 路由中的一部分，以 / 分割，例如 :param
	children []*node       // 子节点，例如 [*filepath] or [a, b, c, :param]，*filepath 为叶节点且与其它兄弟节点冲突，:param 只能有一个且必须在最后
	handlers []HandlerFunc // 路由对应的处理函数
	isWild   bool          // 是否模糊匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 插入节点-递归
func (n *node) insert(pattern string, parts []string, height int, handlers ...HandlerFunc) {
	// pattern 已经匹配完毕，停止递归
	if len(parts) == height {
		if n.pattern != "" {
			panic("pattern conflict")
		}
		n.pattern = pattern
		n.handlers = handlers
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil || child.part != part {
		child = &node{
			part:   part,
			isWild: isWildcard(part),
		}
		n.addChild(child)
	}
	child.insert(pattern, parts, height+1, handlers...)
}

func (n *node) addChild(child *node) {
	if isCatchAll(child.part) && len(n.children) > 0 {
		panic("catch-all conflict with existing children")
	}

	if isParam(child.part) && len(n.children) > 0 && n.children[len(n.children)-1].isWild {
		panic("param conflict with existing children")
	}

	if !isWildcard(child.part) && len(n.children) > 0 && n.children[len(n.children)-1].isWild {
		n.children = append(n.children[:len(n.children)-1], child, n.children[len(n.children)-1])
		return
	}

	n.children = append(n.children, child)
}

func isWildcard(part string) bool {
	return strings.HasPrefix(part, ":") || strings.HasPrefix(part, "*")
}

func isParam(part string) bool {
	return strings.HasPrefix(part, ":")
}

func isCatchAll(part string) bool {
	return strings.HasPrefix(part, "*")
}

// // parsePattern 将 pattern 解析为字符串数组（Only one * is allowed）
// func parsePattern(pattern string) []string {
// 	vs := strings.Split(pattern, "/")

// 	parts := make([]string, 0)
// 	for _, item := range vs {
// 		if item != "" {
// 			parts = append(parts, item)
// 			if item[0] == '*' {
// 				break
// 			}
// 		}
// 	}
// 	return parts
// }

// 查找节点-递归
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}

	return nil
}

// getRoute 查找对应的节点和模糊匹配所对应的解析参数
func (n *node) getRoute(pattern string) (*node, map[string]string) {
	params := make(map[string]string)
	searchParts := parsePattern(pattern)

	res := n.search(searchParts, 0)
	if res != nil {
		parts := parsePattern(res.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
			}
		}
		return res, params
	}
	return nil, nil
}

// // 遍历所有节点
// func (n *node) travel(list *([]*node)) {
// 	if n.pattern != "" {
// 		*list = append(*list, n)
// 	}
// 	for _, child := range n.children {
// 		child.travel(list)
// 	}
// }

// 匹配第一个子节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 匹配所有子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

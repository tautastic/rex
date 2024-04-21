package syntax

import (
	"fmt"

	"github.com/tautastic/rex/utils"
)

type Node struct {
	Label string
	Sub   []*Node
}

func (node Node) nodeStr() (str string) {
	if utils.IsAnyOf(node.Label, []string{
		"Disjunction",
		"Term",
		"Factor",
		"Assertion",
		"Quantifier",
		"Atom",
		"Perl",
		"Control",
		"HexSeq",
		"UniSeq",
		"Class",
		"ClassRange",
		"Literal",
	}) {
		str = fmt.Sprintf("<%v>", node.Label)
	} else {
		str = fmt.Sprintf("\"%v\"", node.Label)
	}
	return
}

func treeStr(t *Node, prefix string, arrow int) (str string) {
	switch arrow {
	case 0:
		str += fmt.Sprintf("%v\n", t.nodeStr())
		break
	case 1:
		str += fmt.Sprintf("%v%v %v\n", prefix, "└──", t.nodeStr())
		prefix += "    "
		break
	case 2:
		str += fmt.Sprintf("%v%v %v\n", prefix, "├──", t.nodeStr())
		prefix += "│   "
		break
	}
	for i := 0; i < len(t.Sub); i++ {
		if i < len(t.Sub)-1 && len(t.Sub) > 1 {
			str += treeStr(t.Sub[i], prefix, 2)
		} else {
			str += treeStr(t.Sub[i], prefix, 1)
		}
	}
	return
}

func (node Node) String() string {
	return treeStr(&node, "", 0)
}

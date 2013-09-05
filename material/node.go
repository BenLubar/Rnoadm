package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
	"math/rand"
	"sync"
)

const totalNodeSize = 1000000

type NodeLike interface {
	world.Visible

	Gather(world.InventoryLike, *big.Int, func(uint64) world.Visible) bool
	Quality() *big.Int
}

type Node struct {
	world.VisibleObject

	gathered uint64 // volume of material gathered from this node.

	mtx sync.Mutex
}

func init() {
	world.Register("resnode", NodeLike((*Node)(nil)))
}

func (n *Node) Save() (uint, interface{}, []world.ObjectLike) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return 1, n.gathered, nil
}

func (n *Node) Load(version uint, data interface{}, attached []world.ObjectLike) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	switch version {
	case 0:
		// do nothing - reset node
	case 1:
		n.gathered = data.(uint64)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (n *Node) Examine() (string, [][][2]string) {
	_, info := n.VisibleObject.Examine()

	return "a resource node.", info
}

func (n *Node) Gather(i world.InventoryLike, toolStrength *big.Int, item func(uint64) world.Visible) bool {
	nodeStrength := n.Outer().(NodeLike).Quality()

	gathered := n.gather(toolStrength, nodeStrength)
	if gathered == 0 {
		return false
	}
	gathered *= 100
	gathered += uint64(rand.Intn(100))
	n.mtx.Lock()
	n.gathered += gathered + 20
	if n.gathered > totalNodeSize {
		gathered -= n.gathered % totalNodeSize
		n.mtx.Unlock()
		n.Position().Remove(n.Outer())
	} else {
		n.mtx.Unlock()
	}

	i.GiveItem(item(gathered))
	return true
}

var eight = big.NewInt(8)

func (n *Node) gather(toolStrength, nodeStrength *big.Int) uint64 {
	var tmp big.Int

	if tmp.Cmp(nodeStrength) == 0 {
		return uint64(rand.Int63n(int64(totalNodeSize-n.gathered))) / 10
	}

	if tmp.Cmp(toolStrength) == 0 {
		return 0
	}

	r := rand.New(rand.NewSource(rand.Int63()))

	var toolScore, nodeScore big.Int
	toolScore.Add(toolScore.Rand(r, toolScore.Sub(toolStrength, tmp.Div(toolStrength, eight))), &tmp)
	nodeScore.Add(nodeScore.Rand(r, nodeScore.Sub(nodeStrength, tmp.Div(nodeStrength, eight))), &tmp)

	tmp.Sub(&toolScore, &nodeScore)
	if tmp.Sign() < 0 {
		return 0
	}
	return uint64(tmp.BitLen())
}

func (n *Node) Quality() *big.Int {
	return big.NewInt(1000)
}

func (n *Node) Blocking() bool {
	return true
}

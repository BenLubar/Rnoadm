package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"sync"
)

type NodeLike interface {
	world.Visible

	Strength() uint64
	Item() world.Visible
}

type Node struct {
	world.VisibleObject

	gathered uint // number of times gathered. affects success chance.

	mtx sync.Mutex
}

func init() {
	world.Register("resnode", NodeLike((*Node)(nil)))
}

func (n *Node) Save() (uint, interface{}, []world.ObjectLike) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	return 0, n.gathered, nil
}

func (n *Node) Load(version uint, data interface{}, attached []world.ObjectLike) {
	n.mtx.Lock()
	defer n.mtx.Unlock()

	switch version {
	case 0:
		n.gathered = data.(uint)
	default:
		panic(fmt.Sprintf("version %d unknown", version))
	}
}

func (n *Node) Gather(i world.InventoryLike, toolStrength uint64) bool {
	nodeStrength := n.Outer().(NodeLike).Strength()

	n.mtx.Lock()

	if toolStrength < 1<<63 && nodeStrength != 0 {
		if nodeStrength >= 1<<63 || toolStrength == 0 {
			n.mtx.Unlock()
			return false
		}
		toolScore := uint64(rand.Int63n(int64(toolStrength-toolStrength/8))) + toolStrength/8
		nodeScore := uint64(rand.Int63n(int64(nodeStrength-nodeStrength/8))) + nodeStrength/8
		if toolScore < nodeScore {
			n.mtx.Unlock()
			return false
		}
	}
	n.gathered++
	gathered := n.gathered
	n.mtx.Unlock()

	if rand.Intn(int(gathered)) != 0 {
		n.Position().Remove(n.Outer())
	}
	i.GiveItem(n.Outer().(NodeLike).Item())
	return true
}

func (n *Node) Strength() uint64 {
	return 1000
}

func (n *Node) Item() world.Visible {
	return n.Outer().(world.Visible)
}

func (n *Node) Blocking() bool {
	return true
}

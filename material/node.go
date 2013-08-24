package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/rand"
	"strconv"
	"sync"
)

type NodeLike interface {
	world.Visible

	Gather(world.InventoryLike, uint64, world.Visible) bool
	Strength() uint64
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

func (n *Node) Examine() (string, [][][2]string) {
	var info [][][2]string

	info = append(info, [][2]string{
		{strconv.FormatUint(n.Outer().(NodeLike).Strength(), 10), "#4fc"},
		{" strength", "#ccc"},
	})

	n.mtx.Lock()
	defer n.mtx.Unlock()

	if n.gathered != 0 {
		info = append(info, [][2]string{
			{strconv.FormatUint(uint64(n.gathered), 10), "#fc4"},
			{" gathered", "#ccc"},
		})
	}

	return "a resource node.", info
}

func (n *Node) Gather(i world.InventoryLike, toolStrength uint64, item world.Visible) bool {
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
	i.GiveItem(item)
	return true
}

func (n *Node) Strength() uint64 {
	return 1000
}

func (n *Node) Blocking() bool {
	return true
}

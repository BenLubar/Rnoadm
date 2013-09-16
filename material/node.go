package material

import (
	"fmt"
	"github.com/BenLubar/Rnoadm/world"
	"math/big"
	"math/rand"
	"sync"
)

type NodeLike interface {
	world.Visible

	Gather(world.InventoryLike, *big.Int, func(uint64) world.Visible) (bool, bool)
	Quality() *big.Int
	Size() uint64
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

func (n *Node) Gather(i world.InventoryLike, maxToolScore *big.Int, item func(uint64) world.Visible) (bool, bool) {
	maxNodeScore := (&big.Int{}).Mul(n.Outer().(NodeLike).Quality(), world.TuningNodeScoreMultiplier)

	gathered := n.gather(maxToolScore, maxNodeScore)
	if gathered == 0 {
		return false, false
	}
	n.mtx.Lock()
	remaining := n.Size() - n.gathered
	n.gathered += gathered + uint64(rand.Intn(20*int(gathered)/100))
	if gathered >= remaining {
		gathered = remaining
		n.mtx.Unlock()
		n.Position().Remove(n.Outer())
	} else {
		n.mtx.Unlock()
	}

	return i.GiveItem(item(gathered)), true
}

func (n *Node) gather(maxToolScore, maxNodeScore *big.Int) uint64 {
	r := rand.New(rand.NewSource(rand.Int63()))
	toolScore := (&big.Int{}).Rand(r, maxToolScore)
	nodeScore := (&big.Int{}).Rand(r, maxNodeScore)

	if toolScore.Cmp(nodeScore) <= 0 {
		return 0
	}

	volume := (&big.Int{}).Sub(toolScore, nodeScore)
	volume.Mul(volume, world.TuningMaxGatherVolume)
	volume.Div(volume, maxNodeScore)

	if volume.Cmp(world.TuningMaxGatherVolume) > 0 {
		return world.TuningMaxGatherVolume.Uint64()
	}

	return volume.Uint64()
}

func (n *Node) Quality() *big.Int {
	return big.NewInt(1000)
}

func (n *Node) Size() uint64 {
	return 10000000
}

func (n *Node) Blocking() bool {
	return true
}

type GatherSchedule struct {
	world.Object

	Tool    world.StatLike
	Target_ NodeLike
	Item    func(uint64) world.Visible
	started bool
}

func (s *GatherSchedule) Act(obj world.Living) (uint, bool) {
	if !s.started {
		s.started = true
		if o, ok := obj.(world.SendMessageLike); ok {
			o.SendMessage("you attempt to collect from the " + s.Target_.Name() + "...")
		}
		return 10, true
	}

	gave, collected := s.Target_.Gather(obj.(world.InventoryLike), s.Tool.Stat(world.StatGathering), s.Item)
	if o, ok := obj.(world.SendMessageLike); ok {
		if gave {
			o.SendMessage("you collect some of the " + s.Target_.Name() + ".")
		} else if collected {
			o.SendMessage("you collect some of the " + s.Target_.Name() + ", but it crumbles as you try to fit it into your pack.")
		} else {
			o.SendMessage("you fail to collect any of the " + s.Target_.Name() + ".")
		}
	}
	return 0, false
}

func (s *GatherSchedule) ShouldSave() bool         { return false }
func (s *GatherSchedule) Target() world.ObjectLike { return s.Target_ }

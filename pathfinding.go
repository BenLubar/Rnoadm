package main

import (
	"github.com/sperre/astar"
)

func FindPath(z *Zone, startX, startY, endX, endY uint8, allowDiagonal bool) *MoveSchedule {
	spaceBlocked := false
	mapData := make(astar.MapData, 256)
	mapData_ := make([]int, 256*256)
	for i := range mapData {
		mapData[i] = mapData_[:256]
		mapData_ = mapData_[256:]
	}
	for x := range mapData {
		for y := range mapData[x] {
			if z.Blocked(uint8(x), uint8(y)) {
				if uint8(x) == endX && uint8(y) == endY {
					spaceBlocked = true
					mapData[x][y] = astar.LAND
				} else if uint8(x) == startX && uint8(y) == startY {
					mapData[x][y] = astar.LAND
				} else {
					mapData[x][y] = astar.WALL
				}
			} else {
				mapData[x][y] = astar.LAND
			}
		}
	}
	nodes := astar.Astar(mapData, int(startX), int(startY), int(endX), int(endY), allowDiagonal)
	path := make(MoveSchedule, len(nodes))
	for i, n := range nodes {
		path[i][0] = uint8(n.X)
		path[i][1] = uint8(n.Y)
	}
	if len(path) > 1 && path[0][0] == startX && path[0][1] == startY {
		path = path[1:]
	}
	if spaceBlocked && len(path) > 0 {
		path = path[:len(path)-1]
	}
	return &path
}

package main

import (
	"github.com/sperre/astar"
)

func FindPath(z *Zone, startX, startY, endX, endY uint8, allowDiagonal bool) [][2]uint8 {
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
				} else {
					mapData[x][y] = astar.WALL
				}
			} else {
				mapData[x][y] = astar.LAND
			}
		}
	}
	nodes := astar.Astar(mapData, int(startX), int(startY), int(endX), int(endY), allowDiagonal)
	path := make([][2]uint8, len(nodes))
	for i, n := range nodes {
		path[i][0] = uint8(n.X)
		path[i][1] = uint8(n.Y)
	}
	if spaceBlocked && len(path) > 0 {
		path = path[:len(path)-1]
	}
	return path
}

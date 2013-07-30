package main

import (
	"github.com/sperre/astar"
)

func FindPath(z *Zone, startX, startY, endX, endY uint8, allowDiagonal bool) [][2]uint8 {
	mapData := astar.NewMapData(256, 256)
	for x := range mapData {
		for y := range mapData[x] {
			if z.Blocked(uint8(x), uint8(y)) {
				mapData[x][y] = astar.WALL
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
	return path
}

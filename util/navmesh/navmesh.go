package navmesh

import (
	"fmt"
)

type NavMesh struct {
	Vertices  []*V3      `json:"vertices"`
	Triangles [][3]int32 `json:"triangles"`
	Dijkstra  Dijkstra   `json:"dijkstra"`

	pathVerticesCache  map[int32][]int32              `json:"-"`
	pathTrianglesCache map[int32]map[int32][][3]int32 `json:"-"`
}

func (nm *NavMesh) InitWithVertices(verticesSlice [][]*V3) {
	nm.Vertices = make([]*V3, 0)
	verticesMap := map[string]int32{}
	index := int32(0)
	for _, vertices := range verticesSlice {
		for _, v := range vertices {
			key := getNavMeshKey(v)
			if _, ok := verticesMap[key]; !ok {
				verticesMap[key] = index
				nm.Vertices = append(nm.Vertices, v)
				index++
			}
		}
	}
	for _, d := range verticesSlice {
		nm.Triangles = append(nm.Triangles, [3]int32{
			verticesMap[getNavMeshKey(d[0])],
			verticesMap[getNavMeshKey(d[1])],
			verticesMap[getNavMeshKey(d[2])],
		})
	}
	nm.Dijkstra.CreateMatrixFromMesh(nm.Vertices, nm.Triangles)
}

// 在 X,Y 平面上寻路
func (nm *NavMesh) FindingPath(src, dest *V3) (*Path, error) {
	if srcID, destID := nm.getTriangleId(src), nm.getTriangleId(dest); srcID < 0 || destID < 0 {
		return nil, fmt.Errorf("finding path error,srcID:%d,destID:%d", srcID, destID)
	} else if srcID == destID {
		return &Path{PathList: []*V3{src, dest}}, nil
	} else if pathTriangle := nm.getPathTriangle(srcID, destID); pathTriangle == nil {
		return nil, fmt.Errorf("no path")
	} else {
		return nm.route(pathTriangle, src, dest)
	}
}

// v 点是否是可行走区域
func (nm *NavMesh) IsWalkable(v *V3) bool {
	return nm.getTriangleId(v) > -1
}

func (nm *NavMesh) route(pathTriangle [][3]int32, src, dest *V3) (*Path, error) {
	path := &Path{PathList: []*V3{src}}
	// 计算临边
	border := nm.createBorder(pathTriangle)
	// 目标点
	vertices := append(nm.Vertices, dest)
	border = append(border, int32(len(vertices))-1, int32(len(vertices))-1)

	// 第一个可视区域
	lineStart := src
	lastVisLeft, lastVisRight, lastPLeft, lastPRight := nm.updateVis(src, vertices, border, 0, 1)
	var res *V3
	for k := 2; k <= len(border)-2; k += 2 {
		curVisLeft, curVisRight, pLeft, pRight := nm.updateVis(lineStart, vertices, border, k, k+1)
		res = lastVisLeft.Cross(curVisRight)
		if res.Z > 0 { // 左拐点
			lineStart = vertices[border[lastPLeft]]
			path.PathList = append(path.PathList, lineStart)
			// 找到一条不共点的边作为可视区域
			i := 2 * (lastPLeft/2 + 1)
			for ; i <= len(border)-2; i += 2 {
				if border[lastPLeft] != border[i] && border[lastPLeft] != border[i+1] {
					lastVisLeft, lastVisRight, lastPLeft, lastPRight = nm.updateVis(lineStart, vertices, border, i, i+1)
					break
				}
			}

			k = i
			continue
		}

		res = lastVisRight.Cross(curVisLeft)
		if res.Z < 0 { // 右拐点
			lineStart = vertices[border[lastPRight]]
			path.PathList = append(path.PathList, lineStart)
			// 找到一条不共点的边
			i := 2 * (lastPRight/2 + 1)
			for ; i <= len(border)-2; i += 2 {
				if border[lastPRight] != border[i] && border[lastPRight] != border[i+1] {
					lastVisLeft, lastVisRight, lastPLeft, lastPRight = nm.updateVis(lineStart, vertices, border, i, i+1)
					break
				}
			}

			k = i
			continue
		}

		res = lastVisLeft.Cross(curVisLeft)
		if res.Z < 0 {
			lastVisLeft = curVisLeft
			lastPLeft = pLeft
		}

		res = lastVisRight.Cross(curVisRight)
		if res.Z > 0 {
			lastVisRight = curVisRight
			lastPRight = pRight
		}
	}
	path.PathList = append(path.PathList, dest)
	return path, nil
}

func (nm *NavMesh) getPathTriangle(srcID, destID int32) [][3]int32 {
	if nm.pathTrianglesCache == nil {
		nm.pathTrianglesCache = map[int32]map[int32][][3]int32{}
	}
	if pathVerticesMap, ok := nm.pathTrianglesCache[srcID]; ok {
		if pathTriangle, ok := pathVerticesMap[destID]; ok {
			return pathTriangle
		}
	}

	// Phase 1. Use Dijkstra to find shortest path on Triangles
	if nm.pathVerticesCache == nil {
		nm.pathVerticesCache = map[int32][]int32{}
	}
	pathVertices := nm.pathVerticesCache[srcID]
	if pathVertices == nil {
		pathVertices = nm.Dijkstra.Run(srcID)
		nm.pathVerticesCache[srcID] = pathVertices
	}

	// Phase 2.  construct path indices
	// Check if this path include src & dest
	if nm.pathTrianglesCache[srcID] == nil {
		nm.pathTrianglesCache[srcID] = map[int32][][3]int32{}
	}
	pathTriangle := [][3]int32{nm.Triangles[destID]}
	prevID := destID
	for {
		curID := pathVertices[prevID]
		if curID == -1 {
			nm.pathTrianglesCache[srcID][destID] = nil
			return nil
		}
		pathTriangle = append([][3]int32{nm.Triangles[curID]}, pathTriangle...)
		if curID == srcID {
			break
		}
		prevID = curID
	}
	nm.pathTrianglesCache[srcID][destID] = pathTriangle
	return pathTriangle
}

func (nm *NavMesh) createBorder(list [][3]int32) []int32 {
	var border []int32
	for k := 0; k < len(list)-1; k++ {
		for _, i := range list[k] {
			for _, j := range list[k+1] {
				if i == j {
					border = append(border, i)
				}
			}
		}
	}
	return border
}

func (nm *NavMesh) updateVis(v0 *V3, vertices []*V3, indices []int32, i1, i2 int) (l, r *V3, left, right int) {
	var leftVec, rightVec, res *V3
	leftVec = vertices[indices[i1]].Sub(v0)
	rightVec = vertices[indices[i2]].Sub(v0)
	res = leftVec.Cross(rightVec)
	if res.Z > 0 {
		return rightVec, leftVec, i2, i1
	} else {
		return leftVec, rightVec, i1, i2
	}
}

func (nm *NavMesh) getTriangleId(v *V3) int32 {
	for k := 0; k < len(nm.Triangles); k++ {
		if inside(v,
			&V3{X: nm.Vertices[nm.Triangles[k][0]].X, Y: nm.Vertices[nm.Triangles[k][0]].Y},
			&V3{X: nm.Vertices[nm.Triangles[k][1]].X, Y: nm.Vertices[nm.Triangles[k][1]].Y},
			&V3{X: nm.Vertices[nm.Triangles[k][2]].X, Y: nm.Vertices[nm.Triangles[k][2]].Y}) {
			return int32(k)
		}
	}
	return -1
}

func inside(pt, v1, v2, v3 *V3) bool {
	b1 := sign(pt, v1, v2) <= 0
	b2 := sign(pt, v2, v3) <= 0
	b3 := sign(pt, v3, v1) <= 0
	return (b1 == b2) && (b2 == b3)
}

func sign(p1, p2, p3 *V3) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

func getNavMeshKey(v *V3) string {
	return fmt.Sprintf("%.0f/%.0f/%.0f", v.X, v.Y, v.Z)
}

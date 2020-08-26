package navmesh

import (
	"container/heap"
	"math"
)

const largeNumber = math.MaxInt32

// Triangle Heap
type WeightedTriangle struct {
	ID     int32  `json:"id"` // triangle ID
	Weight uint32 `json:"w"`
}

type TriangleHeap struct {
	triangles []WeightedTriangle
	indices   map[int32]int
}

func NewTriangleHeap() *TriangleHeap {
	h := new(TriangleHeap)
	h.indices = make(map[int32]int)
	return h
}

func (th *TriangleHeap) Len() int {
	return len(th.triangles)
}

func (th *TriangleHeap) Less(i, j int) bool {
	return th.triangles[i].Weight < th.triangles[j].Weight
}

func (th *TriangleHeap) Swap(i, j int) {
	th.triangles[i], th.triangles[j] = th.triangles[j], th.triangles[i]
	th.indices[th.triangles[i].ID] = i
	th.indices[th.triangles[j].ID] = j
}

func (th *TriangleHeap) Push(x interface{}) {
	th.triangles = append(th.triangles, x.(WeightedTriangle))
	n := len(th.triangles)
	th.indices[th.triangles[n-1].ID] = n - 1
}

func (th *TriangleHeap) Pop() interface{} {
	n := len(th.triangles)
	x := th.triangles[n-1]
	th.triangles = th.triangles[:n-1]
	return x
}

func (th *TriangleHeap) DecreaseKey(id int32, weight uint32) {
	if index, ok := th.indices[id]; ok {
		th.triangles[index].Weight = weight
		heap.Fix(th, index)
		return
	} else {
		heap.Push(th, WeightedTriangle{id, weight})
	}
}

type Mesh struct {
	Vertices  []V3       // vertices
	Triangles [][3]int32 // triangles
}

// Dijkstra
type Dijkstra struct {
	Matrix map[int32][]WeightedTriangle // all edge for nodes
}

// create neighbour matrix
func (d *Dijkstra) CreateMatrixFromMesh(vertices []*V3, triangles [][3]int32) {
	d.Matrix = make(map[int32][]WeightedTriangle)
	for i := 0; i < len(triangles); i++ {
		for j := 0; j < len(triangles); j++ {
			if i == j {
				continue
			}

			if len(intersect(triangles[i], triangles[j])) == 2 {
				x1 := (vertices[triangles[i][0]].X + vertices[triangles[i][1]].X + vertices[triangles[i][2]].X) / 3.0
				y1 := (vertices[triangles[i][0]].Y + vertices[triangles[i][1]].Y + vertices[triangles[i][2]].Y) / 3.0
				x2 := (vertices[triangles[j][0]].X + vertices[triangles[j][1]].X + vertices[triangles[j][2]].X) / 3.0
				y2 := (vertices[triangles[j][0]].Y + vertices[triangles[j][1]].Y + vertices[triangles[j][2]].Y) / 3.0
				weight := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
				d.Matrix[int32(i)] = append(d.Matrix[int32(i)], WeightedTriangle{int32(j), uint32(weight)})
			}
		}
	}
}

func intersect(a [3]int32, b [3]int32) []int32 {
	var inter []int32
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				inter = append(inter, a[i])
			}
		}
	}
	return inter
}

func (d *Dijkstra) Run(srcID int32) []int32 {
	// triangle heap
	h := NewTriangleHeap()
	// min distance records
	dist := make([]uint32, len(d.Matrix))
	for i := 0; i < len(dist); i++ {
		dist[i] = largeNumber
	}
	distLen := int32(len(dist))
	// previous
	prev := make([]int32, len(d.Matrix))
	for i := 0; i < len(prev); i++ {
		prev[i] = -1
	}
	// visit map
	visited := make([]bool, len(d.Matrix))

	// source vertex, the first vertex in Heap
	dist[srcID] = 0
	heap.Push(h, WeightedTriangle{srcID, 0})

	for h.Len() > 0 { // for every un-visited vertex, try relaxing the path
		// pop the min element
		u := heap.Pop(h).(WeightedTriangle)
		if visited[u.ID] {
			continue
		}
		// current known shortest distance to u
		distU := dist[u.ID]
		// mark the vertex as visited.
		visited[u.ID] = true

		// for each neighbor v of u:
		for _, v := range d.Matrix[u.ID] {
			alt := distU + v.Weight // from src->u->v
			if v.ID < distLen &&    // 越界判断
				alt < dist[v.ID] {
				dist[v.ID] = alt
				prev[v.ID] = u.ID
				if !visited[v.ID] {
					h.DecreaseKey(v.ID, alt)
				}
			}
		}
	}
	return prev
}

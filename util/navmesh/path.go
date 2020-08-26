package navmesh

type Path struct {
	PathList     []*V3
	nextIndex    int
	currPosition *V3
}

// 移动一段距离 返回 当前位置,全部寻路是否完成
func (p *Path) Move(moveDistance float64) (*V3, bool) {
	if len(p.PathList) == 0 {
		return nil, true
	} else if p.currPosition == nil {
		p.currPosition = p.PathList[0]
		p.nextIndex = 1
	}

	var distance float64
	for p.nextIndex < len(p.PathList) {
		if p.currPosition, distance = p.moveTo(moveDistance, p.currPosition, p.PathList[p.nextIndex]); distance > 0 {
			p.nextIndex++
		} else {
			return p.currPosition, false
		}
	}
	return p.PathList[len(p.PathList)-1], true
}

// 只算 X,Y 平面
func (p *Path) moveTo(moveDistance float64, fromPosition, toPosition *V3) (*V3, float64) {
	toNextPositionDistance := fromPosition.Distance(toPosition)
	if toNextPositionDistance < moveDistance {
		return toPosition, moveDistance - toNextPositionDistance
	} else {
		rate := moveDistance / toNextPositionDistance
		return NewV3(
			fromPosition.X+(toPosition.X-fromPosition.X)*rate,
			fromPosition.Y+(toPosition.Y-fromPosition.Y)*rate,
			0,
		), 0
	}
}

package snowflake

type overCostActionArg struct {
	ActionType             int32
	TimeTick               int64
	WorkerId               uint16
	OverCostCountInOneTerm int32
	GenCountInOneTerm      int32
	TermIndex              int32
}

func (oa overCostActionArg) OverCostActionArg(workerId uint16, timeTick int64, actionType int32, overCostCountInOneTerm int32, genCountWhenOverCost int32, index int32) {
	oa.ActionType = actionType
	oa.TimeTick = timeTick
	oa.WorkerId = workerId
	oa.OverCostCountInOneTerm = overCostCountInOneTerm
	oa.GenCountInOneTerm = genCountWhenOverCost
	oa.TermIndex = index
}

package types

type SelectColumn interface {
	Column() string
}

type LolCondition interface {
	Render() string
	And(LolCondition) LolCondition
	Or(LolCondition) LolCondition
	Next() LolCondition
	SetNext(LolCondition)
}

type HasNext struct {
	next LolCondition
}

func (this *HasNext) Next() LolCondition {
	return this.next
}

func (this *HasNext) SetNext(cond LolCondition) {
	prev := this.next
	next := this.next
	for next != nil {
		prev = next
		next = next.Next()
	}
	if prev != nil {
		prev.SetNext(cond)
	} else {
		this.next = cond
	}
}
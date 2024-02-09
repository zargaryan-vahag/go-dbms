package aggregator

import (
	"errors"

	"go-dbms/pkg/types"

	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Float | constraints.Integer
}

type comparable interface {
	constraints.Ordered
}

type AggregatorType string

const (
	SUM AggregatorType = "SUM"
	MAX AggregatorType = "MAX"
	MIN AggregatorType = "MIN"
	AVG AggregatorType = "AVG"
	COUNT AggregatorType = "COUNT"
)

var aggregators = map[AggregatorType]struct{}{
	SUM:   {},
	MAX:   {},
	MIN:   {},
	AVG:   {},
	COUNT: {},
}

type AggregatorBase struct {
	Arguments []string
}

func (ab *AggregatorBase) Value() types.DataType {
	panic(errors.New("unimplemented"))
}

func (ab *AggregatorBase) Apply(row map[string]types.DataType) {
	panic(errors.New("unimplemented"))
}

type Aggregator interface {
	Value() types.DataType
	Apply(row map[string]types.DataType)
}

func IsAggregator(fn string) bool {
	_, ok := aggregators[AggregatorType(fn)]
	return ok
}

func New(name AggregatorType, args []string) Aggregator {
	ab := &AggregatorBase{args}
	switch name {
		case AVG: {
			return &AggregationAVG{AggregatorBase: ab}
		}
		case COUNT: {
			return &AggregationCOUNT{AggregatorBase: ab}
		}
		case MAX: {
			return &AggregationMAX{AggregatorBase: ab}
		}
		case MIN: {
			return &AggregationMIN{AggregatorBase: ab}
		}
		case SUM: {
			return &AggregationSUM[int64]{
				AggregatorBase: ab,
				Meta: &types.DataTypeINTEGERMeta{Signed: true, ByteSize: 8},
			}
		}
		default: {
			panic(errors.New("unknown aggregate function"))
		}
	}
}

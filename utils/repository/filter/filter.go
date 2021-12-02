package filter

type Operator string
const (
	OP_EQ      Operator = "$EQ"
	OP_GT      Operator = "$GT"
	OP_GTE     Operator = "$GTE"
	OP_LT      Operator = "$LT"
	OP_LTE     Operator = "$LTE"
	OP_IN      Operator = "$IN"
	OP_LIKE    Operator = "$LIKE"
	OP_INVALID Operator = ""
)

type Filterer interface {
	GenerateStatement() string
}

type Filter []Criteria

type FilterMap map[string]Filterer

type Criteria struct {
	Key string
	Value interface{}
}

type BindVar struct {
	Value interface{}
}

//func (f *Filter) Generate() (string, []interface{}) {
//	where := ""
//	for _, cr := range f {
//		switch cr.Key {
//		case "$or":
//			break
//		case "$and":
//			break
//		}
//	}
//}
//
//func ParseMap(filterMap map[string]interface{}) Filter {
//	for k, v := range filterMap {
//
//	}
//}

//func (f FilterMap) GenerateWhereCriteriaString() string {
//	where := ""
//	for k, v := range f {
//		if len(where) > 0 {
//			where += " AND "
//		}
//	}
//}
//
//func (f FilterMap) GenerateArgsMap() map[string]interface{} {
//
//}
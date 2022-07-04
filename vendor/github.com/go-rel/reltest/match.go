package reltest

import (
	"reflect"

	"github.com/go-rel/rel"
)

type any struct{}

func (any) String() string {
	return "<Any>"
}

var Any interface{} = any{}

func matchQuery(mock rel.Query, input rel.Query) bool {
	return matchTable(mock.Table, input.Table) &&
		matchSelectQuery(mock.SelectQuery, input.SelectQuery) &&
		matchJoinQuery(mock.JoinQuery, input.JoinQuery) &&
		matchFilterQuery(mock.WhereQuery, input.WhereQuery) &&
		matchGroupQuery(mock.GroupQuery, input.GroupQuery) &&
		matchSortQuery(mock.SortQuery, input.SortQuery) &&
		mock.OffsetQuery == input.OffsetQuery &&
		mock.LimitQuery == input.LimitQuery &&
		mock.LockQuery == input.LockQuery &&
		matchSQLQuery(mock.SQLQuery, input.SQLQuery) &&
		mock.UnscopedQuery == input.UnscopedQuery &&
		mock.ReloadQuery == input.ReloadQuery &&
		mock.CascadeQuery == input.CascadeQuery &&
		reflect.DeepEqual(mock.PreloadQuery, input.PreloadQuery)

}

func matchTable(mock string, input string) bool {
	return mock == "" || input == "" || mock == input
}

func matchSelectQuery(mock rel.SelectQuery, input rel.SelectQuery) bool {
	return mock.OnlyDistinct == input.OnlyDistinct && reflect.DeepEqual(mock.Fields, input.Fields)
}

func matchJoinQuery(mocks []rel.JoinQuery, inputs []rel.JoinQuery) bool {
	if len(mocks) != len(inputs) {
		return false
	}

	for i := range mocks {
		if mocks[i].Assoc != "" && mocks[i].Assoc == inputs[i].Assoc {
			continue
		}

		// TODO: argument support any
		if mocks[i].Mode != inputs[i].Mode ||
			mocks[i].Table != inputs[i].Table ||
			mocks[i].From != inputs[i].From ||
			mocks[i].To != inputs[i].To ||
			!reflect.DeepEqual(mocks[i].Arguments, inputs[i].Arguments) {
			return false
		}
	}

	return true
}

func matchFilterQuery(mock rel.FilterQuery, input rel.FilterQuery) bool {
	switch v := mock.Value.(type) {
	case rel.SubQuery:
		if inputSubQuery, _ := input.Value.(rel.SubQuery); v.Prefix != inputSubQuery.Prefix || !matchQuery(v.Query, inputSubQuery.Query) {
			return false
		}
	case rel.Query:
		if bQuery, ok := input.Value.(rel.Query); !ok || !matchQuery(v, bQuery) {
			return false
		}
	default:
		if mock.Type != input.Type ||
			mock.Field != input.Field ||
			(!reflect.DeepEqual(mock.Value, input.Value) && mock.Value != Any) ||
			len(mock.Inner) != len(input.Inner) {
			return false
		}
	}

	for i := range mock.Inner {
		if !matchFilterQuery(mock.Inner[i], input.Inner[i]) {
			return false
		}
	}

	return true
}

func matchGroupQuery(mock rel.GroupQuery, input rel.GroupQuery) bool {
	return reflect.DeepEqual(mock.Fields, input.Fields) && matchFilterQuery(mock.Filter, input.Filter)
}

func matchSortQuery(mocks []rel.SortQuery, inputs []rel.SortQuery) bool {
	return reflect.DeepEqual(mocks, inputs)
}

func matchSQLQuery(mock rel.SQLQuery, input rel.SQLQuery) bool {
	if mock.Statement != input.Statement || len(mock.Values) != len(input.Values) {
		return false
	}

	for i := range mock.Values {
		if mock.Values[i] != input.Values[i] && mock.Values[i] != Any {
			return false
		}
	}

	return true
}

func matchMutators(mocks []rel.Mutator, inputs []rel.Mutator) bool {
	if len(mocks) != len(inputs) {
		return false
	}

	for i := range mocks {
		switch mock := mocks[i].(type) {
		case rel.Mutate:
			if input, ok := inputs[i].(rel.Mutate); !ok || !matchMutate(mock, input) {
				return false
			}
		case rel.Changeset:
			if input, ok := inputs[i].(rel.Changeset); !ok || !reflect.DeepEqual(mock.Changes(), input.Changes()) {
				return false
			}
		default:
			if !reflect.DeepEqual(mock, inputs[i]) {
				return false
			}
		}
	}

	return true
}

func matchMutate(mock rel.Mutate, input rel.Mutate) bool {
	if mock.Type != input.Type || mock.Field != input.Field {
		return false
	}

	if mock.Type == rel.ChangeFragmentOp {
		var (
			mockArgs, _  = mock.Value.([]interface{})
			inputArgs, _ = input.Value.([]interface{})
		)

		if len(mockArgs) != len(inputArgs) {
			return false
		}

		for i := range mockArgs {
			if mockArgs[i] != Any && !reflect.DeepEqual(mockArgs[i], inputArgs[i]) {
				return false
			}
		}

		return true
	}

	return mock.Value == Any || reflect.DeepEqual(mock.Value, input.Value)
}

func matchMutates(mocks []rel.Mutate, inputs []rel.Mutate) bool {
	if len(mocks) != len(inputs) {
		return false
	}

	for i := range mocks {
		if !matchMutate(mocks[i], inputs[i]) {
			return false
		}
	}

	return true
}

func matchContains(mock interface{}, input interface{}) bool {
	var (
		rva = reflect.Indirect(reflect.ValueOf(mock))
		rta = rva.Type()
		rvb = reflect.Indirect(reflect.ValueOf(input))
	)

	for i := 0; i < rva.NumField(); i++ {
		fva := rva.Field(i)
		if fva.IsZero() {
			continue
		}

		fvb := rvb.FieldByName(rta.Field(i).Name)
		if !fvb.IsValid() || !reflect.DeepEqual(fva.Interface(), fvb.Interface()) {
			return false
		}
	}

	return true
}

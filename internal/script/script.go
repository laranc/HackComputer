package script

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func Load(path string) (*lua.LTable, error) {
	l := lua.NewState()
	defer l.Close()
	if err := l.DoFile(path); err != nil {
		return nil, err
	}
	table := l.Get(-1)
	if tbl, ok := table.(*lua.LTable); ok {
		return tbl, nil
	}
	return nil, fmt.Errorf("Error parsing script")
}

func Get[T lua.LValue](table *lua.LTable, value string) T {
	return table.RawGetString(value).(T)
}

func SliceFromTable[T lua.LValue](table *lua.LTable) []T {
	output := make([]T, 0, table.Len())
	table.ForEach(func(key, value lua.LValue) {
		if value != lua.LNil {
			output = append(output, value.(T))
		}
	})
	return output
}

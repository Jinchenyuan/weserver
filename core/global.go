package core

import (
	"sync/atomic"
)

var globalMesa atomic.Value // stores *Mesa

func SetGlobalMesa(m *Mesa) { globalMesa.Store(m) }
func GetGlobalMesa() *Mesa {
	v := globalMesa.Load()
	if v == nil {
		return nil
	}
	return v.(*Mesa)
}

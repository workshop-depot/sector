package sector

import "sync"

//-----------------------------------------------------------------------------

// Factory /
type Factory interface {
	// Fill , must accept only pointer and it fills that pointer
	Fill(interface{}) bool
}

//-----------------------------------------------------------------------------

// FactoryFunc implements Factory for a normal func
type FactoryFunc func(interface{}) bool

// Fill implements Factory
func (v FactoryFunc) Fill(slot interface{}) bool {
	return v(slot)
}

//-----------------------------------------------------------------------------

// FactoryRepo for dynamically piling factories
type FactoryRepo struct {
	m         *sync.RWMutex
	factories []Factory
}

// NewFactoryRepo /
func NewFactoryRepo() *FactoryRepo {
	fp := FactoryRepo{}
	fp.m = new(sync.RWMutex)

	return &fp
}

// Register /
func (v *FactoryRepo) Register(factory Factory) {
	v.m.Lock()
	defer v.m.Unlock()

	v.factories = append(v.factories, factory)
}

// Fill fills targetSlot which is a pointer
func (v *FactoryRepo) Fill(targetSlot interface{}) bool {
	v.m.RLock()
	defer v.m.RUnlock()

	for _, fac := range v.factories {
		if fac.Fill(targetSlot) {
			return true
		}
	}

	return false
}

// FillAll fills all slots, which are pointers, a helper to reduce locking
func (v *FactoryRepo) FillAll(slots ...interface{}) {
	v.m.RLock()
	defer v.m.RUnlock()

	for _, targetSlot := range slots {
		targetSlot := targetSlot
		for _, fac := range v.factories {
			if fac.Fill(targetSlot) {
				break
			}
		}
	}
}

//-----------------------------------------------------------------------------

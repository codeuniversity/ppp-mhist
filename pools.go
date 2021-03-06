package mhist

import (
	"sync"
)

//Pools holds the pools for the different measurement types (only one for now)
type Pools struct {
	Store            *Store
	messagePool      *sync.Pool
	measurementPools map[MeasurementType]*sync.Pool
}

//MeasurementSlices for moving types of Measurements around
type MeasurementSlices map[MeasurementType][]Measurement

//NewPools returns the constructed pool handler
func NewPools(store *Store) *Pools {
	pools := &Pools{
		Store: store,
		messagePool: &sync.Pool{
			New: func() interface{} {
				return &Message{}
			},
		},
	}
	pools.measurementPools = map[MeasurementType]*sync.Pool{
		MeasurementNumerical: &sync.Pool{
			New: func() interface{} {
				slices, ok := grabSlicesFromStore(store)
				if ok {
					numericalSlice := slices[MeasurementNumerical]
					if len(numericalSlice) > 0 {
						measurement := numericalSlice[0]
						rest := numericalSlice[1:]
						slices[MeasurementNumerical] = rest
						pools.fill(slices)
						return measurement
					}
				}
				return &Numerical{}
			},
		},
		MeasurementCategorical: &sync.Pool{
			New: func() interface{} {
				slices, ok := grabSlicesFromStore(store)
				if ok {
					categoricalSlice := slices[MeasurementCategorical]
					if len(categoricalSlice) > 0 {
						measurement := categoricalSlice[0]
						rest := categoricalSlice[1:]
						slices[MeasurementCategorical] = rest
						pools.fill(slices)
						return measurement
					}
				}
				return &Categorical{}
			},
		},
	}
	return pools
}

//GetNumericalMeasurement out of the correct pool
func (pools *Pools) GetNumericalMeasurement() *Numerical {
	return pools.measurementPools[MeasurementNumerical].Get().(*Numerical)
}

//GetCategoricalMeasurement out of the correct pool
func (pools *Pools) GetCategoricalMeasurement() *Categorical {
	return pools.measurementPools[MeasurementCategorical].Get().(*Categorical)
}

//PutMeasurement out of the correct pool
func (pools *Pools) PutMeasurement(m Measurement) {
	pools.measurementPools[m.Type()].Put(m)
}

//GetMessage from MessagePool
func (pools *Pools) GetMessage() *Message {
	return pools.messagePool.Get().(*Message)
}

//PutMessage into MessagePool
func (pools *Pools) PutMessage(m *Message) {
	pools.messagePool.Put(m)
}

func grabSlicesFromStore(store *Store) (slices MeasurementSlices, ok bool) {
	if store.IsOverSoftLimit() {
		slices := store.ShrinkStore()
		if !store.IsOverMaxSize() {
			return slices, true
		}
	}
	return nil, false
}

func (pools *Pools) fill(slices MeasurementSlices) {
	for key, slice := range slices {
		for _, measurement := range slice {
			pools.measurementPools[key].Put(measurement)
		}
	}
}

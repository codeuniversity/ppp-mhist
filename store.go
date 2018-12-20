package mhist

import (
	"fmt"
)

//Store is responsible for handling Storage of different kinds of measurements
type Store struct {
	subscribers  SubscriberSlice
	replications SubscriberSlice
	diskStore    *DiskStore
}

//NewStore from diskstore, that handles subscribers and replications
func NewStore(diskStore *DiskStore) *Store {
	store := &Store{
		diskStore: diskStore,
	}
	store.AddSubscriber(diskStore)
	return store
}

//AddSubscriber to Store
func (s *Store) AddSubscriber(sub Subscriber) {
	s.subscribers = append(s.subscribers, sub)
}

//AddReplication to Store
func (s *Store) AddReplication(rep Subscriber) {
	s.replications = append(s.replications, rep)
}

//Add named measurement to correct Series
func (s *Store) Add(name string, m Measurement, isReplication bool) {
	if !isReplication {
		s.replications.NotifyAll(name, m)
	}
	s.subscribers.NotifyAll(name, m)
}

//GetMeasurementsInTimeRange from disk store
func (s *Store) GetMeasurementsInTimeRange(start, end int64, filterDefinition FilterDefinition) map[string][]Measurement {
	if s.diskStore != nil {
		return s.diskStore.GetMeasurementsInTimeRange(start, end, filterDefinition)
	}
	return map[string][]Measurement{}
}

//GetStoredMetaInfo from Diskstore
func (s *Store) GetStoredMetaInfo() []MeasurementTypeInfo {
	if s.diskStore == nil {
		fmt.Println("no diskstore added to store, can't access metadata")
		return []MeasurementTypeInfo{}
	}

	return s.diskStore.GetAllStoredInfos()
}

//Shutdown diskStore
func (s *Store) Shutdown() {
	s.diskStore.Shutdown()
}

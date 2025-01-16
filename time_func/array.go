package time_func

import (
	"sort"
	"sync"
)

type log struct {
	startTime int
	count     int
	msg       string
}

type safeArray struct {
	array []log
	sync.Mutex
}

func (sa *safeArray) Add(time int, count int, msg string) {
	sa.Lock()
	defer sa.Unlock()
	sa.array = append(sa.array, log{time, count, msg})
}

func (sa *safeArray) ReserveAndClear() []string {
	sa.Lock()
	defer sa.Unlock()
	newArray := make([]log, len(sa.array))
	copy(newArray, sa.array)

	sort.Slice(newArray, func(i, j int) bool {
		return newArray[i].startTime+newArray[i].count < newArray[j].startTime+newArray[j].count
	})
	msgs := make([]string, 0, len(sa.array))
	for _, l := range newArray {
		msgs = append(msgs, l.msg)
	}
	sa.array = make([]log, 0)
	return msgs
}

var defaultSafeArray = &safeArray{}

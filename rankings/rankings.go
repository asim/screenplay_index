package rankings

import (
	//"fmt"
	//"time"
	"container/heap"
)

type rollingCounter struct {
	w        int // sliding window length
	c        int // current window
	counters []map[string]int
}

type object interface {
	ID() string
}

type Item struct {
	value    object
	priority int
	index    int
}

type rankings []*Item

type RankingsManager struct {
	rc *rollingCounter // rolling counter
	n  int             // number of elements to rank
	r  *rankings
}

func newRollingCounter(w int) *rollingCounter {
	return &rollingCounter{
		w: w,
		c: 0,
		counters: []map[string]int{
			make(map[string]int),
		},
	}
}

func newRankings() *rankings {
	r := &rankings{}
	heap.Init(r)
	return r
}

func newRankingsManager(n int, w int) *RankingsManager {
	return &RankingsManager{
		rc: newRollingCounter(w),
		n:  n,
		r:  newRankings(),
	}
}

func (i *Item) Value() object {
	return i.value
}

func (r *rollingCounter) incr(counter string) {
	r.counters[r.c][counter]++
}

func (r *rollingCounter) get(counter string) int {
	return r.counters[r.c][counter]
}

func (r *rollingCounter) getTotal(counter string) int {
	count := 0
	for _, c := range r.counters {
		count += c[counter]
	}
	return count
}

func (r *rollingCounter) slide() {
	r.counters = append(r.counters, make(map[string]int))
	if len(r.counters) > r.w {
		r.counters = r.counters[1:]
	}
	r.c = len(r.counters) - 1
}

func (r rankings) Len() int {
	return len(r)
}

func (r rankings) Less(i, j int) bool {
	return r[i].priority < r[j].priority
}

func (r rankings) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
	r[i].index = i
	r[j].index = j
}

func (r *rankings) Push(x interface{}) {
	n := len(*r)
	item := x.(*Item)
	item.index = n
	*r = append(*r, item)
}

func (r *rankings) Pop() interface{} {
	old := *r
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*r = old[0 : n-1]
	return item
}

func (r *rankings) update(item *Item, value object, priority int) {
	heap.Remove(r, item.index)
	item.value = value
	item.priority = priority
	heap.Push(r, item)
}

func (r *RankingsManager) update(item object) {
	r.rc.incr(item.ID())
	i := r.rc.getTotal(item.ID())

	var meti *Item
	for _, m := range *r.r {
		if m.value.ID() == item.ID() {
			meti = m
			break
		}
	}

	if meti == nil {
		meti = &Item{
			value:    item,
			priority: i,
		}
		heap.Push(r.r, meti)
	} else {
		r.r.update(meti, meti.value, i)
	}

	if r.r.Len() > r.n {
		heap.Pop(r.r)
	}
}

func (r *RankingsManager) GetRankings() *rankings {
	return r.r
}

func (r *RankingsManager) Update(item object) {
	r.update(item)
}

func (r *RankingsManager) Slide() {
	r.rc.slide()
}

func New(n, w int) *RankingsManager {
	return newRankingsManager(n, w)
}

/*

type script struct{
	name string
}

func (s *script) Name() string {
	return s.name
}

func main() {
	r := newRankingsManager(6, 3)

	words := []*script{}

	for i := 0; i<10; i++ {
		words = append(words, &script{fmt.Sprintf("foo%d", i)})
	}

	tick := time.NewTicker(time.Millisecond)
	for {
		select {
		case <-tick.C:
			fmt.Println("slide")
			r.rc.slide()
		default:
			for i := 0; i<1000; i++ {
				j := int(time.Now().UnixNano()) % len(words)
				word := words[j]
				r.update(word)
			}

		}

		for r.r.Len() > 0 {
			i := heap.Pop(r.r).(*Item)
			fmt.Printf("item %#v %s %d %d\n", i, i.value.Name(), r.rc.get(i.value.Name()), r.rc.getTotal(i.value.Name()))
		}
		fmt.Println("")
	}
}

*/

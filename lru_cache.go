package main

import "sync"

import (
	"container/list"
	"errors"
	"fmt"
	"time"
)

type cacheNode struct {
	key interface{}
	value interface{}
	count int
}

type LRUCache struct {
	lock    sync.Mutex
	pList   *list.List
	pList_candidate   *list.List
	cache   map[interface{}]*list.Element
	capacity int
	threshold int
}

func newCacheNode(k,v interface{}) *cacheNode {
	return &cacheNode{key: k, value: v, count: 1}
}

func NewLRUCache(cap int, level int) (*LRUCache){
	return &LRUCache {
		capacity: cap,
		pList: list.New(),
		pList_candidate: list.New(),
		//count: make(map[interface{}]int)
		cache: make(map[interface{}]*list.Element),
		threshold: level,
		}
}

func (lc *LRUCache) Get(key interface{}) (v interface{}, ret bool, e error) {
	if lc == nil || lc.cache == nil {
		e = errors.New("LRU Cache not initialized")
		return v,false, e
	}

	lc.lock.Lock()
	defer lc.lock.Unlock()

	if pElm, ok := lc.cache[key]; ok {
		//exist
		//move to head, higher priority
		lc.pList.MoveToFront(pElm)
		return pElm.Value.(*cacheNode).value,true,nil
	}
	return v,false,nil
}

func (lc *LRUCache) setCache(key, value interface{}) (error) {
	if lc == nil || lc.cache == nil {
		e := errors.New("LRU Cache not initialized")
		return e
	}

	if pElm, ok := lc.cache[key]; ok {
		//exist
		pElm.Value.(*cacheNode).value = value
		pElm.Value.(*cacheNode).count++
		//move to head, higher priority
		lc.pList.MoveToFront(pElm)
		return nil
	}

	//not exist, then create cache node
	item := newCacheNode(key, value)
	newElement := lc.pList.PushFront(item)
	lc.cache[key] = newElement

	if lc.pList.Len() > lc.capacity {
		//delete last one
		lastElement := lc.pList.Back()
		if lastElement == nil {
			return nil
		}
		//delete cache map
		delete(lc.cache, lastElement.Value.(*cacheNode).key)
		//delete cache list
		lc.pList.Remove(lastElement)
	}
	return nil
}

func (lc *LRUCache) SetDirect(key, value interface{}) (error) {
	lc.lock.Lock()
	defer lc.lock.Unlock()
	return lc.setCache(key, value)
}

func (lc *LRUCache) Set(key, value interface{}) (error) {
	if lc == nil || lc.cache == nil {
		e := errors.New("LRU Cache not initialized")
		return e
	}
	lc.lock.Lock()
	defer lc.lock.Unlock()

	if pElm, ok := lc.cache[key]; ok {
		//exist
		pElm.Value.(*cacheNode).value = value
		pElm.Value.(*cacheNode).count++
		//move to head, higher priority
		lc.pList.MoveToFront(pElm)
		return nil
	}
	//not in cache, check 2nd-level list
	for item:=lc.pList_candidate.Front(); item!=nil; item=item.Next() {
		if item.Value.(*cacheNode).key == key {
			//found
			//item.Value.(*cacheNode).value = value // no use in 2nd-level cache
			item.Value.(*cacheNode).count++
			if item.Value.(*cacheNode).count >= lc.threshold {
				//add to cache
				lc.setCache(key, value)
				//delete from 2nd-level list
				lc.pList_candidate.Remove(item)

			} else {
				// or count < threshold, higher priority in 2nd-level list
				lc.pList_candidate.MoveToFront(item)
			}
			return nil
		}
	}
	//not exist, then create cache node
	elm := newCacheNode(key, value)
	//store in 2nd-level cache list
	lc.pList_candidate.PushFront(elm)
	if lc.pList_candidate.Len() > lc.capacity {
		//delete last one
		lastElement := lc.pList_candidate.Back()
		if lastElement == nil {
			return nil
		}
		//cacheNode := lastElement.Value.(*CacheNode)
		lc.pList_candidate.Remove(lastElement)
	}

	return nil
}

func test(lc *LRUCache, k, v interface{}) {
	fmt.Println("test", k, v)
	i, r, err := lc.Get(k)
	fmt.Println("111", k)

	if r {
		if i != nil {
			fmt.Println("cache is", i.([]int))
		}
	} else {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("not exist")
		lc.Set(k, v)
	}
}

func main() {
	lc := NewLRUCache(5, 2)
	go test(lc, "wd", []int{14})
	go test(lc, "ha", []int{10})
	go test(lc, "ha", []int{11})
	go test(lc, "wd", []int{16})
	go test(lc, "ha", []int{1,2,3})
	go test(lc, "wd", []int{12})
	//test(lc, "ha", []int{1,2,3})
	//test(lc, "ha", []int{10})
	//test(lc, "ha", []int{11})

	var elc *LRUCache
	_, r, err := elc.Get(19)
	fmt.Println(r, err)
	time.Sleep(2*time.Second)
}
package crawler

import "sync"

type list struct {
	lock  sync.Mutex
	items map[string]bool
}

func (_this *list) Add(item string) {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	_this.items[item] = true
}

func (_this *list) Remove(item string) {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	delete(_this.items, item)
}

func (_this *list) IsExist(item string) bool {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	return _this.items[item]
}

func (_this *list) Items() map[string]bool {
	return _this.items
}

func (_this *list) Length() int {
	return len(_this.items)
}

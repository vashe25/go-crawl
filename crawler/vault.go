package crawler

type vault struct {
	collection *list
	visited    *list
	q          *list
}

func newVault() *vault {
	return &vault{
		collection: &list{
			items: make(map[string]bool),
		},
		visited: &list{
			items: make(map[string]bool),
		},
		q: &list{
			items: make(map[string]bool),
		},
	}
}

func (_this *vault) collect(url string) {
	_this.collection.Add(url)
}

func (_this *vault) collected() *list {
	return _this.collection
}

func (_this *vault) addVisited(url string) {
	_this.visited.Add(url)
}

func (_this *vault) isVisited(url string) bool {
	return _this.visited.IsExist(url)
}

func (_this *vault) addToQueue(url string) {
	_this.q.Add(url)
}

func (_this *vault) isInQueue(url string) bool {
	return _this.q.IsExist(url)
}

func (_this *vault) isFull() bool {
	return _this.q.Length() == _this.visited.Length()
}

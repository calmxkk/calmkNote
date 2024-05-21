package lru

type CacheLRU struct {
	limit int
	size  int
	cache map[string]*Node
	dl    *DoubleList
}

func NewLRU(limit int) *CacheLRU {
	return &CacheLRU{
		limit: limit,
		size:  0,
		cache: make(map[string]*Node),
		dl:    NewDL(),
	}
}

func (c *CacheLRU) Insert(key string, value interface{}) {
	if _, ok := c.cache[key]; !ok {
		// new node
		node := &Node{key: key, value: value}
		c.cache[key] = node
		c.dl.InsertHead(node)

		c.size++
		if c.size > c.limit {
			c.dl.RemoveLast()
			c.size--
		}
	} else {
		node := c.cache[key]
		node.value = value
		c.dl.MoveToHead(node)
	}
}

func (c *CacheLRU) Get(key string) interface{} {
	if _, ok := c.cache[key]; !ok {
		return nil
	}

	node := c.cache[key]
	c.dl.MoveToHead(node)
	return node.value
}

func TestLRU() {

}

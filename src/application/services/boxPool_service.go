package services

type BoxPool struct {
	pool chan int
}

func NewBoxPool(maxBoxes int) *BoxPool {
	ch := make(chan int, maxBoxes)
	for i := 0; i < maxBoxes; i++ {
		ch <- i
	}
	return &BoxPool{pool: ch}
}

func (bp *BoxPool) Acquire() int {
	return <-bp.pool
}

func (bp *BoxPool) Release(id int) {
	bp.pool <- id
}

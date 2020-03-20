package main

import "log"

/*视频在播放的时候要保持长连接，多路长连接的同时保持的情况下，不停的增加迟早要crash掉server，这里流控的思想就很重要
 被消耗的部分：一个是链接数，一个是带宽 ，还有一个RAM被消耗殆尽，系统就crash掉了
流控通常用的是 bucket token 算法   限定一个箱子，箱子里给20个连接，给一个连接就减一个，释放再还回来，最多20个。
通常我们在处理并发问题的时候，如果说在同时访问bucket，为了保证数据正确，要加锁，加锁就会影响性能
在go中，我们用shared channel instead of shared memory.
*/

type ConnLimiter struct {
	concurrentConn int
	bucket         chan int
}

func newConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: cc,
		bucket:         make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("reached the rate limitation.")
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <-cl.bucket
	log.Printf("new connection coming: %d", c)
}
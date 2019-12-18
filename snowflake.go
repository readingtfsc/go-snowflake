package snowflake

import (
	"sync"
	"time"
)

var node *Node

type Node struct {
	mu            sync.Mutex
	machineID     int64 // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
	sn            int64 // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64 // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
}

func init() {
	node = new(Node)
	node.lastTimeStamp = time.Now().UnixNano() / 1000000
}
func NewNode(mid int64) *Node {
	defer node.mu.Unlock()
	node.mu.Lock()
	node.machineID = mid << 12
	return node
}

func (n *Node) Snowflake() int64 {
	defer n.mu.Unlock()
	n.mu.Lock()
	curTimeStamp := time.Now().UnixNano() / 1000000
	// 同一毫秒
	if curTimeStamp == n.lastTimeStamp {
		n.sn++
		// 序列号占 12 位,十进制范围是 [ 0, 4095 ]
		if n.sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			n.lastTimeStamp = curTimeStamp
			n.sn = 0
		}

		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | n.machineID | n.sn
		return id
	}
	if curTimeStamp > n.lastTimeStamp {
		n.sn = 0
		n.lastTimeStamp = curTimeStamp
		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | n.machineID | n.sn
		return id
	}
	if curTimeStamp < n.lastTimeStamp {
		return 0
	}
	return 0
}

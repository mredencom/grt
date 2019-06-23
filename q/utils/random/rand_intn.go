package random

import (
	"crypto/rand"
	"encoding/binary"
	"os"
)

const (
	// uint32随机数的缓冲区大小。
	BUFFER_SIZE = 10000
)

var (
	// Buffer chan.
	bufferChan = make(chan uint32, BUFFER_SIZE)
)

// 它使用异步goroutine来产生随机数，
// 和缓冲区chan来存储随机数。所以它具有很高的性能
// 生成随机数。
func init() {
	step := 0
	buffer := make([]byte, 1024)
	go func() {
		for {
			if n, err := rand.Read(buffer); err != nil {
				panic(err)
				os.Exit(1)
			} else {
				// 使用缓冲区数据进行一次完整的随机数生成
				for i := 0; i < n-4; {
					bufferChan <- binary.LittleEndian.Uint32(buffer[i : i+4])
					i++
				}
				// 充分利用缓冲区数据，随机索引递增
				for i := 0; i < n; i++ {
					step = int(buffer[0]) % 10
					if step != 0 {
						break
					}
				}
				if step == 0 {
					step = 2
				}
				for i := 0; i < n-4; {
					bufferChan <- binary.BigEndian.Uint32(buffer[i : i+4])
					i += step
				}
			}
		}
	}()
}

// Intn返回一个介于0和最大值之间的int数 -  [0，max）。
// 注意：
// 1.结果大于或等于0，但小于<max>;
// 2.结果编号为32位且小于math.MaxUint32。
func Intn(max int) int {
	n := int(<-bufferChan) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}
	return n
}

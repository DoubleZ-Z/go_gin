package util

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 雪花算法结构体
type Snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

const (
	// epoch 时间偏移量 (2024-01-01 00:00:00)
	epoch int64 = 1704067200000

	// workerIDBits 工作机器ID位数
	workerIDBits uint8 = 10

	// sequenceBits 序列号位数
	sequenceBits uint8 = 12

	// maxWorkerID 最大工作机器ID
	maxWorkerID int64 = -1 ^ (-1 << workerIDBits)

	// maxSequence 最大序列号
	maxSequence int64 = -1 ^ (-1 << sequenceBits)

	// workerIDShift 工作机器ID左移位数
	workerIDShift uint8 = sequenceBits

	// timestampShift 时间戳左移位数
	timestampShift uint8 = sequenceBits + workerIDBits
)

// NewSnowflake 创建新的雪花算法实例
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}

	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

// GenerateUniqueID 生成唯一ID
func (s *Snowflake) GenerateUniqueID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTimestamp := time.Now().UnixNano() / 1000000 // 毫秒时间戳

	// 如果当前时间小于上次时间戳，说明发生了时钟回拨
	if currentTimestamp < s.timestamp {
		panic(fmt.Sprintf("clock moved backwards, refusing to generate id for %d milliseconds", s.timestamp-currentTimestamp))
	}

	// 如果是同一毫秒内生成的ID，则增加序列号
	if currentTimestamp == s.timestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		// 如果序列号达到最大值，则等待下一毫秒
		if s.sequence == 0 {
			for currentTimestamp <= s.timestamp {
				currentTimestamp = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同毫秒内，序列号从0开始
		s.sequence = 0
	}

	s.timestamp = currentTimestamp

	// 生成ID
	id := ((currentTimestamp - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// 全局雪花ID生成器实例
var globalSnowflake *Snowflake
var once sync.Once

// initGlobalSnowflake 初始化全局雪花ID生成器
func initGlobalSnowflake() {
	sf, err := NewSnowflake(1) // 使用默认worker ID为1
	if err != nil {
		panic(err)
	}
	globalSnowflake = sf
}

// GenerateUniqueID 生成唯一ID的全局函数
func GenerateUniqueID() string {
	once.Do(initGlobalSnowflake)
	id := globalSnowflake.GenerateUniqueID()
	return fmt.Sprintf("%d", id)
}

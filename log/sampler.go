/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

// Package log implements the functions, types, and interfaces for the module.
package log

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/origadmin/toolkits/errors"
)

// Sampler 提供了基于速率的日志采样功能
type Sampler struct {
	rate    float64
	counter int32 // 使用atomic操作
	mu      sync.RWMutex
	randMu  sync.Mutex // 保护随机数生成器的并发访问
	rand    *rand.Rand
}

func NewSampler(rate float64) *Sampler {
	pcg := rand.NewPCG(mustCryptoSeed())
	return &Sampler{
		rate:      rate,
		pcgSource: pcg,
		rand:      rand.New(pcg),
	}
}

// ShouldLog 根据采样率决定是否记录日志
func (s *Sampler) ShouldLog() bool {
	// 使用原子操作增加计数器
	counter := atomic.AddInt32(&s.counter, 1)
	if counter > 1000 {
		s.resetCounter()
	}

	s.randMu.Lock()
	defer s.randMu.Unlock()
	return s.rand.Float64() < s.rate
}

// resetCounter 重置计数器，线程安全
func (s *Sampler) resetCounter() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 双重检查，避免重复重置
	if s.counter > 1000 {
		s.counter = 0
		s.rand.Seed(time.Now().UnixNano()) // 使用时间戳作为种子
	}
}

type LoggerWithSampling struct {
	baseLogger log.Logger
	sampler    *Sampler
}

func (l *LoggerWithSampling) Log(level log.Level, keyvals ...any) error {
	if !l.sampler.ShouldLog() {
		return nil
	}
	return l.baseLogger.Log(level, keyvals...)
}

// LevelSampling 提供基于日志级别的采样
// 注意：初始化后不应修改rates和burstCounters，除非持有mu锁
type LevelSampling struct {
	rates         map[log.Level]float64
	burstCounters map[log.Level]int32 // 使用int32以支持原子操作
	mu            sync.RWMutex        // 保护rates和burstCounters
	randMu        sync.Mutex          // 保护随机数生成器
	rand          *rand.Rand
}

func NewLevelSampling(defaultRate float64) *LevelSampling {
	pcg := rand.NewPCG(mustCryptoSeed())
	return &LevelSampling{
		rates: map[log.Level]float64{
			log.LevelDebug: defaultRate,
			log.LevelInfo:  defaultRate,
			log.LevelWarn:  defaultRate,
			log.LevelError: 1.0,
		},
		burstCounters: make(map[log.Level]int),
		pcgSource:     pcg,
		rand:          rand.New(pcg),
	}
}

// ShouldSample 根据日志级别和采样率决定是否采样
func (ls *LevelSampling) ShouldSample(level log.Level) bool {
	// 读锁保护rates的读取
	ls.mu.RLock()
	rate, ok := ls.rates[level]
	ls.mu.RUnlock()

	if !ok {
		rate = 1.0
	}

	// 原子递增计数器
	counter := atomic.AddInt32(&ls.burstCounters[level], 1)
	if counter > 1000 {
		ls.resetCounter(level)
	}

	ls.randMu.Lock()
	defer ls.randMu.Unlock()
	return ls.rand.Float64() < rate
}

// resetCounter 重置指定级别的计数器
func (ls *LevelSampling) resetCounter(level log.Level) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	// 双重检查
	if ls.burstCounters[level] > 1000 {
		ls.burstCounters[level] = 0
		ls.rand.Seed(time.Now().UnixNano())
	}
}

func (ls *LevelSampling) GetRate(level log.Level) float64 {
	if rate, ok := ls.rates[level]; ok {
		return rate
	}
	return 1.0
}

type LevelSampler struct {
	logger  log.Logger
	sampler *LevelSampling
}

func (l *LevelSampler) Log(level log.Level, keyvals ...any) error {
	if !l.sampler.ShouldSample(level) {
		return nil
	}
	return l.logger.Log(level, keyvals...)
}

func NewLevelSampler(logger log.Logger, sampling *LevelSampling) log.Logger {
	return &LevelSampler{
		logger:  logger,
		sampler: sampling,
	}
}

var seedPool = sync.Pool{
	New: func() any {
		return new([16]byte)
	},
}

func cryptoSeed() (uint64, uint64, error) {
	buf := seedPool.Get().(*[16]byte)
	defer func() {
		clear(buf[:])
		seedPool.Put(buf)
	}()

	if _, err := cryptorand.Read(buf[:]); err != nil {
		return 0, 0, errors.Wrapf(err, "crypto/rand failure")
	}

	return binary.BigEndian.Uint64(buf[0:8]),
		binary.BigEndian.Uint64(buf[8:16]), nil
}

// mustCryptoSeed 生成加密安全的随机种子
// 如果无法生成安全的随机数，会使用当前时间戳作为备选方案
func mustCryptoSeed() (uint64, uint64) {
	h, l, err := cryptoSeed()
	if err != nil {
		// 使用时间戳作为备选方案，而不是panic
		nano := time.Now().UnixNano()
		return uint64(nano), uint64(nano >> 32)
	}
	return h, l
}

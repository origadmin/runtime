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

type Sampler struct {
	rate      float64
	counter   int
	pcgSource *rand.PCG
	mu        sync.Mutex
}

func NewSampler(rate float64) *Sampler {
	return &Sampler{
		rate:      rate,
		pcgSource: rand.NewPCG(mustCryptoSeed()),
	}
}

func (s *Sampler) ShouldLog() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	if s.counter > 1000 {
		s.counter = 0
		s.pcgSource.Seed(mustCryptoSeed())
	}
	return rand.New(s.pcgSource).Float64() < s.rate
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

type LevelSampling struct {
	rates         map[log.Level]float64
	burstCounters map[log.Level]int
	pcgSource     *rand.PCG
	mu            sync.Mutex
}

func NewLevelSampling(defaultRate float64) *LevelSampling {
	return &LevelSampling{
		rates: map[log.Level]float64{
			log.LevelDebug: defaultRate,
			log.LevelInfo:  defaultRate,
			log.LevelWarn:  defaultRate,
			log.LevelError: 1.0,
		},
		burstCounters: make(map[log.Level]int),
		pcgSource:     rand.NewPCG(mustCryptoSeed()),
	}
}

func (ls *LevelSampling) ShouldSample(level log.Level) bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	rate, ok := ls.rates[level]
	if !ok {
		rate = 1.0
	}

	if ls.burstCounters[level] > 1000 {
		ls.burstCounters[level] = 0
		ls.pcgSource.Seed(mustCryptoSeed())
	}
	ls.burstCounters[level]++
	return rand.New(ls.pcgSource).Float64() < rate
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
	rate := l.sampler.GetRate(level)
	if rand.Float64() > rate {
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
		return 0, 0, errors.Wrap(err, "crypto/rand failure")
	}

	return binary.BigEndian.Uint64(buf[0:8]),
		binary.BigEndian.Uint64(buf[8:16]), nil
}

func mustCryptoSeed() (uint64, uint64) {
	h, l, err := cryptoSeed()
	if err != nil {
		now := time.Now().UnixNano()
		h, l = uint64(now), uint64(now>>32)
	}
	return h, l
}

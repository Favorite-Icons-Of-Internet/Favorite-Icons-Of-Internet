package main

import (
	"expvar"
	"net/http"
)

const (
	MonitorBind = "localhost:8855"
)

type selfMonitor struct {
	processedItems *expvar.Int
	failedItems    *expvar.Int
}

func NewSelfMonitor() *selfMonitor {
	return &selfMonitor{
		processedItems: expvar.NewInt("ProcessedItems"),
		failedItems:    expvar.NewInt("FailedItems"),
	}
}

func (sm selfMonitor) Start() {
	go http.ListenAndServe(MonitorBind, nil)
}

func (sm selfMonitor) AddProcessed() {
	sm.processedItems.Add(1)
}

func (sm selfMonitor) AddFailed() {
	sm.failedItems.Add(1)
}

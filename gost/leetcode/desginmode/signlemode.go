package main

import "sync"

type Admin struct{}

var instance *Admin
var once sync.Once

func getInstance1() *Admin {
	once.Do(func() {
		instance = &Admin{}
	})
	return instance
}

func getInstance2() *Admin {
	if instance == nil {
		instance = &Admin{}
	}
	return instance
}

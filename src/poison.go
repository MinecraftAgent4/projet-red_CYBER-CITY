package main

import "time"

func (m *Monster) poison() {
	for i := 0; i <= 2; i++ {
		m.AddPV(-10)
		time.Sleep(1 * time.Second)
	}
}

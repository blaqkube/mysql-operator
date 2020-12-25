package controllers

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Store Controller", func() {
	It("Check NewTimeManager has values", func() {
		t := NewTimeManager()
		Expect(t).NotTo(BeNil())
		Expect((*t)[0].Round(1 * time.Second).Seconds()).To(Equal(float64(1)))
	})

	It("Check TimeManager Next for now is 1s", func() {
		t := NewTimeManager().Next(time.Now())
		Expect(t).To(Equal(1 * time.Second))
	})

	It("Check TimeManager Next for initial time is 300s", func() {
		t := NewTimeManager().Next(time.Time{})
		Expect(t).To(Equal(300 * time.Second))
	})

	It("Check TimeManager Next for 10s before now should be 5s", func() {
		t := NewTimeManager().Next(time.Now().Add(-10 * time.Second))
		Expect(t).To(Equal(5 * time.Second))
	})

	It("Check TimeManager Next for 60s before now should be 20s", func() {
		t := NewTimeManager().Next(time.Now().Add(-60 * time.Second))
		Expect(t).To(Equal(20 * time.Second))
	})
})

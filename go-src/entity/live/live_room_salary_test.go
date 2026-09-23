package entity

import (
	"testing"
	"time"
)

func TestLiveRoomIsSalaryEffective(t *testing.T) {
	startsAt := time.Date(2026, time.September, 23, 10, 0, 0, 0, time.Local)
	endsAt := time.Date(2026, time.September, 28, 0, 0, 0, 0, time.Local)
	room := &LiveRoom{
		HasSalary:                true,
		SalaryEffectiveStartTime: &startsAt,
		SalaryEffectiveEndTime:   &endsAt,
	}

	if room.IsSalaryEffective(startsAt.Add(-time.Nanosecond)) {
		t.Fatal("salary should not be effective before the start boundary")
	}
	if !room.IsSalaryEffective(startsAt) {
		t.Fatal("salary should be effective at the start boundary")
	}
	if !room.IsSalaryEffective(endsAt.Add(-time.Nanosecond)) {
		t.Fatal("salary should remain effective before the end boundary")
	}
	if room.IsSalaryEffective(endsAt) {
		t.Fatal("salary should expire at next Monday 00:00:00")
	}
	if room.IsSalaryEffective(endsAt.Add(time.Second)) {
		t.Fatal("salary should remain expired after the expiration boundary")
	}

	room.HasSalary = false
	if room.IsSalaryEffective(endsAt.Add(-time.Second)) {
		t.Fatal("salary should not be effective when has_salary is false")
	}
}

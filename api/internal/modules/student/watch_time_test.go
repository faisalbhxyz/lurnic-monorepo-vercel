package student

import (
	"testing"
	"time"
)

func TestNormalizeWatchInput_ClampsAndDefaults(t *testing.T) {
	svc := &StorefrontService{}
	watchedAt := "2026-09-04T08:15:22Z"
	platform := "android"
	event, err := svc.normalizeWatchInput(WatchTimeInput{
		ClientEventID:  "550e8400-e29b-41d4-a716-446655440000",
		WatchedSeconds: 400,
		WatchDate:      "2026-09-04",
		Timezone:       "Asia/Dhaka",
		WatchedAt:      &watchedAt,
		DevicePlatform: &platform,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event.WatchedSeconds != maxWatchedSecondsPerEvent {
		t.Fatalf("expected clamp to %d, got %d", maxWatchedSecondsPerEvent, event.WatchedSeconds)
	}
	if event.Source != watchSourceEnrolled {
		t.Fatalf("expected default source enrolled, got %s", event.Source)
	}
	if event.WatchDateStr != "2026-09-04" {
		t.Fatalf("watch date = %s", event.WatchDateStr)
	}
}

func TestNormalizeWatchInput_RejectsZeroSeconds(t *testing.T) {
	svc := &StorefrontService{}
	_, err := svc.normalizeWatchInput(WatchTimeInput{
		ClientEventID:  "abc",
		WatchedSeconds: 0,
		WatchDate:      "2026-09-04",
		Timezone:       "Asia/Dhaka",
	})
	if err != errInvalidWatchSeconds {
		t.Fatalf("expected errInvalidWatchSeconds, got %v", err)
	}
}

func TestNormalizeWatchInput_RejectsFutureDate(t *testing.T) {
	svc := &StorefrontService{}
	future := time.Now().In(reportLocation()).AddDate(0, 0, 3).Format("2006-01-02")
	_, err := svc.normalizeWatchInput(WatchTimeInput{
		ClientEventID:  "abc",
		WatchedSeconds: 10,
		WatchDate:      future,
		Timezone:       "Asia/Dhaka",
	})
	if err != errFutureWatchDate {
		t.Fatalf("expected errFutureWatchDate, got %v", err)
	}
}

func TestCalcStreakDaysEndingToday(t *testing.T) {
	loc := reportLocation()
	today := startOfLocalDay(time.Now(), loc)
	daily := []DailyWatchSecondsEntry{
		{Date: today.AddDate(0, 0, -2).Format("2006-01-02"), Seconds: 100},
		{Date: today.AddDate(0, 0, -1).Format("2006-01-02"), Seconds: 200},
		{Date: today.Format("2006-01-02"), Seconds: 50},
	}
	if got := calcStreakDaysEndingToday(daily, today, loc); got != 3 {
		t.Fatalf("streak = %d, want 3", got)
	}

	daily[2].Seconds = 0
	if got := calcStreakDaysEndingToday(daily, today, loc); got != 0 {
		t.Fatalf("streak should break on today empty, got %d", got)
	}
}

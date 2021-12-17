package ztime

import (
	"testing"
	"time"
)

func TestRangeCurrent(t *testing.T) {
	tt := time.Date(2020, 6, 18, 14, 49, 20, 666, time.UTC)
	ns := "999999999"

	tests := []struct {
		t                  time.Time
		p                  Period
		wantStart, wantEnd time.Time
	}{
		{tt, Minute,
			New("2020-06-18 14:49:00.0"),
			New("2020-06-18 14:49:59." + ns),
		},
		{tt, Hour,
			New("2020-06-18 14:00:00.0"),
			New("2020-06-18 14:59:59." + ns),
		},
		{tt, Day,
			New("2020-06-18 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, WeekMonday,
			New("2020-06-15 00:00:00.0"),
			New("2020-06-21 23:59:59." + ns),
		},
		{tt, WeekSunday,
			New("2020-06-14 00:00:00.0"),
			New("2020-06-20 23:59:59." + ns),
		},
		{tt, Month,
			New("2020-06-01 00:00:00.0"),
			New("2020-06-30 23:59:59." + ns),
		},
		{tt, Quarter,
			New("2020-04-01 00:00:00.0"),
			New("2020-06-30 23:59:59." + ns),
		},
		{tt, HalfYear,
			New("2020-01-01 00:00:00.0"),
			New("2020-06-30 23:59:59." + ns),
		},
		{tt, Year,
			New("2020-01-01 00:00:00.0"),
			New("2020-12-31 23:59:59." + ns),
		},

		{time.Date(2020, 1, 1, 14, 49, 20, 666, time.UTC), WeekMonday,
			New("2019-12-30 00:00:00.0"),
			New("2020-01-05 23:59:59." + ns),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := NewRange(tt.t).Current(tt.p)

			if !have.Start.Equal(tt.wantStart) {
				t.Errorf("start wrong\nhave: %s\nwant: %s", have.Start, tt.wantStart)
			}
			if !have.End.Equal(tt.wantEnd) {
				t.Errorf("end wrong\nhave: %s\nwant: %s", have.End, tt.wantEnd)
			}
		})
	}
}

func TestRangeLast(t *testing.T) {
	tt := time.Date(2020, 6, 18, 14, 49, 20, 666, time.UTC)
	ns := "999999999"

	tests := []struct {
		t                  time.Time
		p                  Period
		wantStart, wantEnd time.Time
	}{
		{tt, Minute,
			New("2020-06-18 14:48:20.0"),
			New("2020-06-18 14:49:20." + ns),
		},
		{tt, Hour,
			New("2020-06-18 13:49:00.0"),
			New("2020-06-18 14:49:59." + ns),
		},
		{tt, Day,
			New("2020-06-17 14:00:00.0"),
			New("2020-06-18 14:59:59." + ns),
		},
		{tt, WeekMonday,
			New("2020-06-11 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, WeekSunday,
			New("2020-06-11 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, Month,
			New("2020-05-18 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, Quarter,
			New("2020-03-18 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, HalfYear,
			New("2019-12-18 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
		{tt, Year,
			New("2019-06-18 00:00:00.0"),
			New("2020-06-18 23:59:59." + ns),
		},
	}

	for _, tt := range tests {
		t.Run(tt.p.String(), func(t *testing.T) {
			have := NewRange(tt.t).Last(tt.p)

			if !have.Start.Equal(tt.wantStart) {
				t.Errorf("start wrong\nhave: %s\nwant: %s", have.Start, tt.wantStart)
			}
			if !have.End.Equal(tt.wantEnd) {
				t.Errorf("end wrong\nhave: %s\nwant: %s", have.End, tt.wantEnd)
			}
		})
	}
}

func TestRangeString(t *testing.T) {
	tests := []struct {
		now, start, end, want string
	}{
		// One day
		{
			"2020-06-18",
			"2020-06-18",
			"2020-06-18",
			"Today",
		},
		{
			"2020-06-18",
			"2020-06-17",
			"2020-06-17",
			"Yesterday",
		},
		{
			"2020-06-18",
			"2020-06-16",
			"2020-06-16",
			"Jun 16",
		},
		{
			"2020-06-18",
			"2019-06-16",
			"2019-06-16",
			"Jun 16 2019",
		},

		// Multiple days, less than a week.
		{
			"2020-06-18",
			"2020-06-17",
			"2020-06-18",
			"Yesterday–Today",
		},
		{
			"2020-06-18",
			"2020-06-12",
			"2020-06-18",
			"6 days ago–Today",
		},
		{
			"2020-06-18",
			"2020-06-16",
			"2020-06-17",
			"Jun 16–Jun 17 (2 days)",
		},
		{
			"2020-06-18",
			"2020-06-15",
			"2020-06-17",
			"Jun 15–Jun 17 (3 days)",
		},

		// One week
		{
			"2021-04-05",
			"2021-04-05",
			"2021-04-11",
			"Apr 5–Apr 11 (1 week)",
		},
		{
			"2021-04-05",
			"2021-03-08",
			"2021-03-14",
			"Mar 8–Mar 14 (1 week)",
		},

		// More than a week.
		{
			"2020-06-18",
			"2020-06-27",
			"2020-07-04",
			"Jun 27–Jul 4 (8 days)",
		},

		{
			"2020-06-18",
			"2020-06-05",
			"2020-06-18",
			"13 days ago–Today",
		},
		{
			"2020-06-18",
			"2020-06-04",
			"2020-06-18",
			"2 weeks ago–Today",
		},
		{
			"2020-06-18",
			"2020-06-03",
			"2020-06-18",
			"15 days ago–Today",
		},
		{
			"2020-06-18",
			"2020-05-18",
			"2020-06-18",
			"1 month ago–Today",
		},
		{
			"2020-08-29",
			"2020-05-28",
			"2020-08-29",
			"May 28–Today",
		},
		{
			"2020-06-18",
			"2019-05-18",
			"2020-06-18",
			"13 months ago–Today",
		},

		{
			"2021-04-05",
			"2020-10-05",
			"2021-04-05",
			"6 months ago–Today",
		},

		{
			"2020-06-18",
			"2020-06-01",
			"2020-06-30",
			"June",
		},
		{
			"2020-06-18",
			"2019-06-01",
			"2019-06-30",
			"June 2019",
		},

		{
			"2020-06-18",
			"2020-02-01",
			"2020-06-30",
			"February–June",
		},
		{
			"2020-06-18",
			"2019-02-01",
			"2019-06-30",
			"February 2019–June 2019",
		},

		{
			"2020-06-18",
			"2020-04-02",
			"2020-06-17",
			"Apr 2–Jun 17 (2 months, 15 days)",
		},

		{
			"2020-06-18",
			"2019-01-01",
			"2020-07-22",
			"Jan 1 2019–Jul 22 (18 months, 21 days)",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			SetNow(t, tt.now)

			have := NewRange(New(tt.start)).To(New(tt.end)).String()
			if have != tt.want {
				t.Errorf("\nhave: %q\nwant: %q", have, tt.want)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		a, b time.Time
		p    []Period
		want Diff
	}{
		{time.Time{}, time.Time{}, nil, Diff{}},
		{New("1548-09-18 17:44:55"), New("1548-09-18 17:44:55"), nil, Diff{}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{99}, Diff{}}, // Invalid value

		{New("1548-06-10"), New("1548-06-18"), nil, Diff{Days: 8}},
		{New("1542-06-10 13:04:05"), New("1548-09-18 17:44:55"), nil, Diff{Years: 6, Months: 3, Days: 8, Hours: 4, Mins: 40, Secs: 50}},

		{time.Date(2020, -1000, -1000, -1000, -1000, -1000, -1000, time.UTC), New("1548-09-18 17:44:55"), nil,
			Diff{Years: 385, Days: 3, Hours: 21, Mins: 18, Secs: 24}},
		{New("1548-09-18 17:44:55").AddDate(0, 1, 1), New("1548-09-18 17:44:55"), nil, Diff{Months: 1, Days: 1}},
		{New("1548-09-18 17:44:55").AddDate(0, -1, -1), New("1548-09-18 17:44:55"), nil, Diff{Months: 1, Days: 1}},
		{New("1548-09-18 17:44:55").AddDate(0, -13, -13), New("1548-09-18 17:44:55"), nil, Diff{Years: 1, Months: 1, Days: 13}},
		{New("1548-09-18 17:44:55").AddDate(0, -26, -13), New("1548-09-18 17:44:55"), nil, Diff{Years: 2, Months: 2, Days: 13}},

		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), nil, Diff{Years: 1, Months: 1, Days: 8}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{Day}, Diff{Days: 405}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{Month}, Diff{Months: 13}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{Year}, Diff{Years: 1}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{Month, Day}, Diff{Months: 13, Days: 8}},

		// TODO: make sure this is correct.
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{WeekMonday}, Diff{Weeks: 57}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{WeekMonday, Day}, Diff{Weeks: 57, Days: 6}},
		{New("2020-06-18 14:15:16"), New("2019-05-10 14:15:16"), []Period{Month, WeekMonday, Day}, Diff{Months: 13, Weeks: 1, Days: 1}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := NewRange(tt.a).To(tt.b).Diff(tt.p...).String()
			want := tt.want.String()
			if have != want {
				t.Fatalf("\nhave: %v\nwant: %v", have, want)
			}
		})
	}
}

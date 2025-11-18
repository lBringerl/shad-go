//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	loadChangingDates := make(map[int]struct{})
	hotelState := make(map[int]int)
	for _, guest := range guests {
		loadChangingDates[guest.CheckInDate] = struct{}{}
		loadChangingDates[guest.CheckOutDate] = struct{}{}
		hotelState[guest.CheckInDate] += 1
		hotelState[guest.CheckOutDate] -= 1
	}

	dates := make([]int, 0, len(loadChangingDates))
	for date := range loadChangingDates {
		dates = append(dates, date)
	}
	sort.Ints(dates)

	load := make([]Load, 0)

	currentState := 0
	for _, date := range dates {
		if hotelState[date] == 0 {
			continue
		}

		currentState += hotelState[date]
		load = append(load, Load{
			StartDate:  date,
			GuestCount: currentState,
		})
	}

	return load
}

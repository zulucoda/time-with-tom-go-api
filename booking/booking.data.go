package booking

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

// Creating RW Mutex so that goroutines for write data at the same time
var bookingMap = struct {
	sync.RWMutex
	m map[int]Booking
}{m: make(map[int]Booking)}

func init() {
	fmt.Println("loading bookings...")
	bookMap, err := loadBookingMap()
	bookingMap.m = bookMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf("%d bookings loaded...\n", len(bookingMap.m))
}

func loadBookingMap() (map[int]Booking, error) {
	fileName := "bookings.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist")
	}

	file, _ := ioutil.ReadAll(fileName)
	bookingList := make([]Booking, 0)
	err = json.Unmarshal([]byte(file), &bookingList)
	if err != nil {
		log.Fatal(err)
	}
	bookMap := make(map[int]Booking)
	for i := 0; i < len(bookingList); i++ {
		bookMap[bookingList[i].BookingID] = bookingList[i]
	}
	return bookMap, nil
}

func getBooking(bookingID int) *Booking {
	bookingMap.RLock()
	defer bookingMap.RUnlock()
	if booking, ok := bookingMap.m[bookingID]; ok {
		return &booking
	}
	return nil
}

func removeBooking(bookingID int) {
	bookingMap.Lock()
	defer bookingMap.Unlock()
	delete(bookingMap.m, bookingID)
}

func getBookingList() []Booking {
	bookingMap.RLock()
	bookings := make([]Booking, 0, len(bookingMap.m))
	for _, value := range bookingMap.m {
		bookings = append(bookings, value)
	}
	bookingMap.Unlock()
	return bookings
}

func getBookingIds() []int {
	bookingMap.RLock()
	bookingIds := []int{}
	for key := range bookingMap.m {
		bookingIds = append(bookingIds, key)
	}
	bookingMap.RUnlock()
	sort.Ints(bookingIds)
	return bookingIds
}

func getNextBookingID() int {
	bookingIDs := getBookingIds()
	return bookingIDs[len(bookingIDs)-1] + 1
}

func addOrUpdateBooking(booking Booking) (int, error) {
	// if booking id is set, update, otherwise add
	addOrUpdate := -1
	if booking.BookingID > 0 {
		oldBooking := getBooking(booking.BookingID)
		// if it exist, replace it, otherwise return error
		if oldBooking == nil {
			return 0, fmt.Errorf("booking id [%d] doesnt exist", booking.BookingID)
		}
		addOrUpdate = booking.BookingID
	} else {
		addOrUpdate = getNextBookingID()
		booking.BookingID = addOrUpdate
	}

	bookingMap.Lock()
	bookingMap.m[addOrUpdate] = booking
	bookingMap.Unlock()
	return addOrUpdate, nil
}

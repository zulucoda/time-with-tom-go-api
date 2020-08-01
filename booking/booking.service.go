package booking

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bookingsJson, err := json.Marshal(bookingList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bookingsJson)
		return

	case http.MethodPost:
		var newBooking Booking
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &newBooking)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newBooking.BookingID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newBooking.BookingID = getNextId()
		bookingList = append(bookingList, newBooking)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {

	urlPathSegments := strings.Split(r.URL.Path, "bookings/")
	bookingID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	booking, listItemIndex := findBookingByID(bookingID)
	if booking == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		bookingJSON, err := json.Marshal(booking)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bookingJSON)
		return

	case http.MethodPut:
		var updateBooking Booking
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updateBooking)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateBooking.BookingID != bookingID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		booking = &updateBooking
		bookingList[listItemIndex] = *booking
		w.WriteHeader(http.StatusOK)
		return

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
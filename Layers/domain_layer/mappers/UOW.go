package mappers

import (
	"fmt"
	"github.com/Skellyboy38/SOEN-343-NullPointer/Layers/domain_layer/classes"
)

type userQueue []classes.User
type reservationQueue []classes.Reservation

var UOWSingleTon *UOW

type UOW struct {
	registeredNewUsers            userQueue
	registeredDirtyUsers          userQueue
	registeredDeletedUsers        userQueue
	registeredNewReservations     reservationQueue
	registeredDirtyReservations   reservationQueue
	registeredDeletedReservations []int
	userMapper                    *UserMapper
	ReservationMapper             *ReservationMapper
}

func InitUOW() {
	UOWSingleTon = &UOW{
		[]classes.User{},
		[]classes.User{},
		[]classes.User{},
		[]classes.Reservation{},
		[]classes.Reservation{},
		[]int{},
		MapperBundle.UserMapper,
		MapperBundle.ReservationMapper,
	}
}

func (uow *UOW) RegisterNewUser(object classes.User) {
	uow.registeredNewUsers = append(uow.registeredNewUsers, object)
	fmt.Println(uow.registeredNewUsers)
}

func (uow *UOW) RegisterNewReservation(object classes.Reservation) {
	uow.registeredNewReservations = append(uow.registeredNewReservations, object)
}

func (uow *UOW) RegisterDirtyUser(object classes.User) {
	uow.registeredDirtyUsers = append(uow.registeredDirtyUsers, object)
}

func (uow *UOW) RegisterDirtyReservations(object classes.Reservation) {
	uow.registeredDirtyReservations = append(uow.registeredDirtyReservations, object)
}

func (uow *UOW) RegisterDeleteReservation(id int) {
	uow.registeredDeletedReservations = append(uow.registeredDeletedReservations, id)
}

func (uow *UOW) Commit() error {
	fmt.Println("GOT TO COMMIT")
	fmt.Println(uow.registeredNewUsers)

	processedRegisteredNewUsers := reverseUsers(reduceUserQueue(reverseUsers(uow.registeredNewUsers)))
	processedRegisteredDirtyUsers := reverseUsers(reduceUserQueue(reverseUsers(uow.registeredDirtyUsers)))
	processedRegisteredDeletedUsers := convertToUserIdSlice(
		reverseUsers(
			reduceUserQueue(
				reverseUsers(
					uow.registeredDeletedUsers))))
	MapperBundle.UserMapper.SaveDeleted(processedRegisteredDeletedUsers)
	MapperBundle.UserMapper.SaveDirty(processedRegisteredDirtyUsers)
	MapperBundle.UserMapper.SaveNew(processedRegisteredNewUsers)

	processedRegisteredNewReservations := reverseReservations(reduceReservationQueue(reverseReservations(uow.registeredNewReservations)))
	processedRegisteredDirtyReservations := reverseReservations(reduceReservationQueue(reverseReservations(uow.registeredDirtyReservations)))
	processedRegisteredDeletedReservations := reverseIntArray(
		reduceIntQueue(
			reverseIntArray(
				uow.registeredDeletedReservations)))

	if err := uow.ReservationMapper.SaveDeleted(processedRegisteredDeletedReservations); err != nil {
		fmt.Println(err)
		return err
	}
	if err := uow.ReservationMapper.SaveDirty(processedRegisteredDirtyReservations); err != nil {
		fmt.Println(err)
		return err
	}
	if err := uow.ReservationMapper.SaveNew(processedRegisteredNewReservations); err != nil {
		fmt.Println(err)
		return err
	}

	uow.registeredNewUsers = userQueue{}
	uow.registeredDirtyUsers = userQueue{}
	uow.registeredDeletedUsers = userQueue{}
	uow.registeredNewReservations = reservationQueue{}
	uow.registeredDirtyReservations = reservationQueue{}
	uow.registeredDeletedReservations = []int{}

	return nil
}

func reverseUsers(users []classes.User) []classes.User {
	reversedUsers := []classes.User{}
	for i := len(users) - 1; i >= 0; i-- {
		reversedUsers = append(reversedUsers, users[i])
	}
	return reversedUsers
}

func reverseReservations(reservations []classes.Reservation) []classes.Reservation {
	reversedReservations := []classes.Reservation{}
	for i := len(reservations) - 1; i >= 0; i-- {
		reversedReservations = append(reversedReservations, reservations[i])
	}
	return reversedReservations
}

func reverseIntArray(reservations []int) []int {
	reversedReservations := []int{}
	for i := len(reservations) - 1; i >= 0; i-- {
		reversedReservations = append(reversedReservations, reservations[i])
	}
	return reversedReservations
}

func reduceUserQueue(queue []classes.User) userQueue {
	reducedQueue := []classes.User{}
	exist := make(map[int]classes.User)
	for _, element := range queue {
		_, found := exist[element.StudentId]
		if found {
			continue
		} else {
			reducedQueue = append(reducedQueue, element)
		}
	}
	return reducedQueue
}

func reduceReservationQueue(queue []classes.Reservation) reservationQueue {
	reducedQueue := []classes.Reservation{}
	exist := make(map[int]classes.Reservation)
	for _, element := range queue {
		_, found := exist[element.ReservationId]
		if found {
			continue
		} else {
			reducedQueue = append(reducedQueue, element)
		}
	}
	return reducedQueue
}

func reduceIntQueue(queue []int) []int {
	reducedQueue := []int{}
	exist := make(map[int]int)
	for _, element := range queue {
		_, found := exist[element]
		if found {
			continue
		} else {
			reducedQueue = append(reducedQueue, element)
		}
	}
	return reducedQueue
}

func convertToUserIdSlice(userSlice []classes.User) []int {
	intSlice := []int{}
	for _, x := range userSlice {
		intSlice = append(intSlice, x.StudentId)
	}
	return intSlice
}

package db

const MongoDBNmeEnvName = "MONGO_DB_NAME"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type HotelPaginationQueryFilter struct {
	Page  int64
	Limit int64
}

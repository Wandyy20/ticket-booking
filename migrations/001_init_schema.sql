CREATE TABLE flights (
    id VARCHAR(50) PRIMARY KEY,
    airline VARCHAR(100) NOT NULL,
    origin VARCHAR(100) NOT NULL,
    destination VARCHAR(100) NOT NULL,
    departure_time TIMESTAMP NOT NULL,
    arrival_time TIMESTAMP NOT NULL,
    price NUMERIC NOT NULL,
    total_seats INT NOT NULL
);

CREATE TABLE seats (
    id VARCHAR(100) PRIMARY KEY,
    flight_id VARCHAR(50) NOT NULL,
    seat_number VARCHAR(20) NOT NULL,
    is_booked BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_seats_flight
        FOREIGN KEY (flight_id)
        REFERENCES flights(id)
);

CREATE TABLE bookings (
    id VARCHAR(100) PRIMARY KEY,
    flight_id VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP,

    CONSTRAINT fk_bookings_flight
        FOREIGN KEY (flight_id)
        REFERENCES flights(id)
);

CREATE TABLE booking_seats (
    id VARCHAR(100) PRIMARY KEY,
    booking_id VARCHAR(100) NOT NULL,
    seat_id VARCHAR(100) NOT NULL,
    passenger_name VARCHAR(100) NOT NULL,

    CONSTRAINT fk_booking_seats_booking
        FOREIGN KEY (booking_id)
        REFERENCES bookings(id),

    CONSTRAINT fk_booking_seats_seat
        FOREIGN KEY (seat_id)
        REFERENCES seats(id)
);
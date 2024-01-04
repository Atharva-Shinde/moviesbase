CREATE TABLE IF NOT EXISTS movies(
	id bigserial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- "timestamp with time zone" is a datatype to display both date and time along with time zone
	title text NOT NULL,
	year integer NOT NULL,
	genres text[] NOT NULL,
	runtime integer NOT NULL,
	version integer NOT NULL DEFAULT 1
);
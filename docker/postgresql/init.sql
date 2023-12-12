/*
create table at init
*/
CREATE TABLE users (
  user_id serial PRIMARY KEY,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  date_of_birth TIMESTAMP,
  city VARCHAR(50),
  street_name VARCHAR(50),
  street_address VARCHAR(50),
  zip_code VARCHAR(50),
  state VARCHAR(50),
  country VARCHAR(50),
  latitude double precision,
  longitude double precision
);
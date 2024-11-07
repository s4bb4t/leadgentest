-- +migrate Up
CREATE TABLE IF NOT EXISTS public.buildings (
id SERIAL PRIMARY KEY,
title VARCHAR(255) NOT NULL,
city VARCHAR(255) NOT NULL,
year INT NOT NULL,
floors INT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS public.buildings;

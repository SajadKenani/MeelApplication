BEGIN;

CREATE TABLE public.cars (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INT NOT NULL,
    description TEXT NOT NULL,
    brand TEXT NOT NULL,
    image BYTEA NOT NULL,
    more_images_id TEXT[] NOT NULL,
    _more_images_id TEXT NOT NULL,
    property TEXT[] NOT NULL,
);

ALTER TABLE public.cars add column info text null
ALTER TABLE public.cars add column is_favorite bool null

CREATE TABLE public.cars_images (
    id SERIAL PRIMARY KEY,
    image BYTEA NOT NULL
);

CREATE TABLE public.brands (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image BYTEA NOT NULL
);

CREATE TABLE public.property (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
)

CREATE TABLE public.info {
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL

}

CREATE TABLE public.sub_property (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    property_id INT NOT NULL REFERENCES public.property(id) ON DELETE CASCADE 
)

COMMIT;
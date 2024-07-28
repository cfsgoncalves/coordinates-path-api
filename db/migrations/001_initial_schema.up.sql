-- public.orders definition

CREATE TABLE public.orders (
	id bigserial NOT NULL,
	weight numeric NULL,
	latitude numeric NULL,
	longitude numeric NULL,
	description text NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id)
);

-- public.trucks definition

CREATE TABLE public.trucks (
	plate text NOT NULL,
	max_weight numeric NULL,
	CONSTRAINT trucks_pkey PRIMARY KEY (plate)
);

-- public.order_trucks definition

CREATE TABLE public.order_trucks (
	"date" text NOT NULL,
	order_id bigserial NOT NULL,
	truck_plate text NULL,
	order_sequence int8 NULL,
	CONSTRAINT order_trucks_pkey PRIMARY KEY (date, order_id, truck_plate)
);

-- public.order_trucks foreign keys

ALTER TABLE public.order_trucks ADD CONSTRAINT fk_order_trucks_order FOREIGN KEY (order_id) REFERENCES public.orders(id);
ALTER TABLE public.order_trucks ADD CONSTRAINT fk_order_trucks_truck FOREIGN KEY (truck_plate) REFERENCES public.trucks(plate);


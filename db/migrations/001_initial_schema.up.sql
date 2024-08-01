-- public.orders definition

CREATE TABLE public.orders (
	order_code text NOT NULL,
	weight float NOT NULL,
	latitude float NOT NULL,
	longitude float NOT NULL,
	description text NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (order_code)
);

-- public.trucks definition

CREATE TABLE public.trucks (
	plate text NOT NULL,
	max_weight float NOT NULL,
	CONSTRAINT trucks_pkey PRIMARY KEY (plate)
);

-- public.order_trucks definition

CREATE TABLE public.order_trucks (
	"date" text NOT NULL,
	order_code text NOT NULL,
	truck_plate text NOT NULL,
	order_sequence int NULL,
	order_status text NOT NULL, 
	CONSTRAINT order_trucks_pkey PRIMARY KEY (date, order_code, truck_plate)
);

-- public.order_trucks foreign keys

ALTER TABLE public.order_trucks ADD CONSTRAINT fk_order_trucks_order FOREIGN KEY (order_code) REFERENCES public.orders(order_code);
ALTER TABLE public.order_trucks ADD CONSTRAINT fk_order_trucks_truck FOREIGN KEY (truck_plate) REFERENCES public.trucks(plate);


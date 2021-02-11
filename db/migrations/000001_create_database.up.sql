
CREATE TABLE if not exists public.device
(
    id serial PRIMARY KEY,
    name text,
    address text,
    location text,
    is_active BOOLEAN DEFAULT true
);
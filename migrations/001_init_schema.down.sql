-- Drop order_items first because it depends on orders and products
DROP TABLE IF EXISTS order_items;

-- Drop orders next because it depends on customers
DROP TABLE IF EXISTS orders;

-- Drop customers
DROP TABLE IF EXISTS customers;

-- Drop products (depends on categories)
DROP TABLE IF EXISTS products;

-- Drop categories (self-referencing)
DROP TABLE IF EXISTS categories;

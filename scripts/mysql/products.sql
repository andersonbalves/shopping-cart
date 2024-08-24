
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

INSERT INTO products (name, price) VALUES 
  ('Notebook', 5000.00),
  ('Celular', 3054.59),
  ('Tablet', 2500.29),
  ('Fone de ouvido', 505.19);
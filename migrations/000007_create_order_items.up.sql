
CREATE TABLE IF NOT EXISTS order_items (
  order_id INT NOT NULL,
  item_id INT NOT NULL,
  PRIMARY KEY(order_id,item_id),
  quantity INT NOT NULL,
  instruction VARCHAR(6000),
  complete BOOL DEFAULT FALSE,
  FOREIGN KEY(order_id) REFERENCES orders(id),
  FOREIGN KEY(item_id) REFERENCES items(id)
);

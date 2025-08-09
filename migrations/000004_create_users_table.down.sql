DELETE FROM items WHERE name IN (
  'Garlic Bread', 'Bruschetta', 'Stuffed Mushrooms',
  'Grilled Chicken Alfredo', 'Beef Lasagna', 'Veggie Stir Fry',
  'Cheesecake', 'Chocolate Lava Cake', 'Fruit Salad',
  'Coke', 'Orange Juice', 'Coffee'
);

DELETE FROM categories WHERE name IN (
  'Appetizers', 'Main Courses', 'Desserts', 'Beverages'
);


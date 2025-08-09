INSERT INTO categories (name) VALUES
('Appetizers'),
('Main Courses'),
('Desserts'),
('Beverages');
INSERT INTO items (category_id, name, price, description) VALUES
(1, 'Garlic Bread', 4.50, 'Crispy garlic bread served with marinara sauce.'),
(1, 'Bruschetta', 5.25, 'Toasted bread topped with fresh tomatoes, garlic, and basil.'),
(1, 'Stuffed Mushrooms', 6.00, 'Mushrooms filled with cheese, herbs, and breadcrumbs.');
INSERT INTO items (category_id, name, price, description) VALUES
(2, 'Grilled Chicken Alfredo', 13.99, 'Grilled chicken breast served over fettuccine in Alfredo sauce.'),
(2, 'Beef Lasagna', 12.50, 'Layers of pasta, meat sauce, ricotta, and mozzarella cheese.'),
(2, 'Veggie Stir Fry', 10.99, 'Mixed seasonal vegetables sautéed in a savory soy-ginger sauce.');
INSERT INTO items (category_id, name, price, description) VALUES
(3, 'Cheesecake', 5.50, 'Classic New York-style cheesecake with a graham cracker crust.'),
(3, 'Chocolate Lava Cake', 6.25, 'Warm chocolate cake with a gooey molten center.'),
(3, 'Fruit Salad', 4.75, 'Freshly cut seasonal fruits.');
INSERT INTO items (category_id, name, price, description) VALUES
(4, 'Coke', 2.00, 'Chilled Coca-Cola served with ice.'),
(4, 'Orange Juice', 2.50, 'Fresh-squeezed orange juice.'),
(4, 'Coffee', 2.25, 'Freshly brewed hot coffee.');

INSERT INTO users (name, mail, isAdmin, isCheff, score, password_hash) VALUES
(
  'admin',
  'admin@example.com',
  TRUE,  
  FALSE,
  0,
  '$2a$10$BtwbYoJXZMoC5jMkq9MQc.5Ebej2kKRZ1S3A6KLBdk1B3ImKcjMlS'
),
(
  'chef',
  'chef@example.com',
  FALSE,
  TRUE,
  0,
  '$2a$10$VZkin5ExcvowQzErIG8A7O0zE.9zSmIvQxA.vaGZO2fBrbUp8E5yG'
);

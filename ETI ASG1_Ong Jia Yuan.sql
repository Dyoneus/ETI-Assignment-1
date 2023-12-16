CREATE USER 'user'@'localhost' IDENTIFIED BY
'password';
GRANT ALL ON *.* TO 'user'@'localhost';
CREATE DATABASE carpool;
CREATE DATABASE carpool_trips;

-- Note: Check if the 2 databases isn't stuck on fetching after creating. Restart MySQL workbench and problem is solved.
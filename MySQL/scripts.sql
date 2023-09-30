/*
---------------------------------------------------
-- Produced By: 4Tekies LLC and Penn State Abington - CMPSC 488 Course
-- Author: Mahir Khan, Joshua Ferrell & Joseph Oakes and Team 2 Members
-- Date: 6/27/2023, 09/28/2023
-- Purpose: OurGo holds all necessary MySQL code needed to establish the database
---------------------------------------------------
*/
-- Drop the database if it exists
DROP DATABASE IF EXISTS goengine;

-- DATABASE CHECK
-- This is used to see if the Database exists on the local computer
-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS goengine;

-- Switch to the goengine databaselog
USE goengine;

-- TABLE CHECK
-- This section is used on the user's computer to make sure they have all the proper tables.
-- If the user does not, then the tables are created for them

-- Create the lookup table for user roles
CREATE TABLE IF NOT EXISTS users_roles_lookup (
                                                  user_role NVARCHAR(5) PRIMARY KEY, -- Primary key representing user role
                                                  role_name NVARCHAR(25) -- Name of the role
);

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
                                     user_id CHAR(36) PRIMARY KEY, -- Unique identifier for the user
                                     user_name NVARCHAR(25), -- Name of the user
                                     user_login NVARCHAR(10), -- login credentials for user
                                     user_role NVARCHAR(5), -- User's role
                                     user_password VARBINARY(16), -- Encrypted password
                                     active_or_not BOOLEAN DEFAULT TRUE, -- Flag indicating if the user is active or not
                                     user_date_added DATETIME DEFAULT CURRENT_TIMESTAMP(), -- Date and time the user was added
                                     FOREIGN KEY (user_role) REFERENCES users_roles_lookup (user_role) -- Foreign key referencing user roles
);

-- Creates the logStatusCode lookup table for a reference to the log table
CREATE TABLE IF NOT EXISTS log_status_codes(
                                             status_code VARCHAR(3) PRIMARY KEY,
                                             status_message VARCHAR(250)
);


-- Creates the Log table
-- CREATE TABLE IF NOT EXISTS log (
--                                    logID CHAR (36)PRIMARY KEY, -- GUID for the ID, returns the ID as a case sensitive 36 long string
--                                    statusCode VARCHAR(3), -- 3 letters to store the status code to the DB
--                                    FOREIGN KEY (statusCode) REFERENCES logStatusCodes (statusCode), -- references to another table
--                                    message VARCHAR(250), -- A message about the status of the log
--                                    goEngineArea VARCHAR(250), -- Where the log status is occurring
--                                    dateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When inserting a value, the dateTime automatically updates to the time it occurred
-- );
-- Fix for Duplicate Key Issue:
-- Adds AUTO_INCREMENT to generate unique logID
CREATE TABLE IF NOT EXISTS log (
                                   log_ID INT AUTO_INCREMENT PRIMARY KEY, -- Auto-generated unique ID
                                   status_code VARCHAR(3), -- 3 letters to store the status code to the DB
                                   FOREIGN KEY (status_code) REFERENCES log_status_codes (status_code), -- references to another table
                                   message VARCHAR(250), -- A message about the status of the log
                                   go_engine_area VARCHAR(250), -- Where the log status is occurring
                                   date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When inserting a value, the dateTime automatically updates to the time it occurred
);
-- Creates the webservice table
CREATE TABLE IF NOT EXISTS web_service(
                                         web_service_ID CHAR(36)PRIMARY KEY, -- GUID for creating a unique ID
                                         web_service_description VARCHAR(250), -- A description of the service being offered
                                         customer_ID CHAR(36), -- We are using CHAR(36) for our GUID's, but other options exist
                                         access_token LONGTEXT, -- This lets the customer access the website. LONGTEXT is used to store JWT's of varying lengths
                                         date_active DATE, -- When the token is activated
                                         is_active BOOLEAN -- If the webservice is currently active or not
);

-- Creating url table for CRAB
CREATE TABLE IF NOT EXISTS urls (
                                    id CHAR(36) PRIMARY KEY, -- Unique identifier for the URLs using GUID
                                    url LONGTEXT NOT NULL, -- The URL string for storing the urls
                                    tags JSON, -- Optional JSON field for storing tags related to the URL
                                    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP() -- The timestamp of when the URL was created/added
);

-- Table for TaskManager
CREATE TABLE IF NOT EXISTS tasks (
    task_id CHAR(36) PRIMARY KEY, 
    task_name NVARCHAR(50),
    priority INT,
    status NVARCHAR(20),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for MachineLearningModels
CREATE TABLE IF NOT EXISTS machine_learning_models (
    model_id CHAR(36) PRIMARY KEY,
    model_name NVARCHAR(50),
    weights LONGTEXT,
    biases LONGTEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for WebCrawlers
USE goengine;
CREATE TABLE IF NOT EXISTS webcrawlers (
    crawler_id CHAR(36) PRIMARY KEY,
    source_url LONGTEXT,
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for the scraper engine
CREATE TABLE IF NOT EXISTS scraper_engine (
	engine_id CHAR(36) PRIMARY KEY,
    engine_name NVARCHAR(50),
    engine_description VARCHAR(250),
    time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Stored Procedure to add a new task
DELIMITER //
CREATE PROCEDURE create_task(
    IN p_task_name NVARCHAR(50),
    IN p_priority INT,
    IN p_status NVARCHAR(20)
)
BEGIN
    DECLARE v_task_id CHAR(36);
    SET v_task_id = UUID();
    INSERT INTO tasks (task_id, task_name, priority, status)
    VALUES (v_task_id, p_task_name, p_priority, p_status);
END //
DELIMITER ;

-- Stored Procedure to add a new machine learning model
DELIMITER //
CREATE PROCEDURE create_model(
    IN p_model_name NVARCHAR(50),
    IN p_weights LONGTEXT,
    IN p_biases LONGTEXT
)
BEGIN
    DECLARE v_model_id CHAR(36);
    SET v_model_id = UUID();
    INSERT INTO machine_learning_models (model_id, model_name, weights, biases)
    VALUES (v_model_id, p_model_name, p_weights, p_biases);
END //
DELIMITER ;

-- Stored Procedure to add a new web crawler
DELIMITER //
CREATE PROCEDURE create_webcrawler(
    IN p_source_url LONGTEXT
)
BEGIN
    DECLARE v_crawler_id CHAR(36);
    SET v_crawler_id = UUID();
    INSERT INTO webcrawlers (crawler_id, source_url)
    VALUES (v_crawler_id, p_source_url);
END //
DELIMITER ;

-- Stored Procedure to add a new scraper engine
DELIMITER //
CREATE PROCEDURE create_scraper_engine(
	IN p_engine_name NVARCHAR(50),
    IN p_engine_description VARCHAR(250)
)
BEGIN
	DECLARE v_engine_id CHAR(36);
    -- Generating a unique identifier and assigning it to v_engine_id
    SET v_engine_id = UUID();
    INSERT INTO scraper_engine(engine_id, engine_name, engine_description)
    VALUES (v_engine_id, p_engine_name, p_engine_description);
END //
DELIMITER ;

-- PROCEDURE CHECK

-- Set the delimiter for the following function creation
DELIMITER //

-- Create a function to encrypt passwords
CREATE FUNCTION encrypts_password(password NVARCHAR(10)) RETURNS VARBINARY(16)
    DETERMINISTIC
    NO SQL
BEGIN
    DECLARE encrypted VARBINARY(16);
    SET encrypted = AES_ENCRYPT(password, 'IST888IST888');
    RETURN encrypted;
END //
DELIMITER //

CREATE PROCEDURE get_status_code(IN statusCode VARCHAR(3)) -- Gets the Status code of the Log
BEGIN
    SELECT * FROM log AS l WHERE l.status_code = statusCode;
END //
DELIMITER // -- Needed to not throw an error because of complex coding

DELIMITER //

CREATE PROCEDURE insert_log(
    IN pStatusCode VARCHAR(3),
    IN pMessage VARCHAR(250),
    IN pGoEngineArea VARCHAR(250)
)
BEGIN
    DECLARE pLogID CHAR(36);
    SET pLogID = UUID();

    INSERT INTO log (log_ID, status_code, message, go_engine_area)
    VALUES (pLogID, pStatusCode, pMessage, pGoEngineArea);
END//

DELIMITER ;



DELIMITER //

CREATE PROCEDURE select_all_logs()
BEGIN
    SELECT log_ID, status_code, message, go_engine_area, date_time
    FROM log;
END //

DELIMITER //

DELIMITER //
CREATE PROCEDURE select_all_logs_by_status_code(IN pStatusCode VARCHAR(3))
BEGIN
    SELECT log_ID, status_code, message, go_engine_area, date_time
    FROM log
    WHERE status_code = pStatusCode;
END //

DELIMITER //
DELIMITER //
CREATE PROCEDURE populate_log_status_codes() -- Populates the log's status codes if they aren't already
BEGIN
    IF (SELECT COUNT(*) FROM log_status_codes) = 0 THEN
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('OPR', 'Normal operational mode');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('WAR', 'Warning: issue, application still functional');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('ERR', 'Severe error, application not functional');
    END IF;
END //

DELIMITER ;
DELIMITER //


select * from log
-- Reset the delimiter back to default
DELIMITER ;

DELIMITER //
-- CREATE
CREATE PROCEDURE create_user(
    IN p_name NVARCHAR(25),
    IN p_login NVARCHAR(10),
    IN p_role NVARCHAR(5),
    IN p_password VARBINARY(16),
    IN p_active BOOLEAN
)
BEGIN
    DECLARE v_id CHAR(36);

    SET v_id = UUID();

    INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
    VALUES (v_id, p_name, p_login, p_role, p_password, p_active, CURRENT_TIMESTAMP());
    END //

DELIMITER ;
-- READ
DELIMITER //

CREATE PROCEDURE get_users()
BEGIN
    SELECT user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added
    FROM users;
END //

DELIMITER ;

-- UPDATE
DELIMITER //

CREATE PROCEDURE update_user(
    IN p_id CHAR(36),
    IN p_name NVARCHAR(25),
    IN p_login NVARCHAR(10),
    IN p_role NVARCHAR(5),
    IN p_password VARBINARY(16)
)
BEGIN
    UPDATE users
    SET user_id = p_name,
        user_login = p_login,
        user_role = p_role,
        user_password = p_password
    WHERE user_id = p_id;
END //

DELIMITER ;

-- DELETE
DELIMITER //

CREATE PROCEDURE delete_user(
    IN p_user_id CHAR(36)
)
BEGIN
    DELETE FROM users
    WHERE user_id = p_user_id;
END //

DELIMITER ;

DELIMITER //
CREATE PROCEDURE get_random_url()
BEGIN
    SELECT * FROM urls ORDER BY RAND() LIMIT 1;
END //
DELIMITER ;

-- Procedure to retrieve only the 'url' column from the 'urls' table
DELIMITER //

CREATE PROCEDURE get_urls_only()
BEGIN
    -- Select only the 'url' column from the 'urls' table
    SELECT url FROM urls;
END //
DELIMITER ;

DELIMITER //
CREATE PROCEDURE get_urls_and_tags()
BEGIN
    SELECT url, tags FROM urls;
END //
DELIMITER ;

-- Insert values into users_roles_lookup table
INSERT INTO users_roles_lookup (user_role, role_name)
VALUES
    ('ADM', 'Administrator'),
    ('FAC', 'Faculty'),
    ('STD', 'Student'),
    ('DEV', 'Developer');

-- Insert values into users table
INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
VALUES
    (UUID(), 'Joseph Oakes', 'jxo19', 'ADM', encrypts_password('admin123'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Mahir Khan', 'mrk5928', 'DEV', encrypts_password('dev789'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Joshua Ferrell', 'jmf6913', 'DEV', encrypts_password('std447'), TRUE, CURRENT_TIMESTAMP());


-- Inserting the first record
INSERT INTO urls (id, url, tags)
VALUES (UUID(), 'https://sites.google.com/view/mahirbootstrap/home', '{"tag1": "<a>"}');

-- Inserting the second record
INSERT INTO urls (id, url, tags)
VALUES (UUID(), 'https://www.abington.psu.edu/WPL/mahir-khan', '{"tag2": "<img>"}');

-- Inserting the third record
INSERT INTO urls (id, url, tags)
VALUES (UUID(), 'https://sites.google.com/view/mahirbootstrap/signup?authuser=0', '{"tag3": "<label>"}');

-- Inserting the fourth record
INSERT INTO urls (id, url, tags)
VALUES (UUID(), 'https://sites.google.com/view/golangserver/home', '{"tag4": "<section>"}');

CALL populate_log_status_codes();
INSERT INTO users_roles_lookup (user_role, role_name)
VALUES ('1', 'User');

-- call PopulateLog();
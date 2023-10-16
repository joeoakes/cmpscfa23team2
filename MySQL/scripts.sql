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
                                                  user_role NVARCHAR(5) PRIMARY KEY , -- Primary key representing user role
                                                  role_name NVARCHAR(25) -- Name of the role
);

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
                                     user_id CHAR(36) PRIMARY KEY, -- Unique identifier for the user
                                     user_name NVARCHAR(25), -- Name of the user
                                     user_login NVARCHAR(10), -- login credentials for user
                                     user_role NVARCHAR(5), -- User's role
                                     user_password VARBINARY(255), -- Encrypted password
                                     active_or_not BOOLEAN DEFAULT TRUE, -- Flag indicating if the user is active or not
                                     user_date_added DATETIME DEFAULT CURRENT_TIMESTAMP(), -- Date and time the user was added
                                     FOREIGN KEY (user_role) REFERENCES users_roles_lookup (user_role) -- Foreign key referencing user roles
);

-- Creates the logStatusCode lookup table for a reference to the log table
CREATE TABLE IF NOT EXISTS log_status_codes(
                                             status_code VARCHAR(3) PRIMARY KEY,
                                             status_message VARCHAR(255)
);

-- Fix for Duplicate Key Issue:
-- Adds AUTO_INCREMENT to generate unique logID
CREATE TABLE IF NOT EXISTS log (
                                   log_ID INT AUTO_INCREMENT PRIMARY KEY, -- Auto-generated unique ID
                                   status_code VARCHAR(3), -- 3 letters to store the status code to the DB
                                   FOREIGN KEY (status_code) REFERENCES log_status_codes (status_code), -- references to another table
                                   message VARCHAR(255), -- A message about the status of the log
                                   go_engine_area VARCHAR(255), -- Where the log status is occurring
                                   date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When inserting a value, the dateTime automatically updates to the time it occurred
);

-- Periodically clean the log (anything older than 30 days) 
-- Temporarily disable safe update mode
SET SQL_SAFE_UPDATES = 0;
DELETE FROM log
WHERE date_time < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- Retrieve logs from the start of the day
SELECT* FROM log
WHERE date_time >= CURDATE();

-- get logs from the start of the week
SELECT* FROM log
WHERE date_time >= SUBDATE(CURDATE(), DAYOFWEEK(CURDATE()) - 1);

-- Creates the webservice table
CREATE TABLE IF NOT EXISTS web_service(
                                         web_service_ID CHAR(36)PRIMARY KEY, -- GUID for creating a unique ID
                                         web_service_description VARCHAR(255), -- A description of the service being offered
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
                                    domain LONGTEXT, --
                                    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP() -- The timestamp of when the URL was created/added
);

-- Stored Procedure to add a new prediction
-- DELIMITER //
-- CREATE PROCEDURE create_prediction(
-- 	IN p_engine_id CHAR(36),
-- 	IN p_prediction_info JSON
-- )
-- BEGIN
-- 	DECLARE v_prediction_id CHAR(36);
--     -- generate a unique identifier for the prediction and assign it to v_prediction_id
--     SET v_prediction_id = UUID();
--     INSERT INTO prerdictions (prediction_id, engine_id, prediction_info)
--     VALUES (v_prediction_id, p_engine_id, p_prediction_info);
-- END //
-- DELIMITER;


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

-- Table for predictions
CREATE TABLE IF NOT EXISTS predictions (
    prediction_id INT PRIMARY KEY AUTO_INCREMENT,
    engine_id CHAR(36),
    prediction_tag CHAR(64),  -- New field for clustering similar predictions
    input_data TEXT,
    prediction_info JSON,
    prediction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (engine_id) REFERENCES scraper_engine (engine_id)
);


CREATE TABLE user_sessions (
                               session_id INT PRIMARY KEY AUTO_INCREMENT,
                               user_id CHAR(36),
                               token VARCHAR(255) NOT NULL,
                               time_to_live DATETIME NOT NULL,
                               last_activity DATETIME NOT NULL,
                               scope VARCHAR(255) NOT NULL,
                               FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Create the user_permissions table
CREATE TABLE IF NOT EXISTS user_permissions (
                                                permission_id INT AUTO_INCREMENT PRIMARY KEY, -- Auto-generated unique ID for the permission
                                                user_role NVARCHAR(5), -- User's role
                                                action_name NVARCHAR(50), -- Name of the action or permission
                                                resource_name NVARCHAR(50) -- Name of the resource the permission applies to
);

-- Create the user_token_blacklist table
CREATE TABLE IF NOT EXISTS user_token_blacklist (
                                                    token_id INT AUTO_INCREMENT PRIMARY KEY, -- Auto-generated unique ID for the token
                                                    token VARCHAR(255) NOT NULL, -- The token to be invalidated
                                                    expiry_date DATETIME NOT NULL -- The date and time when the token expires or is invalidated
);

-- A SPROC for user login
DELIMITER //
CREATE PROCEDURE user_login(
    IN p_user_login NVARCHAR(10),
    IN p_user_password VARBINARY(16)
)
BEGIN
    SELECT user_id, user_name, user_role
    FROM users
    WHERE user_login = p_user_login AND user_password = p_user_password AND active_or_not = TRUE;
END //
DELIMITER ;



-- ================================================
-- SECTION: TASK MANAGER SPROCS
-- ================================================
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

-- Stored Procedure to update a task
DELIMITER //
CREATE PROCEDURE update_task(
    IN p_task_id CHAR(36),
    IN p_priority INT,
    IN p_status NVARCHAR(20)
)
BEGIN
    UPDATE tasks
    SET priority = p_priority,
        status = p_status
    WHERE task_id = p_task_id;
END //
DELIMITER ;

-- ================================================
-- SECTION: CUDA SPROCS
-- ================================================

-- Stored Procedure to add a new prediction
DELIMITER //
CREATE PROCEDURE create_prediction(
    IN p_engine_id CHAR(36),
    IN p_prediction_tag CHAR(64),  -- New parameter
    IN p_prediction_info JSON
)
BEGIN
    INSERT INTO predictions (engine_id, prediction_tag, prediction_info)
    VALUES (p_engine_id, p_prediction_tag, p_prediction_info);
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

-- New Stored Procedures
-- Stored Procedure to update a machine learning model
DELIMITER //
CREATE PROCEDURE update_model(
    IN p_model_id CHAR(36),
    IN p_weights LONGTEXT,
    IN p_biases LONGTEXT
)
BEGIN
    UPDATE machine_learning_models
    SET weights = p_weights,
        biases = p_biases
    WHERE model_id = p_model_id;
END //
DELIMITER ;
-- Setting the delimiter for the entire script
DELIMITER //

-- Stored Procedure to delete a machine learning model
CREATE PROCEDURE delete_model(
    IN p_model_id CHAR(36)
)
BEGIN
    DELETE FROM machine_learning_models
    WHERE model_id = p_model_id;
END //


-- ================================================
-- SECTION: LOG SPROCS
-- ================================================
-- Procedure to get status code
CREATE PROCEDURE get_status_code(IN statusCode VARCHAR(3))
BEGIN
    SELECT * FROM log AS l WHERE l.status_code = statusCode;
END //

-- Resetting the delimiter back to ;
DELIMITER ;


DELIMITER ;

DELIMITER //
CREATE PROCEDURE insert_or_update_status_code(
    IN p_status_code VARCHAR(3),
    IN p_status_message VARCHAR(255)
)
BEGIN
    DECLARE existing_count INT;

    -- Check if the status code already exists
    SELECT COUNT(*) INTO existing_count FROM log_status_codes WHERE status_code = p_status_code;

    IF existing_count = 0 THEN
        -- Insert a new status code
        INSERT INTO log_status_codes (status_code, status_message)
        VALUES (p_status_code, p_status_message);
    ELSE
        -- Update the existing status code
        UPDATE log_status_codes
        SET status_message = p_status_message
        WHERE status_code = p_status_code;
    END IF;
END //
DELIMITER ;

-- USE goengine;

-- DELIMITER //

-- DROP PROCEDURE IF EXISTS InsertLog;
-- CREATE PROCEDURE InsertLog(
--     IN pStatusCode VARCHAR(3),
--     IN pMessage VARCHAR(250),
--     IN pGoEngineArea VARCHAR(250)
-- )
-- BEGIN
--     DECLARE pLogID CHAR(36);
--     IF LENGTH(pStatusCode) > 3 THEN
--         SIGNAL SQLSTATE '45000'
--         SET MESSAGE_TEXT = 'Data too long for column pStatusCode';
--         RETURN;
--     END IF;
--     SET pLogID = UUID();
--     INSERT INTO log (logID, statusCode, message, goEngineArea)
--     VALUES (pLogID, pStatusCode, pMessage, pGoEngineArea);
-- END//
-- DELIMITER ;

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
CREATE PROCEDURE populate_log_status_codes() -- Populates the Log's if they aren't already

BEGIN
    IF (SELECT COUNT(*) FROM log_status_codes) = 0 THEN
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('OPR', 'Normal operational mode');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('WAR', 'Warring issue application still functional');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('ERR', 'Severe error application not functional');
    END IF;
END //

DELIMITER ;
DELIMITER //

-- CREATE PROCEDURE PopulateLog()
-- BEGIN
--     DECLARE statusCodeExists INT;

--     SELECT COUNT(*) INTO statusCodeExists FROM logstatuscodes WHERE statusCode IN ('ERR', 'WAR', 'OPR');

--     IF statusCodeExists = 3 THEN
--         IF (SELECT COUNT(*) FROM log) = 0 THEN
--             INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime)
--             VALUES (UUID(), 'ERR', 'An Error has occurred in the following area', 'CARP', NOW());

--             INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime)
--             VALUES (UUID(), 'WAR', 'A Warning has been issued in the following area', 'CRAB', NOW());

--             INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime)
--             VALUES (UUID(), 'OPR', 'Normal Operational Requirements have been met in the following area', 'CUDA', NOW());
--         END IF;
--     ELSE
--         -- If required code is missing from statusCodes, then this error is given
--         SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Missing required status codes in logstatuscodes table.';
--     END IF;
-- END //

-- DELIMITER ;

-- CALL PopulateLog();

select * from log
-- Reset the delimiter back to default
DELIMITER ;

-- ================================================
-- SECTION: CARP SPROCS
-- ================================================

DELIMITER //
-- CREATE
CREATE PROCEDURE create_user(
    IN p_user_name NVARCHAR(25),
    IN p_user_login NVARCHAR(10),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(16),
    IN p_active_or_not BOOLEAN
)
BEGIN
    DECLARE v_user_id CHAR(36);

    SET v_user_id = UUID();

    INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
    VALUES (v_user_id, p_user_name, p_user_login, p_user_role, p_user_password, p_active_or_not, CURRENT_TIMESTAMP());
    SELECT v_user_id;
    END //

DELIMITER ;
-- READ
-- A SPROC to get a specific user
DELIMITER //
CREATE PROCEDURE get_user_by_ID(IN p_user_id CHAR(36))
BEGIN
    SELECT * FROM users WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC to get all users in a specific role
DELIMITER //
CREATE PROCEDURE get_users_by_role(IN p_role NVARCHAR(5))
BEGIN
    SELECT * FROM users WHERE user_role = p_role;
END //
DELIMITER ;

-- A SPROC to get all users
DELIMITER //
CREATE PROCEDURE get_users()
BEGIN
    SELECT user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added
    FROM users;
END //
DELIMITER ;

-- A SPROC to fetch user using username
DELIMITER //
CREATE PROCEDURE fetch_user_id(
    IN p_user_name NVARCHAR(25)
)
BEGIN
    SELECT user_id FROM users WHERE user_name = p_user_name;
END //
DELIMITER ;

DELIMITER //

-- Create a function to encrypt passwords
CREATE FUNCTION encrypts_password(password NVARCHAR(255)) RETURNS VARBINARY(255)
    DETERMINISTIC
    NO SQL
BEGIN
    DECLARE encrypted VARBINARY(255);
    SET encrypted = AES_ENCRYPT(password, 'IST888IST888');
    RETURN encrypted;
END;

DELIMITER // 
-- A SPROC to validate user credentials
CREATE PROCEDURE `validate_user`(IN userLogin VARCHAR(255), IN userPassword NVARCHAR(255))
BEGIN
    DECLARE encryptedPwd VARBINARY(255);
    SET encryptedPwd = encrypts_password(userPassword);
    SELECT COUNT(*) FROM users WHERE user_login = userLogin AND user_password = encryptedPwd INTO @exists;
    SELECT IF(@exists > 0, TRUE, FALSE) AS valid;
END;


-- UPDATE
DELIMITER //

CREATE PROCEDURE update_user(
    IN p_user_id CHAR(36),
    IN p_user_name NVARCHAR(25),
    IN p_user_login NVARCHAR(10),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(16)
)
BEGIN
    UPDATE users
    SET user_name = p_user_name,
        user_login = p_user_login,
        user_role = p_user_role,
        user_password = p_user_password
    WHERE user_id = p_user_id;
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


-- ================================================
-- SECTION: Authentication and Authorization SPROCS
-- ================================================

-- SPROC for authenticating a user
DELIMITER //
CREATE PROCEDURE authenticate_user(
    IN p_user_login NVARCHAR(10),
    IN p_user_password VARBINARY(255)
)
BEGIN
    DECLARE v_user_id CHAR(36);
    DECLARE v_authenticated BOOLEAN;

    -- Check if the login and hashed password match any user
    SELECT user_id INTO v_user_id FROM users
    WHERE user_login = p_user_login AND user_password = p_user_password;

    -- Determine if the user is authenticated
    SET v_authenticated = (v_user_id IS NOT NULL);

    SELECT v_authenticated, v_user_id;
END //
DELIMITER ;

-- SPROC for getting the role of a user
DELIMITER //
CREATE PROCEDURE get_user_role(
    IN p_user_id CHAR(36)
)
BEGIN
    DECLARE v_user_role NVARCHAR(5);

    -- Fetch the role of the user
    SELECT user_role INTO v_user_role FROM users WHERE user_id = p_user_id;

    SELECT v_user_role;
END //
DELIMITER ;

-- SPROC for checking if a user is active
DELIMITER //
CREATE PROCEDURE is_user_active(
    IN p_user_id CHAR(36)
)
BEGIN
    DECLARE v_active BOOLEAN;

    -- Fetch the active status of the user
    SELECT active_or_not INTO v_active FROM users WHERE user_id = p_user_id;

    SELECT v_active;
END //
DELIMITER ;

-- SPROC for authorizing a user based on role
DELIMITER //
CREATE PROCEDURE authorize_user(
    IN p_user_id CHAR(36),
    IN required_role NVARCHAR(5)
)
BEGIN
    DECLARE v_user_role NVARCHAR(5);

    -- Fetch the role of the user
    SELECT user_role INTO v_user_role FROM users WHERE user_id = p_user_id;

    -- Check if the user is authorized to perform the operation
    IF v_user_role = required_role THEN
        SELECT TRUE AS is_authorized;
    ELSE
        SELECT FALSE AS is_authorized;
    END IF;
END //
DELIMITER ;

-- Reset the delimiter back to default
DELIMITER ;

-- UPDATE
-- A SPROC to update a user's role
DELIMITER //
CREATE PROCEDURE update_user_role(
    IN p_user_id CHAR(36),
    IN p_new_role NVARCHAR(5)
)
BEGIN
    UPDATE users
    SET user_role = p_new_role
    WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC to update a user's password
DELIMITER //
CREATE PROCEDURE update_user_password(
    IN p_user_id CHAR(36),
    IN p_new_password VARBINARY(16)
)
BEGIN
    UPDATE users
    SET user_password = p_new_password
    WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC to deactivate a user
DELIMITER //
CREATE PROCEDURE deactivate_user(
    IN p_user_id CHAR(36)
)
BEGIN
    UPDATE users
    SET active_or_not = FALSE
    WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC for user registration
DELIMITER //
CREATE PROCEDURE user_registration(
    IN p_user_name NVARCHAR(25),
    IN p_user_login NVARCHAR(10),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(16),
    IN p_active_or_not BOOLEAN
)
BEGIN
    CALL create_user(p_user_name, p_user_login, p_user_role, p_user_password, p_active_or_not);
END //
DELIMITER ;




-- ================================================
-- SECTION: CRAB SPROCS
-- ================================================

-- Stored Procedure to add a new web crawler
DELIMITER //
DELIMITER //
CREATE PROCEDURE create_webcrawler(
    IN p_source_url LONGTEXT
)
BEGIN
    DECLARE v_crawler_id CHAR(36);
    SET v_crawler_id = UUID();
    INSERT INTO webcrawlers (crawler_id, source_url)
    VALUES (v_crawler_id, p_source_url);
    SELECT v_crawler_id;
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
    SELECT v_engine_id;
END //
DELIMITER ;

-- SPROC to insert URL records into the URLs table
DELIMITER //
CREATE PROCEDURE insert_url(IN p_url LONGTEXT, IN p_tags JSON, IN p_domain LONGTEXT)
BEGIN
    DECLARE v_id CHAR(36);
    SET v_id = UUID();
    INSERT INTO urls (id, url, tags, domain, created_time)
    VALUES (v_id, p_url, p_tags, p_domain, CURRENT_TIMESTAMP);
    SELECT v_id; -- This line returns the generated ID.
END //
DELIMITER ;


-- SPROC to update URL
DELIMITER //
CREATE PROCEDURE update_url(IN p_id CHAR(36), IN p_url LONGTEXT, IN p_tags JSON, IN p_domain LONGTEXT)
BEGIN
    UPDATE urls SET url = p_url, tags = p_tags, domain = p_domain WHERE id = p_id;
END //
DELIMITER ;

-- SPROC for domain-specific queries: get URL tags and domain
DELIMITER //
CREATE PROCEDURE get_url_tags_and_domain(IN p_id CHAR(36))
BEGIN
    SELECT tags, domain FROM urls WHERE id = p_id;
END //
DELIMITER ;

-- SPROC to get URLs from a specific domain
DELIMITER //
CREATE PROCEDURE get_urls_from_domain(IN p_domain LONGTEXT)
BEGIN
    SELECT * FROM urls WHERE domain = p_domain;
END //
DELIMITER ;

-- SPROC to get UUID from URL and domain
DELIMITER //
CREATE PROCEDURE get_Uuid_from_URL_and_domain(IN p_url LONGTEXT, IN p_domain LONGTEXT)
BEGIN
    SELECT id FROM urls WHERE url = p_url AND domain = p_domain;
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

DELIMITER //
CREATE PROCEDURE get_urls_and_tags()
BEGIN
    SELECT url, tags FROM urls;
END //
DELIMITER ;



 -- ================================================
-- SECTION: VALIDATE TOKEN:
-- ================================================

 #Tables are not created yet
 #Tables left to create user_sessions, user_permissions

-- Setting the delimiter for stored procedures
DELIMITER //

-- Procedure to create a new session for a user
CREATE PROCEDURE create_session(
    IN p_user_id CHAR(36),
    IN p_token TEXT
)
BEGIN
    -- Inserting a new session with details and setting an expiration time of 1 hour
    INSERT INTO user_sessions (session_id, user_id, token, time_to_live, last_activity, scope)
    VALUES (UUID(), p_user_id, p_token, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 1 HOUR), CURRENT_TIMESTAMP, 'default');
END //
DELIMITER ;

-- Procedure to validate a user's token
DELIMITER //
CREATE PROCEDURE validate_token(
    IN p_token TEXT
)
BEGIN
    -- Checking if the token is valid and still within its active time frame
    SELECT user_id, time_to_live > CURRENT_TIMESTAMP AS is_valid
    FROM user_sessions
    WHERE token = p_token;
END //
DELIMITER ;

-- Procedure to add a new permission for a user role
DELIMITER //
CREATE PROCEDURE add_permission(
    IN p_user_role NVARCHAR(5),
    IN p_action_name NVARCHAR(100),
    IN p_resource_name NVARCHAR(100)
)
BEGIN
    -- Inserting a new permission record for the given user role, action, and resource
    INSERT INTO user_permissions (permission_id, user_role, action_name, resource_name)
    VALUES (UUID(), p_user_role, p_action_name, p_resource_name);
END //
DELIMITER ;

-- Procedure to check if a user role has a specific permission
DELIMITER //
CREATE PROCEDURE check_permission(
    IN p_user_role NVARCHAR(5),
    IN p_action_name NVARCHAR(100),
    IN p_resource_name NVARCHAR(100)
)
BEGIN
    -- Checking if a permission exists for the given user role, action, and resource
    SELECT COUNT(*) > 0 AS has_permission
    FROM user_permissions
    WHERE user_role = p_user_role AND action_name = p_action_name AND resource_name = p_resource_name;
END //
DELIMITER ;

-- Inserting predefined roles into the user roles lookup table
INSERT INTO users_roles_lookup (user_role, role_name)
VALUES
    ('ADM', 'Administrator'),
    ('1', 'User'),
    ('FAC', 'Faculty'),
    ('STD', 'Student'),
    ('DEV', 'Developer');

-- Inserting sample users into the users table
INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
VALUES
    (UUID(), 'Joesph Oakes', 'jxo19', 'ADM', encrypts_password('admin123'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Mahir Khan', 'mrk5928', 'DEV', encrypts_password('dev789'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Joshua Ferrell', 'jmf6913', 'DEV', encrypts_password('std447'), TRUE, CURRENT_TIMESTAMP());

-- Inserting sample URLs into the URLs table
INSERT INTO urls (id, url, tags)
VALUES
    (UUID(), 'https://sites.google.com/view/mahirbootstrap/home', '{"tag1": "<a>"}'),
    (UUID(), 'https://www.abington.psu.edu/WPL/mahir-khan', '{"tag2": "<img>"}'),
    (UUID(), 'https://sites.google.com/view/mahirbootstrap/signup?authuser=0', '{"tag3": "<label>"}'),
    (UUID(), 'https://sites.google.com/view/golangserver/home', '{"tag4": "<section>"}');

-- Call to the procedure to populate log status codes
CALL populate_log_status_codes();

use goengine;

CALL user_registration('test_user', 'test_login', 'STD', 'test_password', true);


call populate_log_status_codes();


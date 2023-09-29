/*
---------------------------------------------------
-- Produced By: 4Tekies LLC
-- Author: Mahir Khan, Joshua Ferrell & Joseph Oakes
-- Date: 6/27/2023
-- Purpose: OurGo holds all necessary MySQL code needed to establish the database
---------------------------------------------------
*/
-- Drop the database if it exists
DROP DATABASE IF EXISTS goengine;

-- DATABASE CHECK
-- This is used to see if the Database exists on the local computer
-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS goengine;

-- Switch to the goengine database
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
use goengine;
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
CREATE TABLE IF NOT EXISTS logStatusCodes(
                                             statusCode VARCHAR(3) PRIMARY KEY,
                                             statusMessage VARCHAR(250)
);


-- Creates the Log table
CREATE TABLE IF NOT EXISTS log (
                                   logID CHAR (36)PRIMARY KEY, -- GUID for the ID, returns the ID as a case sensitive 36 long string
                                   statusCode VARCHAR(3), -- 3 letters to store the status code to the DB
                                   FOREIGN KEY (statusCode) REFERENCES logStatusCodes (statusCode), -- references to another table
                                   message VARCHAR(250), -- A message about the status of the log
                                   goEngineArea VARCHAR(250), -- Where the log status is occurring
                                   dateTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When inserting a value, the dateTime automatically updates to the time it occurred
);

-- Creates the webservice table
CREATE TABLE IF NOT EXISTS webservice(
                                         webServiceID CHAR(36)PRIMARY KEY, -- GUID for creating a unique ID
                                         webServiceDescription VARCHAR(250), -- A description of the service being offered
                                         customerID CHAR(36), -- We are using CHAR(36) for our GUID's, but other options exist
                                         accessToken LONGTEXT, -- This lets the customer access the website. LONGTEXT is used to store JWT's of varying lengths
                                         dateActive DATE, -- When the token is activated
                                         isActive BOOLEAN -- If the webservice is currently active or not
);

-- Creating url table for CRAB
CREATE TABLE IF NOT EXISTS urls (
                                    id CHAR(36) PRIMARY KEY, -- Unique identifier for the URLs using GUID
                                    url LONGTEXT NOT NULL, -- The URL string for storing the urls
                                    tags JSON, -- Optional JSON field for storing tags related to the URL
                                    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP() -- The timestamp of when the URL was created/added
);

-- Table for TaskManager
USE goengine;
CREATE TABLE IF NOT EXISTS tasks (
    task_id CHAR(36) PRIMARY KEY, 
    task_name NVARCHAR(50),
    priority INT,
    status NVARCHAR(20),
    created_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for MachineLearningModels
USE goengine;
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

-- Stored Procedure to add a new task
USE goengine;
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
USE goengine;
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
USE goengine;
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

-- PROCEDURE CHECK

-- Switch to the goengine database
USE goengine;

-- Set the delimiter for the following function creation
DELIMITER //

-- Create a function to encrypt passwords
CREATE FUNCTION EncryptsPassword(password NVARCHAR(10)) RETURNS VARBINARY(16)
    DETERMINISTIC
    NO SQL
BEGIN
    DECLARE encrypted VARBINARY(16);
    SET encrypted = AES_ENCRYPT(password, 'IST888IST888');
    RETURN encrypted;
END //
DELIMITER //
CREATE PROCEDURE GetStatusCode(IN statusCode VARCHAR(3)) -- Gets the Status code of the Log
BEGIN
    SELECT * FROM log AS l WHERE l.statusCode = statusCode;
END //
DELIMITER // -- Needed to not throw an error because of complex coding


USE goengine;

DELIMITER //

CREATE PROCEDURE InsertLog(
    IN pStatusCode VARCHAR(3),
    IN pMessage VARCHAR(250),
    IN pGoEngineArea VARCHAR(250)
)
BEGIN
    DECLARE pLogID CHAR(36);
    SET pLogID = UUID();

    INSERT INTO log (logID, statusCode, message, goEngineArea)
    VALUES (pLogID, pStatusCode, pMessage, pGoEngineArea);
END//

DELIMITER ;

USE goengine;

DELIMITER //

CREATE PROCEDURE SelectAllLogs()
BEGIN
    SELECT logID, statusCode, message, goEngineArea, dateTime
    FROM log;
END //

DELIMITER //

DELIMITER //
CREATE PROCEDURE SelectAllLogsByStatusCode(IN pStatusCode VARCHAR(3))
BEGIN
    SELECT logID, statusCode, message, goEngineArea, dateTime
    FROM log
    WHERE statusCode = pStatusCode;
END //

DELIMITER //
DELIMITER //
CREATE PROCEDURE PopulateLogStatusCodes() -- Populates the Log's if they aren't already

BEGIN
    IF (SELECT COUNT(*) FROM logStatusCodes) = 0 THEN
        INSERT INTO logStatusCodes (statusCode, statusMessage) VALUES ('OPR', 'Normal operational mode');
        INSERT INTO logStatusCodes (statusCode, statusMessage) VALUES ('WAR', 'Warring issue application still functional');
        INSERT INTO logStatusCodes (statusCode, statusMessage) VALUES ('ERR', 'Severe error application not functional');
    END IF;
END //

DELIMITER ;

-- Set the delimiter //beginning
DELIMITER //

-- Populate Log Procedure
CREATE PROCEDURE PopulateLog()
BEGIN
    DECLARE statusCodeExists INT;
SELECT COUNT(*) INTO statusCodeExists FROM logstatuscodes WHERE statusCode IN ('ERR', 'WAR', 'OPR');
IF statusCodeExists = 3 THEN
        IF (SELECT COUNT(*) FROM log) = 0 THEN
            INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime) VALUES (UUID(), 'ERR', 'An Error has occurred in the following area', 'CARP', NOW());
INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime) VALUES (UUID(), 'WAR', 'A Warning has been issued in the following area', 'CRAB', NOW());
INSERT INTO log (logID, statusCode, message, goEngineArea, dateTime) VALUES (UUID(), 'OPR', 'Normal Operational Requirements have been met in the following area', 'CUDA', NOW());
END IF;
ELSE
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Missing required status codes in logstatuscodes table.';
END IF;
END //
-- Reset the delimiter
DELIMITER ;

-- Create User Procedure //user role removed
DELIMITER //
CREATE PROCEDURE create_user(
    IN v_UUID NVARCHAR(25),
    IN v_account_id NVARCHAR(25), --some redundancy with UUID?
    IN p_user_login NVARCHAR(25) --also referred to as email
    IN p_user_name NVARCHAR(25),
    IN p_user_password_hashed VARBINARY(16),
    IN p_first_name NVARCHAR(25),
    IN p_last_name NVARCHAR(25),
    IN p_date_registered DATE
    IN p_last_login DATE
        )
BEGIN
    DECLARE v_user_id CHAR(36);
    SET v_user_id = UUID();
--//user role removed
INSERT INTO users (UUID, account_id, user_login, user_password_hashed, first_name, last_name, date_registered, last_login) VALUES (v_UUID, v_account_id, p_user_login, p_user_password_hashed, p_first_name, p_last_name, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());
END //
DELIMITER ;

-- Get Users Procedure //user role removed
DELIMITER //
CREATE PROCEDURE get_users()
BEGIN
SELECT account_id, user_login, user_password_hashed, last_login, date_registered FROM users;
END //
DELIMITER ;

-- Update User Procedure //user role removed
DELIMITER //
CREATE PROCEDURE update_user(
    IN p_account_id NVARCHAR(25),
    IN p_user_login NVARCHAR(25), --also referred to as email
    IN p_user_password_hashed VARBINARY(16),
    IN p_last_login DATE
        )
BEGIN
UPDATE users SET account_id = p_account_id, user_login = p_user_login, user_password_hashed = p_user_password_hashed WHERE user_login = p_user_login, last_login = p_last_login;
END //
DELIMITER ;
--end evan code

-- DELETE
USE goengine;
DELIMITER //

CREATE PROCEDURE delete_user(
    IN p_user_id CHAR(36)
)
BEGIN
    DELETE FROM users
    WHERE user_id = p_user_id;
END //

DELIMITER ;

USE goengine;

DELIMITER //
CREATE PROCEDURE GetRandomURL()
BEGIN
    SELECT * FROM urls ORDER BY RAND() LIMIT 1;
END //
DELIMITER ;

USE goengine;


-- Procedure to retrieve only the 'url' column from the 'urls' table
DELIMITER //

CREATE PROCEDURE GetURLsOnly()
BEGIN
    -- Select only the 'url' column from the 'urls' table
    SELECT url FROM urls;
END //

DELIMITER ;

use goengine;
DELIMITER //
CREATE PROCEDURE GetUrlsAndTags()
BEGIN
    SELECT url, tags FROM urls;
END //
DELIMITER ;

USE goengine;
-- Insert values into users_roles_lookup table
INSERT INTO users_roles_lookup (user_role, role_name)
VALUES
    ('ADM', 'Administrator'),
    ('FAC', 'Faculty'),
    ('STD', 'Student'),
    ('DEV', 'Developer');

use goengine;
-- Insert values into users table
INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
VALUES
    (UUID(), 'Joesph Oakes', 'jxo19', 'ADM', EncryptsPassword('admin123'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Mahir Khan', 'mrk5928', 'DEV', EncryptsPassword('dev789'), TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Joshua Ferrell', 'jmf6913', 'DEV', EncryptsPassword('std447'), TRUE, CURRENT_TIMESTAMP());


USE goengine;

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

call PopulateLogStatusCodes();
call PopulateLog();
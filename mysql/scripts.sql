/*
---------------------------------------------------
-- Produced By: CMPSC 488 Fall 2023 Team 2
-- Author: Matthew Assali, Sara Becker, Emily Carpenter, Matthew Finn, Hansi Seitaj, Evan Green, Binh Hoang, Eni Vejseli
-- Date: 12/05/2023
-- Purpose: holds all necessary mysql code needed to establish the database
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
                                     user_login NVARCHAR(36), -- login credentials for user
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
                                   log_ID BINARY(36) PRIMARY KEY, -- Use BINARY(16) to store UUIDs
                                   status_code VARCHAR(3),
                                   FOREIGN KEY (status_code) REFERENCES log_status_codes (status_code),
                                   message VARCHAR(255),
                                   go_engine_area VARCHAR(255),
                                   date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


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

CREATE TABLE scrapedData (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             domain VARCHAR(255),
                             title VARCHAR(255),
                             url VARCHAR(500),
                             description TEXT,
                             price VARCHAR(100),
                             source VARCHAR(255),
                             timestamp DATETIME
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


-- Table for predictions
-- Table for K-Nearest Neighbors Predictions
CREATE TABLE IF NOT EXISTS knn_predictions (
                                               prediction_id VARCHAR(36) PRIMARY KEY,
                                               query_identifier VARCHAR(255),
                                               input_data VARCHAR(255),
                                               prediction_info VARCHAR(255),
                                               prediction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for Linear Regression Predictions
CREATE TABLE IF NOT EXISTS linear_regression_predictions (
                                                             prediction_id VARCHAR(36) PRIMARY KEY,
                                                             query_identifier VARCHAR(255),
                                                             input_data TEXT(255),
                                                             prediction_info VARCHAR(255),
                                                             prediction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);

-- Table for Naive Bayes Predictions
CREATE TABLE IF NOT EXISTS naive_bayes_predictions (
                                                       prediction_id VARCHAR(36) PRIMARY KEY,
                                                       query_identifier VARCHAR(255),
                                                       input_data VARCHAR(255),
                                                       prediction_info LONGTEXT,
                                                       prediction_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);



CREATE TABLE IF NOT EXISTS user_sessions (
                                             session_id CHAR(36) PRIMARY KEY,
                                             user_id CHAR(36),
                                             token VARCHAR(255) NOT NULL,
                                             time_to_live DATETIME NOT NULL,
                                             last_activity DATETIME NOT NULL,
                                             scope VARCHAR(255) NOT NULL,
                                             FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Create the user_permissions table
CREATE TABLE IF NOT EXISTS user_permissions (
                                                permission_id CHAR(36) PRIMARY KEY, -- Auto-generated unique ID for the permission
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
-- Create the refresh_tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
                                              token_id CHAR(36) PRIMARY KEY,
                                              user_id CHAR(36),
                                              token VARBINARY(255),
                                              expiry DATETIME,
                                              FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE ETFs (
                      etf_id INT AUTO_INCREMENT PRIMARY KEY,
                      title VARCHAR(255) NOT NULL,
                      replication VARCHAR(255),
                      earnings VARCHAR(255),
                      total_expense_ratio VARCHAR(255),
                      tracking_difference VARCHAR(255),
                      fund_size VARCHAR(255),
                      isin VARCHAR(255) UNIQUE NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



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


DELIMITER //

CREATE PROCEDURE insert_log(
    IN pStatusCode VARCHAR(3),
    IN pMessage VARCHAR(250),
    IN pGoEngineArea VARCHAR(250)
)
BEGIN
    DECLARE pLogID BINARY(16);
    SET pLogID = UNHEX(REPLACE(UUID(), '-', '')); -- Generate a UUID and convert it to binary

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
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('200', 'Normal operational mode');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('WAR', 'Warring issue application still functional');
        INSERT INTO log_status_codes (status_code, status_message) VALUES ('400', 'Severe error application not functional');
    END IF;
END //

DELIMITER ;
DELIMITER //


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
    IN p_user_login NVARCHAR(36),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(255),
    IN p_active_or_not BOOLEAN
)
BEGIN
    DECLARE v_user_id CHAR(36);

    SET v_user_id = UUID();

    INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
    VALUES (v_user_id, p_user_name, p_user_login, p_user_role, p_user_password, p_active_or_not, CURRENT_TIMESTAMP());
    SELECT v_user_id;
END //

DELIMITER //
CREATE PROCEDURE get_user_by_login(
    IN p_user_login NVARCHAR(10)
)
BEGIN
    SELECT * FROM users WHERE user_login = p_user_login;
END //
DELIMITER ;

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


-- UPDATE
DELIMITER //

CREATE PROCEDURE update_user(
    IN p_user_id CHAR(36),
    IN p_user_name NVARCHAR(25),
    IN p_user_login NVARCHAR(36),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(255)
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
CREATE PROCEDURE update_url(IN p_id LONG, IN p_url LONG, IN p_tags JSON, IN p_domain LONG)
BEGIN
    UPDATE urls SET url = p_url, tags = p_tags, domain = p_domain WHERE id = p_id;
END //
DELIMITER ;

-- SPROC for domain-specific queries: get URL tags and domain
DELIMITER //
CREATE PROCEDURE get_url_tags_and_domain(IN p_id LONG)
BEGIN
    SELECT tags, domain FROM urls WHERE id = p_id;
END //
DELIMITER ;

-- SPROC to get URLs from a specific domain
DELIMITER //
CREATE PROCEDURE get_urls_from_domain(IN p_domain LONG)
BEGIN
    SELECT * FROM urls WHERE domain = p_domain;
END //
DELIMITER ;

-- SPROC to get UUID from URL and domain
DELIMITER //
CREATE PROCEDURE get_Uuid_from_URL_and_domain(IN p_url LONG, IN p_domain LONG)
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

-- SPROC for Inserting or Updating ETF data
# DELIMITER //
# CREATE PROCEDURE InsertOrUpdateETFData(
#     IN p_title VARCHAR(255),
#     IN p_replication VARCHAR(255),
#     IN p_earnings VARCHAR(255),
#     IN p_total_expense_ratio VARCHAR(255),
#     IN p_tracking_difference VARCHAR(255),
#     IN p_fund_size VARCHAR(255),
#     IN p_isin VARCHAR(255)
# )
# BEGIN
#     IF NOT EXISTS (SELECT * FROM ETFs WHERE isin = p_isin) THEN
#         INSERT INTO ETFs (title, replication, earnings, total_expense_ratio, tracking_difference, fund_size, isin)
#         VALUES (p_title, p_replication, p_earnings, p_total_expense_ratio, p_tracking_difference, p_fund_size, p_isin);
#     ELSE
#         UPDATE ETFs
#         SET
#             title = p_title,
#             replication = p_replication,
#             earnings = p_earnings,
#             total_expense_ratio = p_total_expense_ratio,
#             tracking_difference = p_tracking_difference,
#             fund_size = p_fund_size
#         WHERE isin = p_isin;
#     END IF;
# END //
DELIMITER ;


-- SPROC for Retrieving ETF data by ISIN
# DELIMITER //
# CREATE PROCEDURE FetchETFByISIN(
#     IN p_isin VARCHAR(255)
# )
# BEGIN
#     SELECT * FROM ETFs WHERE isin = p_isin;
# END //
# DELIMITER ;

-- SPROC for Deleting ETF data by ISIN
DELIMITER //
CREATE PROCEDURE DeleteETFByISIN(
    IN p_isin VARCHAR(255)
)
BEGIN
    DELETE FROM ETFs WHERE isin = p_isin;
END //
DELIMITER ;

-- SPROC for Listing All ETFs
DELIMITER //
CREATE PROCEDURE ListAllETFs()
BEGIN
    SELECT * FROM ETFs;
END //
DELIMITER ;

-- (Storing data sprocs):
DELIMITER $$
CREATE PROCEDURE InsertOrUpdateETFData(
    IN p_title VARCHAR(255),
    IN p_replication VARCHAR(255),
    IN p_earnings VARCHAR(255),
    IN p_totalExpenseRatio VARCHAR(255),
    IN p_trackingDifference VARCHAR(255),
    IN p_fundSize VARCHAR(255),
    IN p_isin VARCHAR(255)
)
BEGIN
    INSERT INTO ETFs (title, replication, earnings, total_expense_ratio, tracking_difference, fund_size, isin)
    VALUES (p_title, p_replication, p_earnings, p_totalExpenseRatio, p_trackingDifference, p_fundSize, p_isin)
    ON DUPLICATE KEY UPDATE
                         title = VALUES(title),
                         replication = VALUES(replication),
                         earnings = VALUES(earnings),
                         total_expense_ratio = VALUES(total_expense_ratio),
                         tracking_difference = VALUES(tracking_difference),
                         fund_size = VALUES(fund_size);
END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE FetchETFByISIN(IN p_isin VARCHAR(255))
BEGIN
    SELECT * FROM ETFs WHERE isin = p_isin;
END$$
DELIMITER ;

DELIMITER $$
CREATE PROCEDURE UpdateFundSizeByISIN(IN p_isin VARCHAR(255), IN p_fundSize VARCHAR(255))
BEGIN
    UPDATE ETFs SET fund_size = p_fundSize WHERE isin = p_isin;
END$$
DELIMITER ;

DELIMITER $$
# CREATE PROCEDURE DeleteETFByISIN(IN p_isin VARCHAR(255))
# BEGIN
#     DELETE FROM ETFs WHERE isin = p_isin;
# END$$
# DELIMITER ;

--

-- ================================================
-- SECTION: Authorization SPROCS
-- ================================================

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

DELIMITER //
CREATE PROCEDURE get_permissions_for_role(
    IN p_user_role NVARCHAR(5)
)
BEGIN
    -- Fetch all permissions associated with the given user role
    SELECT action_name, resource_name
    FROM user_permissions
    WHERE user_role = p_user_role;
END//
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



-- ================================================
-- SECTION: Authentication SPROCS:
-- ================================================

-- SPROC for authenticating a user
DELIMITER //
CREATE PROCEDURE authenticate_user(
    IN p_user_login NVARCHAR(36)
)
BEGIN
    DECLARE v_user_id CHAR(36);
    DECLARE v_hashed_password LONGTEXT; -- Changed to LONGTEXT

    SELECT user_id, user_password INTO v_user_id, v_hashed_password FROM users
    WHERE user_login = p_user_login;

    SELECT v_user_id, v_hashed_password;
END //
DELIMITER ;

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
END ;
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

DELIMITER  //
-- Procedure to validate a user's token
CREATE PROCEDURE validate_refresh_token(
    IN p_token VARCHAR(255)
)
BEGIN
    SELECT user_id, token, expiry
    FROM refresh_tokens
    WHERE token = p_token;
END //

-- Procedure to issue a new refresh token
CREATE PROCEDURE issue_refresh_token(
    IN p_user_id CHAR(36),
    IN p_token VARBINARY(255)
)
BEGIN
    DELETE FROM refresh_tokens WHERE user_id = p_user_id;
    INSERT INTO refresh_tokens (token_id, user_id, token, expiry)
    VALUES (UUID(), p_user_id, p_token, DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 7 DAY));
END //

DELIMITER ;

DELIMITER //
CREATE PROCEDURE logout_user(
    IN p_user_id CHAR(36)
)
BEGIN
    DELETE FROM user_sessions WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC to update a user's password
DELIMITER //
CREATE PROCEDURE change_user_password(
    IN p_user_id CHAR(36),
    IN p_new_password VARBINARY(255)
)
BEGIN
    UPDATE users
    SET user_password = p_new_password
    WHERE user_id = p_user_id;
END //
DELIMITER ;

-- A SPROC for user registration
DELIMITER //
CREATE PROCEDURE user_registration(
    IN p_user_name NVARCHAR(25),
    IN p_user_login NVARCHAR(36),
    IN p_user_role NVARCHAR(5),
    IN p_user_password VARBINARY(255),
    IN p_active_or_not BOOLEAN
)
BEGIN
    CALL create_user(p_user_name, p_user_login, p_user_role, p_user_password, p_active_or_not);
END //
DELIMITER ;

-- Setting the delimiter for stored procedures
DELIMITER //


-- Sproc for invalidate Token and refresh_token
DELIMITER //
CREATE PROCEDURE invalidate_token(
    IN p_user_id CHAR(36)
)
BEGIN
    DELETE FROM user_sessions WHERE user_id = p_user_id;
END //
DELIMITER ;


-- A SPROC for user login
DELIMITER //
CREATE PROCEDURE user_login(
    IN p_user_login NVARCHAR(36),
    IN p_user_password VARBINARY(255)
)
BEGIN
    SELECT user_id, user_name, user_role
    FROM users
    WHERE user_login = p_user_login AND user_password = p_user_password AND active_or_not = TRUE;
END //
DELIMITER ;


-- ================================================
-- SECTION: INSERTS & CALLS
-- ================================================
-- Inserting predefined roles into the user roles lookup table
INSERT INTO users_roles_lookup (user_role, role_name)
VALUES
    ('ADM', 'Administrator'),
    ('USR', 'User'),
    ('DEV', 'Developer');

# -- Inserting sample users into the users table
# INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
# VALUES
#     (UUID(), 'Joesph Oakes', 'jxo19', 'ADM', 'admin123', TRUE, CURRENT_TIMESTAMP()),
#     (UUID(), 'Mahir Khan', 'mrk5928', 'DEV', 'dev789', TRUE, CURRENT_TIMESTAMP()),
#     (UUID(), 'Joshua Ferrell', 'jmf6913', 'DEV', 'std447', TRUE, CURRENT_TIMESTAMP());

-- Inserting sample users into the users table
INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
VALUES
    ('9c0f0ac1-8d78-11ee-b6e0-4c796ed97681', 'test1', 'test1@test.com', 'USR', '$2a$10$6nsLKZMGjnG4osvBN3AbUOIvOnYXXZVrbcgdYY419OYUsGzqDlDMG', TRUE, '2023-11-27 17:59:24'),
    ('a2eb8427-8d78-11ee-b6e0-4c796ed97681', 'test2', 'test2@test.com', 'USR', '$2a$10$M8s0NhMKr24C6bSwlWBfY.4pPSnWtHIAAVY5qKRPfnoXZAFvzcmgW', TRUE, '2023-11-27 17:59:36'),
    (UUID(), 'hansi', 'hansi@hansi.com', 'USR', '$2a$10$C4ZoMvNpBqJ8MB9LMLzQye2uXvQKPujw1SXccnuLJ/frYoG6GUOZy', TRUE, CURRENT_TIMESTAMP());

-- Inserting admin users into the users table with a hashed password
INSERT INTO users (user_id, user_name, user_login, user_role, user_password, active_or_not, user_date_added)
VALUES
    (UUID(), 'Matthew Assali', 'mfa5498@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Hansi Seitaj', 'hjs5684@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Eni Vejseli', 'emv5319@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Sara Becker', 'sqb6198@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Emily Carpenter', 'esc5316@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Matthew Finn', 'mkf5480@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Evan M Green', 'emg5555@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Binh Thanh Hoang', 'bth5241@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP()),
    (UUID(), 'Shiv Patel', 'sbp5769@psu.edu', 'ADM', '$2a$10$hashedPasswordOfPassword', TRUE, CURRENT_TIMESTAMP());

INSERT INTO linear_regression_predictions(prediction_id,query_identifier, input_data, prediction_info) VALUES (UUID(),'Gas Prices Prediction for the year 2023','(1978.000000, 51.900000, 0.652000),(1979.000000, 70.200000, 0.882000),(1980.000000, 97.500000, 1.221000),(1981.000000, 108.500000, 1.353000),(1982.000000, 102.800000, 1.281000),(1983.000000, 99.400000, 1.225000),(1984.000000, 97.800000, 1.198000),(1985.000000, 98.600000, 1.196000),(1986.000000, 77.000000, 0.931000),(1987.000000, 80.100000, 0.957000),(1988.000000, 80.800000, 0.964000),(1989.000000, 88.500000, 1.060000),(1990.000000, 101.000000, 1.217000),(1991.000000, 99.200000, 1.196000),(1992.000000, 99.000000, 1.190000),(1993.000000, 97.700000, 1.173000),(1994.000000, 98.200000, 1.174000),(1995.000000, 99.800000, 1.205000),(1996.000000, 105.900000, 1.288000),(1997.000000, 105.800000, 1.291000),(1998.000000, 91.600000, 1.115000),(1999.000000, 100.100000, 1.221000),(2000.000000, 128.600000, 1.563000),(2001.000000, 124.000000, 1.531000),(2002.000000, 116.000000, 1.441000),(2003.000000, 135.100000, 1.638000),(2004.000000, 159.700000, 1.923000),(2005.000000, 194.700000, 2.338000),(2006.000000, 219.900000, 2.635000),(2007.000000, 237.959000, 2.849000),(2008.000000, 277.457000, 3.317000),(2009.000000, 201.555000, 2.401000),(2010.000000, 238.594000, 2.836000),(2011.000000, 301.694000, 3.577000),(2012.000000, 311.470000, 3.695000),(2013.000000, 302.577000, 3.584000),(2014.000000, 290.889000, 3.425000),(2015.000000, 212.007000, 2.510000),(2016.000000, 187.602000, 2.204000),(2017.000000, 211.770000, 2.469000),(2018.000000, 240.599000, 2.794000),(2019.000000, 232.003000, 2.698000),(2020.000000, 194.130000, 2.242000),(2021.000000, 264.017000, 3.133000),(2022.000000, 347.747000, 4.192000)','The prediction for gas prices in the year 2023 is: $4.34');
# INSERT INTO naive_bayes_predictions(prediction_id,query_identifier, input_data, prediction_info) VALUES (UUID(),'Top 3 Tech Jobs with most demand skills','','Most Demand Skills: Python, Java, React, Design
# Top Jobs for ''Tech'' Domain:
# Job Title: Software Engineer- Full Stack, URL: https://www.indeed.com/pagead/clk?mo=r&ad=-6NYlbfkN0DlG8eQ-E_eElbjQUAvh1Po9TMRdgcP7wLVIjfv-ORu9JgO3eKYfgel2s47D-0NuUWA6aOloEkMl1IgL2hYHZkjrYA-GpXd3cRH7rjuHefL0mCXddEQODNrSSFhAN_q86B
# wYK24GFcTVHJEC9nRhV0RUcKYvwRjXScpZ009kyw4OR_v2YKKc2InCPsjDB7hyOIo_Ax-9w4_N3v8Lnog5GjG_5b7R8kFUkqoYV7j9pKCQeBa2xGWX7B2SiBMOaOevJdezWYuZzF382QPrFDp1LQKNJHr7HgjSBX2MaK-XwFGOxJ2wseNKBQtxnGof6uDnmPTaY0UaxiOcDjdWqNTkH7szpsZI4rMdJE0Ivm
# XkoJdAhALHSscdNjYESTSZ-5H6mPpHLhVi0wM23pvfz3JZ-wZ-v0qzrEgXFswCZ_ozv4gyi7ImtWXRikSguzAilQVFSFY6LUCxOexrhW2Wjnlmb2x2Pe0zE8ZxslmEi5ffDuQ-1F4fWSMUjgqt3XVxc0NQ3i3Bs2pARaR_feRRS8Y9Q_yKr7UKDBkudRZysshrqCJB1msTCa3_8xMSlJbNb-GZToITyRSf2S
# FQlkEm3PqNzJIHYoArF9cmBYA_Fbbpo6lpJuDBvAGXLQEuIyLGdlN0yL-xBZb1YVhB6_PsIuqxoM7HORIxfTO5bjc9fmr0M7SUXWDsJUzVF8CiNsj_PmPGaZ-q56kEz58-h251XPKu8eQBEC53k-B_uU-1Q9wBh8-SztED0XY2HUvxzIwwEkkmFTh6Cjy-ssztoSf9vpS1TGdmQzNg2xAPU9nCRUz5g8p4WN
# Nyr3Y3MfnAmHgZ-6iKBfmMf7W0HQVyfNMJvFeTKD9RT5Zb0T-QCg-tfhLq8KMgHWO7APPnV0KQUDYct2MnioKi-HROkAVKOwwGOJl6mBZZ9wQg4c0dEWjzJjIjVCvC9KbY7SAVUphz5yYZ4lHASl_juQk7g26NewaVd2YgARC&xkcb=SoCT-_M3HcRK65RgW50IbzkdCdPP&p=3&fvj=0&vjs=3, Company: Motion Recruitment, Location: Location Not Found, Salary: $120,000 - $150,000 a year
# Description: Proficiency in a range of programming languages and technologies such as JavaScript, Python, Java, Ruby, or others, for both front-end and back-end development.
# Front-End Development: React, Angular, Vue.js and proficiency in HTML, CSS, and JavaScript.
# Back-End Development: Node.js, Django, Flask, Spring,
# Full Stack Proficiency: Ability to seamlessly work on both front-end and back-end components, understanding the interplay between them to build cohesive applications.
# Proficiency in a range of programming languages and technologies such as JavaScript, Python, Java, Ruby, or others, for both front-end and back-end development.
# Front-End Development: React, Angular, Vue.js and proficiency in HTML, CSS, and JavaScript.
# Back-End Development: Node.js, Django, Flask, Spring,
# Full Stack Proficiency: Ability to seamlessly work on both front-end and back-end components, understanding the interplay between them to build cohesive applications.
# Database Knowledge: Familiarity with database systems like SQL, NoSQL, and the ability to design efficient data models.
# Adaptability and Eagerness to Learn: Willingness to adapt to new technologies and learn continuously in a fast-paced environment.
# 50% front end
# 50% back end
# 50% front end
# 50% back end
# Design, develop, and deploy scalable, efficient, and high-quality software applications across the full tech stack.
# Design, develop, and deploy scalable, efficient, and high-quality software applications across the full tech stack.
# Ensure code quality through best practices, testing, and code reviews.
# Troubleshoot and debug issues to maintain the performance and reliability of software applications
# Healthcare insurance, 100% company-paid premiums
# Stock options
# Work life balance, vacation time
# Healthcare insurance, 100% company-paid premiums
# Stock options
# Work life balance, vacation time
# Matching Skills: 4 [design, Python, Java, React]
#
# Job Title: Software Developer, URL: https://www.indeed.com/rc/clk?jk=62253a0443b245d8&bb=qX68GcmwpUT_7GZdeTpvaQzYU8YaOhuxU9Fif5hbMMSyOikzvNT_fXC6oyzz7zq0&xkcb=SoBC67M3HcRKiURSZR0PbzkdCdPP&fccid=3c871c128fb1bf83&vjs=3, Company: LifeLens Technologies, Location: Ivyland, PA 18974, Salary: Salary Not Found
# Description: Java, JavaScript, Python or equivalent
# Understanding of HITRUST, GDPR, and/or HIPAA security compliance.
# Familiarity with React/Redux or equivalent interface frameworks.
# Familiarity with Docker container runtime engine.
# Helping support on-device application development and integration with back-end services.
# Bachelor’s Degree in Computer Science, Computer Engineering, or a related technical discipline
# Matching Skills: 3 [Python, Java, React]
#
# Job Title: URBN Platform Services Engineer, URL: https://www.indeed.com/rc/clk?jk=2cdad9fc6c68b401&bb=jcqAAeUUNaDqg-nH3cQUctIZ9ZoGo2H-L8MmmHGec9TsJ6NJ8bQwd-EeMKKAi_SE&xkcb=SoCn67M3HcRKspyU5p0KbzkdCdPP&fccid=a97e193d3b71f26b&vjs=3, Company: URBN, Location: Location Not Found, Salary: Salary Not Found
# Description: Design and development of microservices and event driven systems
# Develop new product features and advance the design of existing code for the Ecommerce applications
# Develop functional designs while interacting with enterprise architects, developers, business partners as well as end users
# React to possible production issues
# Enjoys working in a collaborative team setting
# Understanding of OO principles and design patterns
# Self-starter, ready to learn new concepts and technologies. Willing to take on responsibility for new features and develop them from inception to completion
# Proficient in Python/Golang
# Strong understanding of RESTful services design
# Strong understanding of asynchronous message processing technologies
# Enjoys working across teams and sharing knowledge with other engineers
# Familiarity with search solutions like Elastic Search, Solr
# Demonstrated expertise with Unit Testing frameworks and mocking libraries
# Familiar with Django or other ORM frameworks
# Exposure to data science and Machine learning technologies
# Understanding of cloud hosting platforms and infrastructure');

INSERT INTO naive_bayes_predictions (prediction_id, query_identifier, input_data, prediction_info)
VALUES (
           UUID(),
           'Top 3 Tech Jobs',
           'Software Release DevOps Engineer',
           'C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\SoftwareEng_top_jobs.json'
       );

INSERT INTO naive_bayes_predictions (prediction_id, query_identifier, input_data, prediction_info)
VALUES (
           UUID(),
           'Top 3 Law Jobs',
           'Law Related Skills', -- Replace with actual skills
           'C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\Law_top_jobs.json'
       );

INSERT INTO naive_bayes_predictions (prediction_id, query_identifier, input_data, prediction_info)
VALUES (
           UUID(),
           'Top 3 Business Jobs',
           'Business Related Skills', -- Replace with actual skills
           'C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\Business_top_jobs.json'
       );
# C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\Business_top_jobs.json
# C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\Law_top_jobs.json
# C:\\Users\\Public\\GoLandProjects\\JustAFork\\Nbc_output\\SoftwareEng_top_jobs.json

-- Inserting sample URLs into the URLs table
INSERT INTO urls (id, url, tags)
VALUES
    (UUID(), 'http://books.toscrape.com/', '{"tag1": "<a>"}');

-- Call to the procedure to populate log status codes
CALL populate_log_status_codes();

use goengine;

CALL user_registration('test_user', 'test_login', 'ADM', 'test_password', true);


call populate_log_status_codes();
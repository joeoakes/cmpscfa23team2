
# User Interaction Flow Instructions

This document outlines the flow of user interactions through the system components for different cases, from the initial user request to the final output delivered back to the user.

## E-commerce (Price Prediction)

- **User Interaction**: The user queries the API for a price prediction of a specific product.
- **Web UI to CUDA**: The Web UI processes the query, formats the request, and sends it to the CUDA component.
- **CUDA to CRAB**: CUDA determines the necessary data for the prediction and requests it from the CRAB component.
- **CRAB to DAL/SQL**: CRAB fetches historical price data, user reviews, and product metadata from the SQL database through the Data Access Layer (DAL).
- **SQL to CRAB**: The SQL database sends the requested data back to CRAB.
- **CUDA Processing**: CUDA processes the data with a machine learning model to predict the product's price.
- **Output to User**: The predicted price is sent back to the user through the Web UI.

## Social Media (Sentiment Analysis)

- **User Interaction**: The user requests sentiment analysis of a trending topic.
- **Web UI to CUDA**: The request is passed to CUDA.
- **CUDA to CRAB**: CUDA instructs CRAB to gather relevant social media posts.
- **CRAB to DAL/SQL**: CRAB retrieves social media data from the SQL database via DAL.
- **SQL to CRAB**: The data is returned to CRAB.
- **CUDA Processing**: CUDA performs sentiment analysis using NLP algorithms.
- **Output to User**: The sentiment analysis result is displayed to the user.

## News and Articles (Trend Detection)

- **User Interaction**: The user seeks the latest trends in a particular news category.
- **Web UI to CUDA**: The user's interest is translated into a CUDA query.
- **CUDA to CRAB**: CUDA requests recent news articles from CRAB.
- **CRAB to DAL/SQL**: CRAB pulls the latest articles from the SQL database.
- **SQL to CRAB**: Articles are sent back to CRAB.
- **CUDA Processing**: CUDA applies trend detection algorithms on the articles.
- **Output to User**: The detected trends are sent to the user.

## Real Estate (Market Value Prediction)

- **User Interaction**: The user requests an estimation of a property's market value.
- **Web UI to CUDA**: The request is processed and forwarded to CUDA.
- **CUDA to CRAB**: CUDA asks CRAB for relevant real estate listings.
- **CRAB to DAL/SQL**: CRAB fetches the requested data from the SQL database.
- **SQL to CRAB**: The real estate data is sent to CRAB.
- **CUDA Processing**: CUDA predicts the property's value using regression models.
- **Output to User**: The estimated value is provided to the user.

## Job Market (Industry Trend Analysis)

- **User Interaction**: The user inquires about the demand for a specific skill set.
- **Web UI to CUDA**: The inquiry goes to CUDA.
- **CUDA to CRAB**: CUDA requests job listings data from CRAB.
- **CRAB to DAL/SQL**: CRAB retrieves job data from the SQL database.
- **SQL to CRAB**: The data is returned to CRAB.
- **CUDA Processing**: CUDA analyzes the data for demand trends.
- **Output to User**: The analysis is presented to the user.

## Weather (Weather Pattern Forecast)

- **User Interaction**: The user asks for a weather forecast.
- **Web UI to CUDA**: The forecast request is routed to CUDA.
- **CUDA to CRAB**: CUDA requests historical weather data from CRAB.
- **CRAB to DAL/SQL**: CRAB pulls the needed weather data from the SQL database.
- **SQL to CRAB**: The database sends back the weather data.
- **CUDA Processing**: CUDA uses time series forecasting to predict the weather.
- **Output to User**: The weather forecast is delivered to the user.

# General Flow

1. The user makes a request via an API.
2. CUDA processes the logic and decides what data is needed.
3. CRAB retrieves the data from a SQL database through DAL.
4. CUDA processes and predicts using the retrieved data.
5. The results are returned to the user through the Web UI.

This flow represents the full cycle of data handling and processing from the user's request to the delivery of the machine learning-powered prediction.

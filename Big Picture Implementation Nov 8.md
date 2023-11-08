# Prediction Engine API Interaction Flow

## E-commerce (Price Prediction)
- **User Interaction**: Request for a price prediction of a product.
- **ML Algorithms**: Linear Regression, Random Forest.
- **Flow**:
  - Web UI → CUDA: Formats and sends the request.
  - CUDA → CRAB: Instructs data retrieval.
  - CRAB → DAL/SQL: Fetches historical price data and user reviews.
  - DAL/SQL → CRAB: Returns the requested data.
  - CUDA: Processes the data to predict the price.
  - User: Receives the predicted price.

## Social Media (Sentiment Analysis and Collaborative Filtering)
- **User Interaction**: Sentiment analysis request for a trending topic.
- **ML Algorithms**: NLP for Sentiment Analysis, Collaborative Filtering for personalized content recommendation.
- **Flow**:
  - Web UI → CUDA: Transmits the user's request.
  - CUDA → CRAB: Gathers relevant social media posts.
  - CRAB → DAL/SQL: Retrieves data from the database.
  - DAL/SQL → CRAB: Sends back the data.
  - CUDA: Performs sentiment analysis and applies collaborative filtering.
  - User: Gets sentiment results and content recommendations.

## News and Articles (Trend Detection)
- **User Interaction**: Query for latest trends in news.
- **ML Algorithms**: NLP for Topic Modeling, SVM for Classification.
- **Flow**:
  - Web UI → CUDA: Processes the query.
  - CUDA → CRAB: Requests recent news articles.
  - CRAB → DAL/SQL: Pulls the latest articles.
  - DAL/SQL → CRAB: Delivers the articles.
  - CUDA: Applies trend detection algorithms.
  - User: Receives the trending topics.

## Real Estate (Market Value Prediction)
- **User Interaction**: Market value estimation request for a property.
- **ML Algorithms**: Regression Trees, K-Nearest Neighbors.
- **Flow**:
  - Web UI → CUDA: Forwards the request.
  - CUDA → CRAB: Asks for real estate listings.
  - CRAB → DAL/SQL: Fetches property data.
  - DAL/SQL → CRAB: Returns the data.
  - CUDA: Predicts property value.
  - User: Obtains estimated market value.

## Job Market (Industry Trend Analysis)
- **User Interaction**: Inquiry about the demand for specific job skills.
- **ML Algorithms**: Time Series Analysis, Naive Bayes Classifier.
- **Flow**:
  - Web UI → CUDA: Sends the inquiry.
  - CUDA → CRAB: Requests job listing data.
  - CRAB → DAL/SQL: Retrieves job data.
  - DAL/SQL → CRAB: Sends back the data.
  - CUDA: Analyzes for demand trends.
  - User: Gets the trend analysis.

## Weather (Weather Pattern Forecast)
- **User Interaction**: Weather forecast request.
- **ML Algorithms**: ARIMA for Time Series Forecasting, Neural Networks.
- **Flow**:
  - Web UI → CUDA: Routes the forecast request.
  - CUDA → CRAB: Requests historical weather data.
  - CRAB → DAL/SQL: Pulls weather data.
  - DAL/SQL → CRAB: Delivers the data.
  - CUDA: Uses forecasting models to predict weather.
  - User: Receives the weather forecast.

Each flow is designed to be a seamless interaction from the user's request to the final output, leveraging machine learning algorithms for accurate and efficient predictions.

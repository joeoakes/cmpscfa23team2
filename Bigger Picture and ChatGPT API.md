# Prediction Engine API Interaction Flow with ChatGPT Integration

## E-commerce (Price Prediction)
- **ChatGPT Feature**: Interprets user's natural language queries for product details and clarifies prediction outcomes.
- **User Interaction**: Request for a price prediction of a product.
- **ML Algorithms**: Linear Regression, Random Forest.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT formats and sends the user's request to CUDA.
  - CUDA → CRAB: Instructs data retrieval.
  - CRAB → DAL/SQL: Fetches historical price data and user reviews.
  - DAL/SQL → CRAB: Returns the requested data.
  - CUDA: Processes the data to predict the price.
  - User: Receives the predicted price via ChatGPT.

## Social Media (Sentiment Analysis and Collaborative Filtering)
- **ChatGPT Feature**: Refines search queries for sentiment analysis and explains results to the user.
- **User Interaction**: Sentiment analysis request for a trending topic.
- **ML Algorithms**: NLP for Sentiment Analysis, Collaborative Filtering for personalized content recommendation.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT transmits the refined user's request to CUDA.
  - CUDA → CRAB: Gathers relevant social media posts.
  - CRAB → DAL/SQL: Retrieves data from the database.
  - DAL/SQL → CRAB: Sends back the data.
  - CUDA: Performs sentiment analysis and applies collaborative filtering.
  - User: Gets sentiment results and content recommendations via ChatGPT.

## News and Articles (Trend Detection)
- **ChatGPT Feature**: Assists in formulating user queries for trend detection and summarizes output.
- **User Interaction**: Query for the latest trends in news.
- **ML Algorithms**: NLP for Topic Modeling, SVM for Classification.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT processes and forwards the query to CUDA.
  - CUDA → CRAB: Requests recent news articles.
  - CRAB → DAL/SQL: Pulls the latest articles.
  - DAL/SQL → CRAB: Delivers the articles.
  - CUDA: Applies trend detection algorithms.
  - User: Receives the trending topics via ChatGPT.

## Real Estate (Market Value Prediction)
- **ChatGPT Feature**: Gathers additional property details from users for accurate predictions.
- **User Interaction**: Market value estimation request for a property.
- **ML Algorithms**: Regression Trees, K-Nearest Neighbors.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT forwards the detailed request to CUDA.
  - CUDA → CRAB: Asks for real estate listings.
  - CRAB → DAL/SQL: Fetches property data.
  - DAL/SQL → CRAB: Returns the data.
  - CUDA: Predicts property value.
  - User: Obtains estimated market value via ChatGPT.

## Job Market (Industry Trend Analysis)
- **ChatGPT Feature**: Clarifies inquiries about job market trends and interprets analysis results.
- **User Interaction**: Inquiry about the demand for specific job skills.
- **ML Algorithms**: Time Series Analysis, Naive Bayes Classifier.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT sends a clarified inquiry to CUDA.
  - CUDA → CRAB: Requests job listing data.
  - CRAB → DAL/SQL: Retrieves job data.
  - DAL/SQL → CRAB: Sends back the data.
  - CUDA: Analyzes for demand trends.
  - User: Gets the trend analysis via ChatGPT.

## Weather (Weather Pattern Forecast)
- **ChatGPT Feature**: Collects detailed forecast requests and communicates predictions.
- **User Interaction**: Weather forecast request.
- **ML Algorithms**: ARIMA for Time Series Forecasting, Neural Networks.
- **Flow**:
  - Web UI → ChatGPT → CUDA: ChatGPT routes the detailed forecast request to CUDA.
  - CUDA → CRAB: Requests historical weather data.
  - CRAB → DAL/SQL: Pulls weather data.
  - DAL/SQL → CRAB: Delivers the data.
  - CUDA: Uses forecasting models to predict weather.
  - User: Receives the weather forecast via ChatGPT.

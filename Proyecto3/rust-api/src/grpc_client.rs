pub mod weathertweet {
    tonic::include_proto!("wethertweet"); // <- coincide con package del proto
}

use weathertweet::weather_tweet_service_client::WeatherTweetServiceClient;
use weathertweet::WeatherTweetRequest;
use weathertweet::WeatherTweetResponse;

pub async fn send_to_grpc(
    grpc_url: String,
    municipality: i32,
    temperature: i32,
    humidity: i32,
    weather: i32,
) -> Result<WeatherTweetResponse, Box<dyn std::error::Error>> {
    let mut client = WeatherTweetServiceClient::connect(grpc_url).await?;

    let request = tonic::Request::new(WeatherTweetRequest {
        municipality,
        temperature,
        humidity,
        weather,
    });

    let response = client.send_tweet(request).await?;
    Ok(response.into_inner())
}
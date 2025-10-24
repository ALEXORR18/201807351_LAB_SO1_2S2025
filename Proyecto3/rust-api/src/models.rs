use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, Clone)] // <- agregamos Clone
pub struct WeatherInput {
    pub municipality: String,
    pub temperature: i32,
    pub humidity: i32,
    pub weather: String,
}

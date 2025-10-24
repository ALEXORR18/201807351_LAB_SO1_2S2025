use actix_web::{web, App, HttpServer, HttpResponse, Result, middleware::Logger};
use env_logger;
use log::info;
use chrono::Utc;
use dotenvy;

mod models;
mod grpc_client;

use models::WeatherInput;

fn map_municipality(name: &str) -> i32 {
    match name.to_lowercase().as_str() {
        "mixco" => 1,
        "guatemala" => 2,
        "amatitlan" => 3,
        "chinautla" => 4,
        _ => 0, // municipalities_unknown
    }
}

fn map_weather(weather: &str) -> i32 {
    match weather.to_lowercase().as_str() {
        "sunny" | "soleado" => 1,
        "cloudy" | "nublado" => 2,
        "rainy" | "lluvioso" => 3,
        "foggy" | "neblina" => 4,
        _ => 0, // weathers_unknown
    }
}

async fn root() -> Result<HttpResponse> {
    info!("Endpoint raíz accedido");
    Ok(HttpResponse::Ok().body("¡Hola! Bienvenido a la API del Clima"))
}

async fn post_clima(datos: web::Json<WeatherInput>) -> Result<HttpResponse> {
    info!("POST /clima con datos: {:?}", datos);

    let datos_clone = datos.into_inner();
    let grpc_url = std::env::var("GRPC_SERVER_URL")
        .unwrap_or("http://127.0.0.1:50051".to_string());

    let municipio_enum = map_municipality(&datos_clone.municipality);
    let weather_enum = map_weather(&datos_clone.weather);

    // Clonar nuevamente para usar en json de respuesta
    let datos_response = datos_clone.clone();

    let _ = tokio::spawn(async move {
        if let Err(e) = grpc_client::send_to_grpc(
            grpc_url,
            municipio_enum,
            datos_clone.temperature,
            datos_clone.humidity,
            weather_enum
        ).await {
            eprintln!("Error gRPC: {:?}", e);
        }
    });

    Ok(HttpResponse::Ok().json(serde_json::json!({
        "datos_recibidos": datos_response,
        "status": "ok",
        "timestamp": Utc::now()
    })))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenvy::dotenv().ok();

    match std::env::var("GRPC_SERVER_URL") {
        Ok(url) => println!("Variable GRPC_SERVER_URL cargada: {}", url),
        Err(_) => println!("GRPC_SERVER_URL no encontrada, usando valor por defecto"),
    }

    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));
    let port = 8080;
    info!("Rust REST API corriendo en http://0.0.0.0:{}", port);

    HttpServer::new(|| {
        App::new()
            .wrap(Logger::default())
            .wrap(Logger::new("%a %{User-Agent}i"))
            .route("/", web::get().to(root))
            .route("/clima", web::post().to(post_clima))
    })
    .bind(("0.0.0.0", port))?
    .run()
    .await
}

from locust import HttpUser, task, between
import json
import random

# Lista de municipios y su correspondiente valor enum según el proto
MUNICIPALIDADES = ["mixco", "guatemala", "amatitlan", "chinautla"]
MUNICIPALIDADES_NUM = [1, 2, 3, 4]

# Lista de climas y su correspondiente valor enum según el proto
CLIMA = ["sunny", "cloudy", "rainy", "foggy"]
CLIMA_NUM = [1, 2, 3, 4]

# Funciones para generar temperatura y humedad realista según el clima
def temp(clima):
    if clima == "sunny":    return random.randint(25, 35)
    elif clima == "cloudy": return random.randint(20, 25)
    elif clima == "rainy":  return random.randint(15, 20)
    else:                   return random.randint(10, 15)  # foggy

def humidity(clima):
    if clima == "sunny":    return random.randint(30, 50)
    elif clima == "cloudy": return random.randint(50, 70)
    elif clima == "rainy":  return random.randint(70, 90)
    else:                   return random.randint(80, 100)  # foggy

class ClimaAPITest(HttpUser):
    # Tiempo de espera entre tareas (en segundos)
    wait_time = between(5, 7)

    def on_start(self):
        """Se ejecuta al iniciar cada usuario simulado"""
        print("Iniciando usuario Locust para pruebas de carga...")
        # Verifica que la API está viva
        response = self.client.get("/")
        if response.status_code == 200:
            print("Endpoint raíz activo ✅")
        else:
            print(f"Error en endpoint raíz: {response.status_code} ❌")

    @task(1)
    def test_root(self):
        """Prueba ligera al endpoint raíz"""
        response = self.client.get("/")
        if response.status_code == 200:
            print("Endpoint raíz OK")
        else:
            print(f"Error {response.status_code} en endpoint raíz")

    @task(10)
    def test_clima(self):
        """Envía tweets simulados al endpoint /clima"""
        climaRandom = random.choice(CLIMA)
        municipioRandom = random.choice(MUNICIPALIDADES)

        tweetGenerado = {
            "municipality": MUNICIPALIDADES[MUNICIPALIDADES.index(municipioRandom)],
            "temperature": temp(climaRandom),
            "humidity": humidity(climaRandom),
            "weather": CLIMA[CLIMA.index(climaRandom)]
        }

        # Post al endpoint de la API Rust
        with self.client.post("/clima", json=tweetGenerado, catch_response=True) as response:
            if response.status_code == 200:
                print(f"Tweet enviado correctamente: {json.dumps(tweetGenerado)} ✅")
                response.success()
            else:
                print(f"Error al enviar tweet: {response.status_code} ❌")
                print(f"Datos del tweet: {json.dumps(tweetGenerado)}")
                response.failure(f"Error {response.status_code}")

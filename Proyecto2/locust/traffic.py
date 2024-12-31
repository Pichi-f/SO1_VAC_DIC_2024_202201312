from locust import HttpUser, between, task
from random import randrange
import json

class ReadFile():
    def __init__(self):
        self.calificaciones = []

    def getData(self):
        size = len(self.calificaciones)
        if size > 0:
            index = randrange(0, size - 1) if size > 1 else 0
            return self.calificaciones.pop(index)
        else:
            print("No hay más datos")
            return None

    def loadFile(self):
        try:
            with open("data.json", "r", encoding="utf-8") as file:
                self.calificaciones = json.loads(file.read())
        except Exception as e:
            print(f"Error al cargar el archivo: {e}")

class trafficData(HttpUser):
    wait_time = between(1, 2)
    reader = ReadFile()
    reader.loadFile()

    @task
    def sendMessage(self):
        data = self.reader.getData()
        if data is not None:
            res = self.client.post("/insert", json=data)
            response = res.json()
            print(response)
        else:
            print("No hay más datos")
            self.stop(True)
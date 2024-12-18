# Manual Técnico

## Introducción
Este manual técnico explica detalladamente los componentes utilizados en un proyecto que integra monitoreo, gestión de datos, desarrollo de software y despliegue en contenedores, todos ellos desplegados en instancias de máquinas virtuales (VM) en Google Cloud Platform (GCP). Los componentes principales son:

- Grafana
- Base de datos MySQL
- Módulos de kernel en C
- Agente en Golang
- API desarrollada con Node.js
- Contenedores y Docker

### Arquitectura General
- Una instancia de VM que contiene los contenedores de Grafana, MySQL y la API Node.js. Esta instancia se encarga del monitoreo y gestión de datos.
- Un grupo de instancias de VM con los módulos de kernel en C, administrados por un agente desarrollado en Golang que opera en un contenedor. Este agente recopila métricas y las envía a la API Node.js.

---

## 1. Grafana
Grafana es una plataforma de análisis y monitoreo de código abierto que permite crear paneles interactivos para visualizar métricas en tiempo real.

### Características:
- Soporte para múltiples fuentes de datos, incluyendo MySQL.
- Personalización de paneles con gráficos, tablas y alertas.
- Capacidad de integración con sistemas como Prometheus, InfluxDB, y ElasticSearch.

### Instalación:
1. Configurar un contenedor Docker con Grafana en la VM:
   ```bash
   docker run -d -p 3000:3000 --name grafana grafana/grafana
   ```
2. Acceder a la interfaz en el navegador utilizando `http://<IP_VM_GCP>:3000`.

### Configuración:
- Configurar MySQL como fuente de datos en Grafana.
- Crear paneles personalizados para visualizar métricas provenientes de la base de datos MySQL.

---

## 2. Base de datos MySQL
MySQL es un sistema de gestión de bases de datos relacional, utilizado para almacenar y gestionar los datos recopilados del sistema.

### Características:
- Transacciones ACID.
- Soporte para consultas SQL avanzadas.
- Integración con múltiples aplicaciones.

### Instalación:
1. Configurar un contenedor Docker con MySQL en la VM:
   ```bash
   docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root mysql
   ```
2. Configurar MySQL para aceptar conexiones externas desde la API Node.js.

### Configuración:
- Crear tablas para almacenar métricas enviadas por el agente Golang a través de la API Node.js.
- Configurar conexiones seguras entre MySQL y las demás aplicaciones.

---

## 3. Módulos de Kernel en C
Un módulo de kernel es un programa que extiende la funcionalidad del kernel de Linux sin necesidad de recompilarlo.

### Funcionalidad:
- Monitorizar métricas del sistema, como uso de CPU, memoria y operaciones de I/O.
- Enviar datos al agente Golang para su procesamiento y reenvío a la API Node.js.

### Desarrollo:
1. Crear el archivo fuente del módulo (`module.c`). Ejemplo básico:
   ```c
   #include <linux/module.h>
   #include <linux/kernel.h>

   static int _insert(void) {
        proc_create("cpu_202201312", 0, NULL, &operaciones);
        printk(KERN_INFO "Creado el archivo /proc/cpu_202201312\n - Clase 4\n");
        return 0;
    }
   ```
2. Compilar el módulo:
   ```bash
   make -C /lib/modules/$(uname -r)/build M=$(pwd) modules
   ```
3. Cargar el módulo:
   ```bash
   sudo insmod module.ko
   ```

### Integración:
- Usar el agente Golang para recopilar métricas expuestas por el módulo a través de `procfs` o syscalls personalizados.

---

## 4. Agente en Golang
El agente en Golang se ejecuta en contenedores desplegados en un grupo de instancias de VM. Su función principal es recopilar métricas desde los módulos de kernel y enviarlas a la API Node.js.

### Funcionalidad:
- Ejecutar tareas de monitoreo periódicas.
- Comunicar datos a la API Node.js mediante solicitudes HTTP.

### Instalación:
1. Desarrollar el agente:
   ```go
   package main

   import (
       "encoding/json"
       "net/http"
       "os"
   )

   type Process struct {
	Pid    int     `json:"pid"`
	Name   string  `json:"name"`
	User   int     `json:"user"`
	State  int     `json:"state"`
	Ram    float64 `json:"ram"`
	Father int     `json:"father"`
   }

   type Cpu struct {
      Usage     float64   `json:"percentage_used"`
      Pc        int       `json:"pc"`
      Processes []Process `json:"tasks"`
   }

   type Ram struct {
      Total float64 `json:"total_ram"`
      Free  float64 `json:"free_ram"`
      Used  float64 `json:"used_ram"`
      Perc  float64 `json:"percentage_used"`
      Pc    int     `json:"pc"`
   }
   ```
2. Empaquetar el agente en un contenedor Docker:
   ```dockerfile
   FROM golang:alpine AS build
   RUN apk add --no-cache git
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN go build -o my-app
   FROM alpine:latest
   WORKDIR /app
   COPY --from=build /app/my-app ./
   ENTRYPOINT ["./my-app"]
   EXPOSE 5200
   ```
3. Desplegar el contenedor en cada instancia del grupo de VM.

---

## 5. API con Node.js
La API desarrollada con Node.js actúa como intermediario, procesando métricas enviadas por el agente y almacenándolas en la base de datos MySQL.

### Estructura:
- Framework: Express.js
- Base de datos: MySQL mediante `mysql2` o `sequelize`.

### Desarrollo:
1. Inicializar el proyecto:
   ```bash
   npm init -y
   npm install express mysql2
   ```
2. Crear un servidor básico:
   ```javascript
   const mysql = require('mysql2/promise');
   const express = require('express');
   const dotenv = require('dotenv');

   // Cargar variables de entorno desde .env
   dotenv.config();

   const connection = mysql.createPool({
      host: process.env.DB_HOST,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD,
      database: process.env.DB_NAME,
      port: process.env.DB_PORT,
   });

   const app = express();
   const PORT = process.env.PORT || 4000;

   // Middleware para parsear JSON
   app.use(express.json());

   // Ruta básica
   app.get('/', (req, res) => {
      res.send('¡Hola, API funcionando!');
   });
   ```

### Despliegue:
- Implementar la API en la misma instancia de VM que contiene Grafana y MySQL, dentro de un contenedor Docker.

---

## 6. Contenedores y Docker
Docker es una plataforma que permite empaquetar aplicaciones y sus dependencias en contenedores portables.

### Instalación:
1. Instalar Docker en todas las instancias de VM:
   ```bash
   sudo apt-get install docker.io
   ```
2. Crear los contenedores para cada componente (Grafana, MySQL, API Node.js, Agente Golang).

### Despliegue:
- Utilizar Google Cloud Platform para gestionar las instancias y los contenedores.
- Orquestar los contenedores en diferentes VMs para garantizar escalabilidad y modularidad.

---

## Conclusión
Este proyecto utiliza una arquitectura distribuida en Google Cloud Platform, combinando Grafana, MySQL, módulos de kernel, un agente Golang, Node.js y Docker. Esta solución asegura un monitoreo eficiente y una infraestructura escalable y modular.
DROP TABLE IF EXISTS grafana_db.RAM, grafana_db.CPU;
CREATE TABLE IF NOT EXISTS grafana_db.RAM (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usado FLOAT,
    disponible FLOAT,
    tiempo DATETIME,
    pc VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS grafana_db.CPU (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usado FLOAT,
    disponible FLOAT,
    tiempo DATETIME,
    pc VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS grafana_db.PROCESOS (
    pid INT PRIMARY KEY,
    nombre VARCHAR(255),
    usuario VARCHAR(255),
    estado VARCHAR(255),
    ram_usada VARCHAR(255),
    pc VARCHAR(255)
);

SELECT * FROM grafana_db.RAM;
SELECT * FROM grafana_db.CPU;

SELECT * FROM grafana_db.PROCESOS;

TRUNCATE TABLE grafana_db.RAM;
TRUNCATE TABLE grafana_db.CPU;

TRUNCATE TABLE grafana_db.PROCESOS;
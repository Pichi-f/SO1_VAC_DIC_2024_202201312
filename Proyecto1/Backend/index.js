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

app.post('/cpu', async (req, res) => {
    const {percentage_used, pc} = req.body;
    const query = 'INSERT INTO CPU(usado, disponible, tiempo, pc) VALUES (?, ?, CURTIME(), ?);';
    await connection.query(query, [percentage_used, 100 - percentage_used, pc]);
    console.log(req.body);
    res.send('¡CPU recibida!');
});

app.post('/ram', async (req, res) => {
    const {percentage_used, pc} = req.body;
    const query = 'INSERT INTO RAM(usado, disponible, tiempo, pc) VALUES (?, ?, CURTIME(), ?);';
    await connection.query(query, [percentage_used, 100 - percentage_used, pc]);
    console.log(req.body);
    res.send('RAM recibida!');
});

// Iniciar servidor
app.listen(PORT, () => {
    console.log(`Servidor corriendo en http://localhost:${PORT}`);
});
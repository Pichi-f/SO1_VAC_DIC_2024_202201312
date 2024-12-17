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
    let ipAddress = req.header('x-forwarded-for')  || req.socket.remoteAddress;
    if (ipAddress.substr(0, 7) == "::ffff:") {
        ipAddress = ipAddress.substr(7)
    }
    const {percentage_used, tasks} = req.body;
    const query = 'INSERT INTO CPU(usado, disponible, tiempo, pc) VALUES (?, ?, CURTIME(), ?);';
    await connection.query(query, [percentage_used, 100 - percentage_used, ipAddress]);
    var campos = "";
    var valores = [];
    tasks.forEach(task => {
        const {pid, name, user, state, ram} = task;
        campos +=(campos!==""?",":"")+ "(?, ?, ?, ?, ?, ?)";
        valores.push(pid, name, user, state, ram, ipAddress);
    });
    const query2 = `INSERT INTO grafana_db.PROCESOS (pid, nombre, usuario, estado, ram_usada, pc)
VALUES
    ${campos}
ON DUPLICATE KEY UPDATE
    nombre = VALUES(nombre),
    usuario = VALUES(usuario),
    estado = VALUES(estado),
    ram_usada = VALUES(ram_usada), 
    pc = VALUES(pc);`;
    await connection.query(query2, valores);
    res.send('¡CPU recibida!');
});

app.post('/ram', async (req, res) => {
    let ipAddress = req.header('x-forwarded-for')  || req.socket.remoteAddress;
    if (ipAddress.substr(0, 7) == "::ffff:") {
        ipAddress = ipAddress.substr(7)
    }
    const {percentage_used} = req.body;
    const query = 'INSERT INTO RAM(usado, disponible, tiempo, pc) VALUES (?, ?, CURTIME(), ?);';
    await connection.query(query, [percentage_used, 100 - percentage_used, ipAddress]);
    res.send('RAM recibida!');
});

// Iniciar servidor
app.listen(PORT, () => {
    console.log(`Servidor corriendo en http://localhost:${PORT}`);
});
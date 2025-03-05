const express = require("express");
const mysql = require("mysql2/promise");
const client = require("prom-client");

const app = express();
const port = process.env.PORT || 3000;
const replicaId = process.env.HOSTNAME || "local-" + Math.floor(Math.random() * 1000);
const capstone = process.env.CAPSTONE || "No capstone value set";

const dbConfig = {
  host: process.env.DB_HOST || "mysql",
  user: process.env.DB_USER || "root",
  password: process.env.DB_PASSWORD || "root",
  database: process.env.DB_NAME || "testdb",
};

// Prometheus Metrics Setup
const register = new client.Registry();
client.collectDefaultMetrics({ register });
const dbQueryCounter = new client.Counter({
  name: "db_insert_total",
  help: "Total number of database inserts",
  registers: [register],
});
const dbQueryGauge = new client.Gauge({
  name: "db_last_insert_value",
  help: "Last inserted value",
  registers: [register],
});

async function initDB() {
  const conn = await mysql.createConnection(dbConfig);
  await conn.execute(`
    CREATE TABLE IF NOT EXISTS entries (
      id INT AUTO_INCREMENT PRIMARY KEY,
      replica VARCHAR(255) NOT NULL,
      value INT NOT NULL,
      timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `);
  conn.end();
}

async function insertEntry() {
  try {
    const conn = await mysql.createConnection(dbConfig);
    const value = Math.floor(Math.random() * 100);
    await conn.execute("INSERT INTO entries (replica, value) VALUES (?, ?)", [replicaId, value]);
    conn.end();
    dbQueryCounter.inc(); // Increase counter metric
    dbQueryGauge.set(value); // Set last insert value
  } catch (error) {
    console.error("Database Insert Error:", error);
  }
}

// Start the automatic insertion every second
setInterval(insertEntry, 1000);

app.get("/", async (req, res) => {
  try {
    const conn = await mysql.createConnection(dbConfig);
    const [rows] = await conn.execute("SELECT * FROM entries ORDER BY id DESC LIMIT 3");
    conn.end();
    res.json({
      message: "Last 3 Entries",
      replica: replicaId,
      data: rows,
      capstone: capstone,
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Prometheus Metrics Endpoint
app.get("/metrics", async (req, res) => {
  res.set("Content-Type", register.contentType);
  res.end(await register.metrics());
});

app.listen(port, () => {
  console.log(`Replica ${replicaId} running on port ${port}`);
  initDB();
});

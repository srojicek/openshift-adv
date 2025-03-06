const express = require("express");
const mysql = require("mysql2/promise");

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

async function initDB() {
  const conn = await mysql.createConnection(dbConfig);
  await conn.execute(`
    CREATE TABLE IF NOT EXISTS entries (
      id INT AUTO_INCREMENT PRIMARY KEY,
      replica VARCHAR(255) NOT NULL,
      value INT NOT NULL
    )
  `);
  await conn.execute("INSERT INTO entries (replica, value) VALUES (?, ?)", [
    replicaId,
    Math.floor(Math.random() * 100),
  ]);
  conn.end();
}

app.get("/", async (req, res) => {
  try {
    const conn = await mysql.createConnection(dbConfig);
    const [rows] = await conn.execute("SELECT * FROM entries ORDER BY id DESC LIMIT 3");
    conn.end();
    res.json({
      message: "All Entries",
      replica: replicaId,
      data: rows,
      capstone: capstone
    });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

app.listen(port, () => {
  console.log(`Replica ${replicaId} running on port ${port}`);
  initDB();
});
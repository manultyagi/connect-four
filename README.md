# 🎮 4 in a Row — Real-Time Multiplayer Game  
### Backend Engineering Intern Assignment

A real-time, backend-driven implementation of the classic **4 in a Row (Connect Four)** game.

The game supports:
- **Player vs Player**
- **Player vs Competitive Bot**
- **Real-time gameplay via WebSockets**
- **Kafka-based analytics**
- **Persistent leaderboard backed by PostgreSQL**

---

## 🧠 Overview

The primary focus of this project is **backend correctness**, **real-time systems**, and **clean architecture**, rather than UI polish.

### Key Highlights
- Real-time multiplayer gameplay using **WebSockets**
- Automatic **bot fallback** if no opponent joins within **10 seconds**
- Competitive bot that blocks wins and builds winning paths
- **Kafka-based analytics**, fully decoupled from gameplay
- Persistent leaderboard using **PostgreSQL**
- Fully **hosted live** for evaluation

---

## 🌍 Live Demo

🔗 **Live Application URL**  
👉 https://connect-four-game.onrender.com

> Kafka analytics and database persistence are demonstrated locally.  
> In production, gameplay runs independently for stability.

---

## 📦 Tech Stack

### Backend
- **Go (Golang)**
- Gorilla WebSocket
- In-memory game state
- PostgreSQL (optional, feature-flagged)
- Kafka (optional, feature-flagged)

### Frontend
- Vanilla **HTML, CSS, JavaScript**
- Served directly by the Go server

### Analytics
- Kafka (KRaft mode, no Zookeeper)
- Kafka consumer for analytics aggregation

---

## 🕹 Game Features

### 1️⃣ Player Matchmaking
- Players enter a username and wait for an opponent
- If no opponent joins within **10 seconds**, a **competitive bot** starts the game automatically
- If another player joins within 10 seconds, the game starts as **Player vs Player**

---

### 2️⃣ Competitive Bot
The bot:
- Plays valid moves only
- Blocks the opponent’s immediate winning moves
- Attempts to create its own winning opportunities
- Responds quickly and deterministically (not random)

---

### 3️⃣ Real-Time Gameplay
- Turn-based gameplay using **WebSockets**
- Both players see updates instantly after every move
- The **server is authoritative** over game state and turns

---

### 4️⃣ Game State Handling
- Active games are stored **in-memory**
- Completed games can be stored in **PostgreSQL**
- Database is optional and **feature-flagged** for production safety

---

### 5️⃣ Leaderboard
- Tracks number of games won per player
- Exposed via the `/leaderboard` endpoint
- Displayed on the frontend
- Uses PostgreSQL when enabled

---

### 6️⃣ Kafka Analytics (Bonus)
Analytics are fully **decoupled from gameplay** using Kafka.

Tracked metrics include:
- Average game duration
- Most frequent winners
- Games per day / hour
- User-specific statistics

> Kafka is disabled in production but fully functional locally.

---

## 🚀 Running the App Locally

### 🔹 Prerequisites
- Go **≥ 1.24**
- PostgreSQL
- Kafka (KRaft mode, optional)

---

### ▶️ Run Without Kafka & Database (Gameplay Only)


go run .

- Bot and PvP gameplay work

- Leaderboard is disabled

- Kafka analytics disabled

### ▶️ Run With PostgreSQL (Leaderboard Enabled)
1️⃣ Start PostgreSQL
- Create a database and tables:


CREATE TABLE players (
  username TEXT PRIMARY KEY,
  wins INT
);

CREATE TABLE games (
  id SERIAL PRIMARY KEY,
  player1 TEXT,
  player2 TEXT,
  winner TEXT,
  moves INT
);

2️⃣ Enable DB and start server

set ENABLE_DB=true

go run .

### ▶️ Run With Kafka Analytics
1️⃣ Start Kafka (KRaft mode)

bin/windows/kafka-server-start.bat config/kraft/server.properties

- Create topic (one-time):


bin/windows/kafka-topics.bat --create \
  --topic game-events \
  --bootstrap-server localhost:9092 \
  --partitions 1 \
  --replication-factor 1
  
2️⃣ Start Analytics Consumer

cd analytics

go run .

3️⃣ Start Game Server

set ENABLE_KAFKA=true

set ENABLE_DB=true

go run .

-Analytics output will appear in the consumer terminal.


### 🧪 Production Notes
- Kafka and DB are feature-flagged

- Core gameplay runs independently

- Mirrors real-world production architecture

### 📂 Project Structure
connect-four/

├── analytics/

│   └── main.go

├── public/

│   ├── index.html

│   ├── app.js

│   ├── leaderboard.html

│   └── style.css

├── server.go

├── game.go

├── bot.go

├── db.go

├── kafka_producer.go

├── main.go

├── go.mod

└── README.md
### 🏁 Summary
This project demonstrates:

- Real-time backend systems

- WebSocket communication

- Matchmaking and bot logic

- Decoupled analytics with Kafka

- Production-safe architecture

- Clean, testable Go code

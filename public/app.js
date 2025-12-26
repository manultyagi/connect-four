let previousBoard = null;
let ws;
let myTurn = false;
let myPlayerNumber = null;

const boardDiv = document.getElementById("board");
const status = document.getElementById("status");

// Create empty board immediately
function createEmptyBoard() {
    const board = Array.from({ length: 6 }, () => Array(7).fill(0));
    renderBoard(board);
}

function connect() {
    const username = document.getElementById("username").value.trim();
    if (!username) return;

    createEmptyBoard();

    const protocol = window.location.protocol === "https:" ? "wss" : "ws";
    const wsUrl = `${protocol}://${window.location.host}/ws`;
    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        ws.send(JSON.stringify({ username }));
        status.innerText = "Connected as " + username;

        // 🔑 CRITICAL FIX:
        // Allow Player 1 to make the first move immediately.
        // Server will still validate turns.
        myTurn = true;
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.status) {
            status.innerText = data.status;
        }

        if (data.playerNumber) {
            myPlayerNumber = data.playerNumber;
        }

        if (data.board) {
            renderBoard(data.board);
        }

        if (data.turn !== undefined && myPlayerNumber !== null) {
            myTurn = data.turn === myPlayerNumber;
        }

        if (data.winner && data.winner !== 0) {
            setTimeout(() => {
                alert(
                    data.winner === -1
                        ? "Draw!"
                        : data.winner === myPlayerNumber
                            ? "You win!"
                            : "You lose!"
                );
            }, 400);
        }
    };

    ws.onerror = () => {
        status.innerText = "WebSocket connection failed";
        myTurn = false;
    };

    ws.onclose = () => {
        myTurn = false;
    };
}

function renderBoard(board) {
    boardDiv.innerHTML = "";

    for (let row = 0; row < 6; row++) {
        for (let col = 0; col < 7; col++) {
            const cell = document.createElement("div");
            cell.className = "cell";

            if (board[row][col] === 1) cell.classList.add("player1");
            if (board[row][col] === 2) cell.classList.add("player2");

            if (
                previousBoard &&
                previousBoard[row][col] === 0 &&
                board[row][col] !== 0
            ) {
                cell.classList.add("drop");
            }

            cell.onclick = () => {
                if (!myTurn) return;
                myTurn = false;
                ws.send(JSON.stringify({ column: col }));
            };

            boardDiv.appendChild(cell);
        }
    }

    previousBoard = JSON.parse(JSON.stringify(board));
}

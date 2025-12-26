let previousBoard = null;
let ws;
let myTurn = false;
let myPlayerNumber = null;

const boardDiv = document.getElementById("board");
const status = document.getElementById("status");

function connect() {
    const username = document.getElementById("username").value.trim();
    if (!username) return;

    ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
        ws.send(JSON.stringify({ username }));
        status.innerText = "Connected as " + username;
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.status) status.innerText = data.status;
        if (data.playerNumber) myPlayerNumber = data.playerNumber;
        if (data.board) renderBoard(data.board);
        if (data.turn) myTurn = data.turn === myPlayerNumber;

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
}

function renderBoard(board) {
    boardDiv.innerHTML = "";

    for (let r = 0; r < 6; r++) {
        for (let c = 0; c < 7; c++) {
            const cell = document.createElement("div");
            cell.className = "cell";

            if (board[r][c] === 1) cell.classList.add("player1");
            if (board[r][c] === 2) cell.classList.add("player2");

            if (previousBoard && previousBoard[r][c] === 0 && board[r][c] !== 0) {
                cell.classList.add("drop");
            }

            cell.onclick = () => {
                if (!myTurn) return;
                myTurn = false;
                ws.send(JSON.stringify({ column: c }));
            };

            boardDiv.appendChild(cell);
        }
    }

    previousBoard = JSON.parse(JSON.stringify(board));
}

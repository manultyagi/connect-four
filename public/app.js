let previousBoard = null;
let ws;
let myTurn = false;
let myPlayerNumber = null;

const boardDiv = document.getElementById("board");
const status = document.getElementById("status");

function renderEmptyBoard() {
    const board = Array.from({ length: 6 }, () => Array(7).fill(0));
    renderBoard(board);
}

function connect() {
    const username = document.getElementById("username").value.trim();
    if (!username) return;

    renderEmptyBoard();

    const protocol = location.protocol === "https:" ? "wss" : "ws";
    ws = new WebSocket(`${protocol}://${location.host}/ws`);

    ws.onopen = () => {
        ws.send(JSON.stringify({ username }));
        status.innerText = "Connecting...";
    };

    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.status) {
            status.innerText = data.status;
        }

        if (data.playerNumber !== undefined) {
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
            }, 300);
        }
    };

    ws.onclose = () => {
        myTurn = false;
        status.innerText = "Disconnected";
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

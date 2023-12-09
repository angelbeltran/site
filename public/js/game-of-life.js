const gridDiv = document.querySelector("#game");

// Build DOM and initial state

gridDiv.style.border = 'solid 1px black';
gridDiv.style.display = 'flex';
gridDiv.style.flexDirection = 'row';

// the grid is composed of columns
const gridElements = Array(100);
let gridState = Array(100);

function initGrid() {
	for (let x = 0; x < 100; x += 1) {
		gridElements[x] = Array(100);
		gridState[x] = Array(100);

		const colDiv = document.createElement("div");
		colDiv.style.flex = '1% 1 1';
		colDiv.style.display = 'flex';
		colDiv.style.flexDirection = 'column';

		for (let y = 0; y < 100; y += 1) {
			const isAlive = (Math.random() * 2) >= 1;
			gridState[x][y] = isAlive;

			const squareDiv = document.createElement("div");
			gridElements[x][y] = squareDiv;

			squareDiv.style.border = 'solid 1px black';
			squareDiv.style.flex = '1% 1 1';

			colDiv.append(squareDiv);
		}

		gridDiv.append(colDiv);
	}
}

function updateGridElementsColor() {
	gridElements.forEach((col, x) => {
		col.forEach((cellDiv, y) => {
			const isAlive = gridState[x][y];

			if (isAlive) {
				cellDiv.style.backgroundColor = 'black';
			} else {
				cellDiv.style.backgroundColor = 'gray';
			}
		});
	});
}

function getLiveNeighborCount(x, y) {
	const leftX = (x + 99) % 100;
	const rightX = (x + 1) % 100;
	const aboveY = (y + 99) % 100;
	const belowY = (y + 1) % 100;

	return [
		gridState[leftX][aboveY],
		gridState[leftX][y],
		gridState[leftX][belowY],
		gridState[x][aboveY],
		gridState[x][belowY],
		gridState[rightX][aboveY],
		gridState[rightX][y],
		gridState[rightX][belowY],
	].map(b => b ? 1 : 0)
		.reduce((a, b) => a + b, 0);
}

function calculateNextGridState() {
	const newGridState = Array(100);
	for (let x = 0; x < 100; x += 1) {
		newGridState[x] = Array(100);
		for (let y = 0; y < 100; y += 1) {
			let isAlive = gridState[x][y];
			const count = getLiveNeighborCount(x, y);

			isAlive = isAlive
				? (count == 2 || count == 3)
				: (count == 3);

			newGridState[x][y] = isAlive;
		}
	}
	return newGridState;
}

function tick() {
	gridState = calculateNextGridState();
	updateGridElementsColor();
}

initGrid();
updateGridElementsColor();

setInterval(tick, 1000);

type PlayerState = {
	pos: number;
	energy: number;
	action: string;
	advanced: boolean;
	player?: number;
};

let turn = $state<number>(0);
let maxTurns = $state<number>(20);
let lastAction = $state<string>('');
let gameOver = $state<boolean>(false);
let winner = $state<string>('');
let waitingForOpponent = $state<boolean>(false);
let opponent = $state<PlayerState>({
	pos: 4, energy: 10, action: '', advanced: false
});
let you = $state<PlayerState>({
	pos: 2, energy: 10, action: '', advanced: false
});

export function gameState() {
	return {
		get turn() { return turn; },
		set turn(value) { turn = value; },
		get maxTurns() { return maxTurns; },
		set maxTurns(value) { maxTurns = value; },
		get lastAction() { return lastAction; },
		set lastAction(value) { lastAction = value; },
		get gameOver() { return gameOver; },
		set gameOver(value) { gameOver = value; },
		get winner() { return winner; },
		set winner(value) { winner = value; },
		get waitingForOpponent() { return waitingForOpponent; },
		set waitingForOpponent(value) { waitingForOpponent = value; },
		get opponent() { return opponent; },
		set opponent(value) { opponent = value; },
		get you() { return you; },
		set you(value) { you = value; }
	};
}

let rooms = $state<string[]>([]);
let userState = $state('');
let currentRoom = $state('');
let error = $state<string | undefined>('');

export function useState() {
	return {
		get userState() { return userState; },
		set userState(value) { userState = value; },
		get rooms() { return rooms; },
		set rooms(value) { rooms = value; },
		get currentRoom() { return currentRoom; },
		set currentRoom(value) { currentRoom = value; },
		get error() { return error; },
		set error(value) { error = value; },
	};
}
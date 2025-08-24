// Update the import path to the correct store file and exported member
import { useState } from './stores/state.svelte';
import { gameMessageHandler } from './stores/gameMessageHandler';
import { lobbyMessageHandler } from './stores/lobbyMessageHandler';
import { PLAYER_ID } from './constants/player';

const states = useState();
const decoder = new TextDecoder("utf-8");
let ws: WebSocket;
/**
 * Create a websocket connection
 * @param {string} socketURL
 * @param params
 * @returns
 */
export const connect = (socketURL: string) => {
	ws = new WebSocket(socketURL);

	if (!ws) {
		// Store an error in our state.  The function will be
		// called with the current state;  this only adds the
		// error.
		states.error = 'Unable to connect';
		return;
	}
	//ws.binaryType = 'arraybuffer';

	ws.addEventListener('open', () => {
		states.userState = 'connected';
		ws.send(JSON.stringify({ type: "init_client", payload: { playerId: PLAYER_ID } }));
		send("join_room", {
			"room": "default",
			"playerId": PLAYER_ID,
		});
		ws.send(JSON.stringify({ type: "list_rooms" }));
	});

	ws.addEventListener('message', ({ data }) => {
		//const msg = data instanceof ArrayBuffer ? JSON.parse(decoder.decode(data)) : JSON.parse(data);
		const msg = JSON.parse(data);
		if (states.currentRoom) {
			gameMessageHandler(msg);
		} else {
			lobbyMessageHandler(msg);
		}
	});

	ws.addEventListener('close', (message) => {
		console.log('Disconnected:', message);
		//state.update((state) => ({ ...state, error: message }));
	});

	ws.addEventListener('error', (err) => {
		console.log('websocket error:', err);
		states.error = 'Server encountered an error';
	});
};

export const send = (type: string, message: Record<string, unknown>): void => {
	if (!ws || ws.readyState !== WebSocket.OPEN) {
		console.error('WebSocket is not open.');
		return;
	}
	message = Object.assign({ type }, { payload: message });

	ws.send(JSON.stringify(message));
};

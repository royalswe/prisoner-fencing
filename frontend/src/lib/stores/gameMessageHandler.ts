import { LOBBY_EVENT as EVENT } from '../constants/events';
import { gameState } from './gameState.svelte';

const gs = gameState();


export function gameMessageHandler(msg: any) {
    console.log('Received game message:', msg);
    const payload = msg.payload;

    switch (msg.type) {
        case EVENT.error:
            console.log('lobbyerror', msg);
            break;
        case 'GAME_ACTION_RESULT':
            gs.turn = payload.turn;
            gs.maxTurns = payload.maxTurns;
            gs.lastAction = payload.lastAction;
            gs.gameOver = payload.gameOver;
            gs.winner = payload.winner;
            gs.status = payload.status;
            gs.opponent = payload.playerStates.opponent || {};
            gs.you = payload.playerStates.you || {};
            break;
        case 'UPDATE_STATUS':
            gs.status = payload.status;
            break;
        default:
            console.log('unknown emit from server', msg);
            break;
    }
}
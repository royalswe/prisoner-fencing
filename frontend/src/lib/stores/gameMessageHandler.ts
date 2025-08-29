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
            if (payload.turn !== undefined) gs.turn = payload.turn;
            if (payload.maxTurns !== undefined) gs.maxTurns = payload.maxTurns;
            if (payload.lastAction !== undefined) gs.lastAction = payload.lastAction;
            if (payload.gameOver !== undefined) gs.gameOver = payload.gameOver;
            if (payload.winner !== undefined) gs.winner = payload.winner;
            if (payload.status !== undefined) gs.status = payload.status;
            if (payload.playerStates?.opponent !== undefined) gs.opponent = payload.playerStates.opponent;
            if (payload.playerStates?.you !== undefined) gs.you = payload.playerStates.you;
            break;
        case 'UPDATE_STATUS':
            gs.status = payload.status;
            break;
        default:
            console.log('unknown emit from server', msg);
            break;
    }
}
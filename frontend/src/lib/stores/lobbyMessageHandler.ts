import { LOBBY_EVENT as EVENT } from '../constants/events';
import { useState } from '../stores/state.svelte';

const states = useState();


export function lobbyMessageHandler(msg: any) {
    switch (msg.type) {
        case EVENT.listRooms:
            states.rooms = msg.payload.rooms || [];
            break;
        case EVENT.joinRoom:
            states.currentRoom = msg.payload.room;
            break;
        case EVENT.chat:
            console.log(`Message from room ${msg.room}: ${msg.payload}`);
            //states.chat.push({ sender: msg.clientId, message: msg.msg, datetime: msg.datetime });
            break;
        case EVENT.initClient:
            console.log('Client initialized:', msg.payload);
            break;
        case EVENT.error:
            console.log('lobbyerror', msg);
            break;
        default:
            console.log('unknown emit from server', msg);
            break;
    }
}
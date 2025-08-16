// let _playerId: string | undefined = undefined;
// export function setPlayerId(id: string) {
//     if (_playerId) throw new Error("Player ID already set");
//     _playerId = id;
// }
// export function getPlayerId() {
//     return _playerId;
// }

export const PLAYER_ID = generatePlayerId();

// function that generates a random player ID
function generatePlayerId() {
    return Math.random().toString(36).substring(2, 15);
}
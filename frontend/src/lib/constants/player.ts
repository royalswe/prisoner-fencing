let _playerId: string | undefined = undefined;

function generatePlayerId() {
    return Math.random().toString(36).substring(2, 15);
}

export const PLAYER_ID = (() => {
    if (!_playerId) {
        _playerId = generatePlayerId();
    }
    return _playerId;
})();
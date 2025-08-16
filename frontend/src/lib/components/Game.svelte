<script lang="ts">
  import { send } from "../ws";
  import { gameState } from "../stores/gameState.svelte";
  import { PLAYER_ID } from "../constants/player";

  // Game state from store
  const gs = gameState();

  const { room } = $props<{ room: string }>();

  // Board icons for actions
  const actionIcons: Record<string, string> = {
    WAIT: "⏸️",
    RETREAT: "⬅️",
    ADVANCE: "➡️",
    ATTACK: "⚔️",
    COUNTER: "🛡️",
    DEFAULT: "👤",
    PLAYER: "🟥",
    OPPONENT: "🟦",
  };

  // Actions
  const actions = [
    { key: "1", name: "WAIT", desc: "+1 Energy" },
    { key: "2", name: "RETREAT", desc: "Move away, -1 Energy" },
    { key: "3", name: "ADVANCE", desc: "Move toward, double attack next turn" },
    { key: "4", name: "ATTACK", desc: "Attack, -1 Energy on miss" },
    {
      key: "5",
      name: "COUNTER",
      desc: "Counter, reflect if attacked, else -2 Energy",
    },
  ];

  // Send action to server
  function handleAction(key: string) {
    console.log(
      `Handling action: ${key} for player ${PLAYER_ID} in room ${room}`
    );

    if (gs.gameOver) return;
    send("game_action", { room, playerId: PLAYER_ID, action: key });
  }

  // Board rendering helper
  function getBoardActions() {
    const arr = Array(7).fill(null);
    // Only invert positions for player 2
    if (gs.you.player === 2) {
      arr[6 - gs.you.pos] = "PLAYER";
      arr[6 - gs.opponent.pos] = "OPPONENT";
    } else {
      arr[gs.you.pos] = "PLAYER";
      arr[gs.opponent.pos] = "OPPONENT";
    }
    return arr;
  }
</script>

<main>
  last actions:{gs.lastAction}
  <div
    style="display: flex; justify-content: center; margin: 1em 0; gap: 0.5em;"
  >
    {#each getBoardActions() as val, i}
      <div
        style="width: 48px; height: 48px; display: flex; align-items: center; justify-content: center; font-size: 2em; border: 2px solid #ccc; border-radius: 8px;position: relative;
        background: {val === 'PLAYER'
          ? 'red'
          : val === 'OPPONENT'
            ? 'blue'
            : '#f8f8f8'};"
      >
        {#if val === "PLAYER"}
          <span style="border-color: red;"
            >{actionIcons[gs.you.action] || actionIcons.DEFAULT}</span
          >
        {:else if val === "OPPONENT"}
          {#if gs.opponent.action == "RETREAT"}
            <span style="border-color: blue;"
              >{actionIcons["ADVANCE"] || actionIcons.DEFAULT}</span
            >
          {:else if gs.opponent.action == "ADVANCE"}
            <span style="border-color: blue;"
              >{actionIcons["RETREAT"] || actionIcons.DEFAULT}</span
            >
          {:else}
            <span style="border-color: blue;"
              >{actionIcons[gs.opponent.action] || actionIcons.DEFAULT}</span
            >
          {/if}
        {:else}
          <span>{actionIcons.DEFAULT}</span>
        {/if}
      </div>
    {/each}
  </div>
  <h2>Prisoner's Fencing - Room: {room}</h2>
  <div>Turn: {gs.turn} / {gs.maxTurns}</div>
  <div style="display: flex; gap: 2em; margin: 1em 0;">
    <div>
      <strong>You</strong><br />
      Energy: {gs.you.energy}<br />
      Position: {gs.you.player === 2 ? 6 - gs.you.pos : gs.you.pos}
    </div>
    <div>
      <strong>Opponent</strong><br />
      Energy: {gs.opponent.energy}<br />
      Position: {gs.you.player === 2 ? 6 - gs.opponent.pos : gs.opponent.pos}
    </div>
  </div>
  <div style="margin-bottom: 1em;">
    <strong>Actions:</strong>
    <div style="display: flex; gap: 1em; flex-wrap: wrap;">
      {#each actions as act}
        <button disabled={gs.gameOver} onclick={() => handleAction(act.name)}>
          {actionIcons[act.name] || actionIcons.DEFAULT}
          {act.name} ({act.key})<br /><small>{act.desc}</small>
        </button>
      {/each}
    </div>
  </div>
  <div style="margin: 1em 0; color: #444;">
    {#if gs.lastAction}
      {gs.lastAction}
    {:else}
      Waiting for opponent...
    {/if}
  </div>
  {#if gs.gameOver}
    <div style="font-size: 1.5em; color: darkred;">{gs.winner}</div>
  {/if}
</main>

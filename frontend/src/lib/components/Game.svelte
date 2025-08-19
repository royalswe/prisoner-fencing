<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { send } from "../ws";
  import { gameState } from "../stores/gameState.svelte";
  import { PLAYER_ID } from "../constants/player";

  // Game state from store
  const gs = gameState();

  const { room } = $props<{ room: string }>();

  // Board icons for actions
  const actionIcons: Record<string, string> = {
    WAIT: "/img/wait.png",
    RETREAT: "/img/retreat.png",
    ADVANCE: "/img/advance.png",
    ATTACK: "/img/knife.png",
    COUNTER: "/img/counter.png",
    DEFAULT: "/img/default.png",
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

  onMount(() => {
    window.addEventListener("keydown", handleKeydown);
  });
  onDestroy(() => {
    window.removeEventListener("keydown", handleKeydown);
  });

  function handleKeydown(event: KeyboardEvent) {
    const action = actions.find((a) => a.key === event.key)?.name;
    if (action) {
      handleAction(action);
      event.preventDefault();
    }
  }

  // Send action to server
  function handleAction(key: string) {
    if (gs.gameOver) return;
    const action = actions.find((a) => a.name === key);
    if (!action) return;
    // Remove highlight from all buttons
    const buttons = document.querySelectorAll("button[aria-label]");
    buttons.forEach((btn) => btn.classList.remove("js-highlight"));
    // Highlight selected action
    const button = document.querySelector(
      `button[aria-label="${action.name}"]`
    );

    if (button) {
      button.classList.add("js-highlight");
    }
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

<h2>Prisoner's Fencing - Room: {room}</h2>
{#if gs.gameOver}
  <p class="game-over">
    {gs.winner}
  </p>
{/if}
{gs.lastAction || "Choose an action!"}
<div class="board-row">
  {#each getBoardActions() as val}
    <div
      class="board-cell {val === 'PLAYER'
        ? 'player-cell'
        : val === 'OPPONENT'
          ? 'opponent-cell'
          : ''}"
    >
      {#if val === "PLAYER"}
        <span>
          <img
            src={actionIcons[gs.you.action] || actionIcons.DEFAULT}
            alt={gs.you.action}
            class="action-img"
          />
        </span>
      {:else if val === "OPPONENT"}
        <span>
          <img
            src={actionIcons[gs.opponent.action] || actionIcons.DEFAULT}
            alt={gs.opponent.action}
            class="action-img reverse-img"
          />
        </span>
      {:else}
        <span></span>
      {/if}
    </div>
  {/each}
</div>
<div>Turn: {gs.turn} / {gs.maxTurns}</div>
<div class="player-info-row">
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
<div class="actions-section">
  <strong>Actions:</strong>
  <div class="actions-row">
    {#each actions as act}
      <button
        aria-label={act.name}
        disabled={gs.gameOver}
        onclick={() => handleAction(act.name)}
      >
        <img
          src={actionIcons[act.name] || actionIcons.DEFAULT}
          alt={act.name}
          class="action-btn-img"
        />
        {act.name} ({act.key})<br /><small>{act.desc}</small>
      </button>
    {/each}
  </div>
</div>

<style>
  .board-row {
    display: flex;
    justify-content: center;
    margin: 1em 0;
    gap: 0.5em;
  }
  .board-cell {
    width: 76px;
    height: 76px;
    border: 2px solid #646cff;
    border-radius: 10px;
    position: relative;
    background: linear-gradient(135deg, #e0c3fc 0%, #8ec5fc 100%);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.07);
  }
  :root[data-theme="dark"] .board-cell {
    background: linear-gradient(135deg, #434343 0%, #262626 100%);
    border-color: #43cea2;
  }
  .player-cell {
    background: linear-gradient(135deg, #ff5858 0%, #f09819 100%) !important;
    border: 2px solid #d7263d;
  }
  .opponent-cell {
    background: linear-gradient(135deg, #43cea2 0%, #185a9d 100%) !important;
    border: 2px solid #185a9d;
  }
  .action-img {
    width: 100%;
    height: 100%;
  }
  .player-info-row {
    display: flex;
    gap: 2em;
    margin: 1em 0;
    justify-content: center;
  }
  .actions-section {
    margin-bottom: 1em;
    margin-top: 1em;
  }
  .actions-row {
    display: flex;
    gap: 1em;
    flex-wrap: wrap;
    justify-content: center;
  }
  .action-btn-img {
    width: 28px;
    height: 28px;
    vertical-align: middle;
    margin-bottom: 0.2em;
  }
  .reverse-img {
    -webkit-transform: scaleX(-1);
    transform: scaleX(-1);
  }
  button[aria-label] {
    background: linear-gradient(90deg, #f7971e 0%, #ffd200 100%);
    color: #222;
    border: none;
    border-radius: 10px;
    padding: 0.7em 1.2em;
    font-weight: bold;
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.07);
    transition: background 0.2s;
  }
  :root[data-theme="dark"] button[aria-label] {
    background: linear-gradient(90deg, #232526 0%, #414345 100%);
    color: #eee;
  }
  button[aria-label]:hover {
    filter: brightness(1.1);
  }

  .game-over {
    font-size: 1.5em;
    color: #23cb00;
    font-weight: 700;
    margin-bottom: 1em;
  }
</style>

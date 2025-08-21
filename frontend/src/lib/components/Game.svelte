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

  let youEnergyChange = $state(0);
  let opponentEnergyChange = $state(0);

  // keep last value in a normal variable
  let youLastEnergy = gs.you.energy;
  let opponentLastEnergy = gs.opponent.energy;

  // run every time `energy` changes
  $effect(() => {
    youEnergyChange = gs.you.energy - youLastEnergy;
    youLastEnergy = gs.you.energy;
    opponentEnergyChange = gs.opponent.energy - opponentLastEnergy;
    opponentLastEnergy = gs.opponent.energy;
  });
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
    <div class="info-table">
      <div class="info-row">
        <span class="info-label">Energy:</span>
        <span class="info-value">
          {#if youEnergyChange !== 0}
            <span
              style="color: {youEnergyChange > 0
                ? 'green'
                : 'red'}; font-weight: bold; margin-left: 0.5em;"
            >
              ({youEnergyChange > 0 ? "+" : ""}{youEnergyChange})
            </span>
          {/if}
          {gs.you.energy}
        </span>
      </div>
      <div class="info-row">
        <span class="info-label">Position:</span>
        <span class="info-value"
          >{gs.you.player === 2 ? 6 - gs.you.pos : gs.you.pos}</span
        >
      </div>
      <div class="info-row">
        <span class="info-label">Advanced:</span>
        <span class="info-value">{gs.you.advanced ? "True" : "False"}</span>
      </div>
    </div>
  </div>
  <div>
    <strong>Opponent</strong><br />
    <div class="info-table">
      <div class="info-row">
        <span class="info-label">Energy:</span>
        <span class="info-value">
          {#if opponentEnergyChange !== 0}
            <span
              style="color: {opponentEnergyChange > 0
                ? 'green'
                : 'red'}; font-weight: bold; margin-left: 0.5em;"
            >
              ({opponentEnergyChange > 0 ? "+" : ""}{opponentEnergyChange})
            </span>
          {/if}
          {gs.opponent.energy}
        </span>
      </div>
      <div class="info-row">
        <span class="info-label">Position:</span>
        <span class="info-value"
          >{gs.you.player === 2 ? 6 - gs.opponent.pos : gs.opponent.pos}</span
        >
      </div>
      <div class="info-row">
        <span class="info-label">Advanced:</span>
        <span class="info-value">{gs.opponent.advanced ? "True" : "False"}</span
        >
      </div>
    </div>
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
  .info-table {
    display: flex;
    flex-direction: column;
    margin-top: 0.5em;
  }
  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .info-label {
    text-align: left;
    min-width: 110px;
    font-weight: 600;
  }
  .info-value {
    text-align: right;
    font-weight: 600;
  }
</style>

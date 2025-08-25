<script lang="ts">
  export let open = false;
  export let onClose: () => void;
</script>

{#if open}
  <button
    class="modal-backdrop"
    aria-label="Close How to Play"
    onclick={onClose}
    onkeydown={(e) => (e.key === "Enter" || e.key === " ") && onClose()}
  ></button>
  <div class="modal">
    <button class="close-btn" onclick={onClose}>&times;</button>
    <h2>How to Play</h2>
    <div class="rules">
      <p>
        <strong>Prisoner's Fencing</strong> is a 2-player turn-based strategy
        game. The objective is to outlast your opponent by having more energy or
        <span style="color:#d32f2f">KILLING</span> them within a set number of turns.
      </p>
      <h3>🕹️ Actions Available to Players</h3>
      <ul>
        <li>
          <strong>WAIT (1-Key):</strong> Stay in current position.
          <em>+1 Energy</em>
        </li>
        <li>
          <strong>RETREAT (2-Key):</strong> Move one square away from opponent.
          <em>-1 Energy</em>
        </li>
        <li>
          <strong>ADVANCE (3-Key):</strong> Move one square toward opponent.
          <em>Double attack damage next turn only.</em><br />
          If both players advance to the same open square, a coin flip decides who
          wins the position. <em>-1 Energy</em>
        </li>
        <li>
          <strong>ATTACK (4-Key):</strong> Attempt to attack opponent.
          <em>-1 Energy on miss</em><br />
          <strong>Damage:</strong> Normal = 3, Advanced = 6<br />
        </li>
        <li>
          <strong>COUNTER (5-Key):</strong> Attempt to counter an incoming
          attack.<br />
          <strong>If Opponent is Attacking:</strong> Reflects damage back at no
          energy cost.<br />
          <strong>If Opponent is Not Attacking:</strong> Lose 2 Energy.
        </li>
      </ul>
      <h3>🏁 Winning Conditions</h3>
      <ul>
        <li>
          <strong>Primary:</strong> Eliminate opponent by reducing their energy to
          0.
        </li>
        <li>
          <strong>Secondary:</strong> Have more energy when the total number of turns
          is reached.
        </li>
        <li>
          <strong>Draw:</strong> Both players have equal energy at the end of the
          final turn.
        </li>
      </ul>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.3);
    z-index: 1000;
  }
  .modal {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: #f0f4ff;
    border-radius: 12px;
    padding: 2em 2.5em;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.15);
    z-index: 1001;
    min-width: 320px;
    max-width: 90vw;
  }
  .modal h2 {
    margin-top: 0;
    color: #2a3a5e;
    font-size: 1.5em;
    margin-bottom: 1em;
  }
  .rules {
    font-size: 1em;
    color: #222;
    margin-bottom: 0.5em;
  }
  .rules h3 {
    margin-top: 1.2em;
    margin-bottom: 0.5em;
    font-size: 1.1em;
    color: #2a3a5e;
  }
  .rules ul {
    margin: 0 0 0.7em 0;
    padding-left: 1.2em;
  }
  .rules li {
    margin-bottom: 0.5em;
  }
  .close-btn {
    position: absolute;
    top: 0.7em;
    right: 1em;
    background: none;
    border: none;
    font-size: 2em;
    color: #888;
    cursor: pointer;
    line-height: 1;
  }
  .close-btn:hover {
    color: #d32f2f;
  }
</style>

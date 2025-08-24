<script lang="ts">
  import { onMount } from "svelte";
  import { useState } from "../stores/state.svelte";
  import { PLAYER_ID } from "../constants/player";
  import { connect, send } from "../ws";

  const states = useState();
  let newRoomName = "";
  let newMessage = "";

  onMount(() => {
    connect(
      window.location.protocol !== "https:"
        ? `ws://${window.location.hostname}:8080/ws`
        : "wss://prisonerfencing-server.fumlig.com/ws"
    );
  });

  function createRoom() {
    if (newRoomName.trim()) {
      send("join_room", {
        room: newRoomName.trim(),
        PLAYER_ID,
      });

      newRoomName = "";
    }
  }

  function joinRoom(room: string) {
    send("join_room", {
      room,
      PLAYER_ID,
    });
  }

  function sendMessage(message: string) {
    send("send_message", {
      message: message.trim(),
      from: PLAYER_ID,
    });
  }
</script>

<main>
  <div class="lobby-grid-container">
    <header class="lobby-header">
      <h1>Prisoner Fencing</h1>
      <h2>Lobby</h2>
      <div class="your-id">
        <strong>Your ID:</strong>
        {PLAYER_ID}
      </div>
    </header>
    <main class="lobby-main-area">
      {#if !states.userState}
        <div>Connecting to lobby...</div>
      {:else}
        <h3 class="mb-0">Available Rooms</h3>
        {#if states.rooms.length === 0}
          <div>No rooms yet.</div>
        {:else}
          <ul class="room-list">
            {#each states.rooms as room}
              <li>
                <span>{room}</span>
                <button onclick={() => joinRoom(room)}>Join</button>
              </li>
            {/each}
          </ul>
        {/if}
        <div class="input-row">
          <input
            placeholder="New room name"
            bind:value={newRoomName}
            onkeydown={(e) => e.key === "Enter" && createRoom()}
          />
          <button onclick={createRoom}>Create Room</button>
        </div>
        {#if states.error}
          <div class="error">{states.error}</div>
        {/if}
      {/if}
    </main>
  </div>
</main>

<style>
  .lobby-grid-container {
    display: grid;
    grid-template-areas:
      "header"
      "main";
    justify-content: center;
    width: 100vw;
    box-sizing: border-box;
    overflow: hidden;
  }
  .lobby-header {
    grid-area: header;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 1em 0 0.5em 0;
  }
  .lobby-header h1 {
    letter-spacing: 2px;
    font-size: 2.5em;
    margin-bottom: 1em;
  }
  .lobby-header h2 {
    margin-top: 0;
    font-size: 1.5em;
    margin-bottom: 0.5em;
  }
  .your-id {
    font-size: 1.1em;
    background: #fff7e6;
    color: #222;
    border-radius: 8px;
    padding: 0.3em 1em;
    margin-bottom: 1rem;
    display: inline-block;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.07);
  }
  .lobby-main-area {
    grid-area: main;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    width: 100%;
    max-width: 500px;
    margin: 0 auto;
    padding: 1em 0;
  }
  .room-list {
    list-style: none;
    padding: 0 1rem;
    width: 100%;
  }
  .room-list li {
    background: #fff7e6;
    color: #222;
    border-radius: 8px;
    margin-bottom: 0.5rem;
    padding: 0.5rem 1rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 1.1em;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.07);
  }
  .room-list li button {
    background: linear-gradient(90deg, #ff6a00 0%, #ee0979 100%);
    color: #fff;
    border: none;
    border-radius: 6px;
    padding: 0.3em 1em;
    font-weight: bold;
    cursor: pointer;
  }
  .room-list li button:hover {
    filter: brightness(1.2);
  }

  :root[data-theme="dark"] .input-row button {
    background: linear-gradient(90deg, #fc00ff 0%, #00dbde 100%);
    color: #fff;
  }

  :root[data-theme="dark"] .room-list li button {
    background: linear-gradient(90deg, #43cea2 0%, #185a9d 100%);
  }
  .input-row {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
    padding: 1rem 0;
    width: 100%;
  }
  .input-row input {
    flex: 1;
    padding: 0.5em 1em;
    border-radius: 8px;
    border: 1px solid #ccc;
    font-size: 1em;
    background: #fff;
    color: #222;
  }
  .input-row input:focus {
    border-color: #646cff;
    outline: none;
  }
  .input-row button {
    background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
    color: #222;
    border: none;
    border-radius: 8px;
    padding: 0.5em 1.2em;
    font-weight: bold;
    cursor: pointer;
    transition: background 0.2s;
  }
  .input-row button:hover {
    filter: brightness(1.2);
  }
  .error {
    color: #d32f2f;
    margin-top: 1rem;
    font-weight: bold;
  }
</style>

<script lang="ts">
  import { onMount } from "svelte";
  import HowToPlayModal from "./HowToPlayModal.svelte";
  import { useState } from "../stores/state.svelte";
  import { PLAYER_ID } from "../constants/player";
  import { connect, send } from "../ws";

  const states = useState();
  let showHowToPlay = false;
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
      <button class="how-to-play-btn" onclick={() => (showHowToPlay = true)}
        >How to Play</button
      >
      <HowToPlayModal
        open={showHowToPlay}
        onClose={() => (showHowToPlay = false)}
      />

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
  @import "../styles/lobby.css";
</style>

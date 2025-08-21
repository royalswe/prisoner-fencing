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
  <h1 style="letter-spacing: 2px; font-size: 2.5em; margin-bottom: 0.5em;">
    Prisoner Fencing
  </h1>
  <section class="lobby-card">
    <h2 style="margin-top: 0;">Lobby</h2>
    <div class="your-id">
      <strong>Your ID:</strong>
      {PLAYER_ID}
    </div>
    {#if !states.userState}
      <div>Connecting to lobby...</div>
    {:else}
      <div>
        <h3 style="margin-bottom: 0.5em;">Available Rooms</h3>
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
      </div>
      {#if states.error}
        <div class="error">{states.error}</div>
      {/if}
    {/if}
  </section>
</main>

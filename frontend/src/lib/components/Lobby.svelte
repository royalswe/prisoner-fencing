<script lang="ts">
  import { onMount } from "svelte";
  import Game from "./Game.svelte";
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
  <h1>Prisoner Fencing</h1>
  state: {states.userState}
  {#if states.currentRoom}
    <h2>Game in Progress</h2>
    <Game room={states.currentRoom} />
  {:else}
    <section>
      <h2>Lobby</h2>
      <div>
        <strong>Your ID:</strong>
        {PLAYER_ID}
      </div>
      {#if !states.userState}
        <div>Connecting to lobby...</div>
      {:else}
        <div>
          <h3>Available Rooms</h3>
          {#if states.rooms.length === 0}
            <div>No rooms yet.</div>
          {:else}
            <ul>
              {#each states.rooms as room}
                <li>
                  {room} <button onclick={() => joinRoom(room)}>Join</button>
                </li>
              {/each}
            </ul>
          {/if}
          <input
            placeholder="New room name"
            bind:value={newRoomName}
            onkeydown={(e) => e.key === "Enter" && createRoom()}
          />
          <button onclick={createRoom}>Create Room</button>

          <input
            placeholder="Message"
            bind:value={newMessage}
            onkeydown={(e) => e.key === "Enter" && sendMessage(newMessage)}
          />
          <button onclick={() => sendMessage(newMessage)}>Send Message</button>
        </div>
        {#if states.error}
          <div style="color: red">{states.error}</div>
        {/if}
      {/if}
    </section>
  {/if}
</main>

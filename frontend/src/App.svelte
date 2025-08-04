<script lang="ts">
  import svelteLogo from "./assets/svelte.svg";
  import viteLogo from "/vite.svg";
  import { onMount } from "svelte";

  // Generate a random player id
  function randomId() {
    return Math.random().toString(36).substring(2, 10);
  }
  const playerId = randomId();

  let ws: WebSocket | null = null;
  let connected = false;
  let rooms: string[] = [];
  let currentRoom = "";
  let newRoomName = "";
  let newMessage = "";
  let lobbyError = "";

  function connect() {
    let wsUrl = "";
    if (
      window.location.hostname === "localhost" ||
      window.location.hostname === "127.0.0.1"
    ) {
      wsUrl = `ws://${window.location.hostname}:8080/ws`;
    } else {
      wsUrl = "wss://prisonerfencing-server.fumlig.com/ws";
    }
    ws = new WebSocket(wsUrl);
    ws.onopen = () => {
      connected = true;
      console.log(`WebSocket connected as player ${playerId}`);

      // Request room list
      ws?.send(JSON.stringify({ type: "list_rooms", playerId }));
    };
    ws.onmessage = (event) => {
      console.log(`Received message: ${event.data}`);

      try {
        const msg = JSON.parse(event.data);
        if (msg.type === "room_list") {
          rooms = msg.rooms;
        } else if (msg.type === "new_message") {
          console.log(`Message from room ${msg.room}: ${msg.payload}`);
        } else if (msg.type === "joined_room") {
          currentRoom = msg.room;
        } else if (msg.type === "error") {
          lobbyError = msg.message;
        }
      } catch (e) {
        // ignore
      }
    };
    ws.onclose = () => {
      connected = false;
      currentRoom = "";
      setTimeout(connect, 2000); // try to reconnect
    };
  }

  onMount(() => {
    connect();
    return () => ws?.close();
  });

  function createRoom() {
    if (ws && newRoomName.trim()) {
      ws.send(
        JSON.stringify({
          type: "create_room",
          payload: {
            room: newRoomName.trim(),
            playerId,
          },
        })
      );
      newRoomName = "";
    }
  }

  function joinRoom(room: string) {
    if (ws) {
      ws.send(
        JSON.stringify({ type: "join_room", payload: { room, playerId } })
      );
    }
  }

  function sendMessage(message: string) {
    if (ws) {
      ws.send(
        JSON.stringify({
          type: "send_message",
          payload: {
            message: message.trim(),
            from: playerId,
          },
        })
      );
    }
  }
</script>

<main>
  <div>
    <a href="https://vite.dev" target="_blank" rel="noreferrer">
      <img src={viteLogo} class="logo" alt="Vite Logo" />
    </a>
    <a href="https://svelte.dev" target="_blank" rel="noreferrer">
      <img src={svelteLogo} class="logo svelte" alt="Svelte Logo" />
    </a>
  </div>
  <h1>Prisoner Fencing</h1>

  <section>
    <h2>Lobby</h2>
    <div>
      <strong>Your ID:</strong>
      {playerId}
    </div>
    {#if !connected}
      <div>Connecting to lobby...</div>
    {:else}
      {#if currentRoom}
        <div>
          <h3>In Room: {currentRoom}</h3>
          <button
            on:click={() => {
              currentRoom = "";
              ws?.send(JSON.stringify({ type: "leave_room", playerId }));
            }}>Leave Room</button
          >
        </div>
      {:else}
        <div>
          <h3>Available Rooms</h3>
          {#if rooms.length === 0}
            <div>No rooms yet.</div>
          {:else}
            <ul>
              {#each rooms as room}
                <li>
                  {room} <button on:click={() => joinRoom(room)}>Join</button>
                </li>
              {/each}
            </ul>
          {/if}
          <input
            placeholder="New room name"
            bind:value={newRoomName}
            on:keydown={(e) => e.key === "Enter" && createRoom()}
          />
          <button on:click={createRoom}>Create Room</button>

          <input
            placeholder="Message"
            bind:value={newMessage}
            on:keydown={(e) => e.key === "Enter" && sendMessage(newMessage)}
          />
          <button on:click={() => sendMessage(newMessage)}>Send Message</button>
        </div>
      {/if}
      {#if lobbyError}
        <div style="color: red">{lobbyError}</div>
      {/if}
    {/if}
  </section>
</main>

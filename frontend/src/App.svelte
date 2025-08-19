<script lang="ts">
  import ThemeSwitcher from "./lib/components/ThemeSwitcher.svelte";
  import Lobby from "./lib/components/Lobby.svelte";
  import Game from "./lib/components/Game.svelte";
  import { useState } from "./lib/stores/state.svelte";

  const states = useState();
</script>

<main>
  <ThemeSwitcher />
  {#if states.currentRoom}
    <Game room={states.currentRoom} />
  {:else}
    <Lobby />
  {/if}
</main>

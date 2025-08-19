<script lang="ts">
  const getTheme = () => {
    // Check localStorage first
    const stored = localStorage.getItem("theme");
    if (stored) return stored;

    // Fall back to system preference
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  };

  const applyTheme = (theme: string) => {
    document.documentElement.dataset.theme = theme;
  };

  const toggleTheme = () => {
    const newTheme = getTheme() === "dark" ? "light" : "dark";
    // update icon on button
    const button = document.querySelector(".theme-toggle button");
    if (button) {
      button.textContent = newTheme === "light" ? "🌙" : "☀️";
    }
    applyTheme(newTheme);
    localStorage.setItem("theme", newTheme);
  };

  // Apply theme immediately
  const theme = getTheme();
  applyTheme(theme);
</script>

<div class="theme-toggle">
  <button onclick={toggleTheme} aria-label="Toggle dark/light mode">
    {theme === "light" ? "🌙" : "☀️"}
  </button>
</div>

<style>
  .theme-toggle {
    position: absolute;
    top: 1rem;
    right: 1rem;
    z-index: 100;
  }
  button[aria-label] {
    padding: 0.5em 1em;
    border-radius: 1em;
    border: none;
    background: var(--toggle-bg);
    color: var(--toggle-fg);
    cursor: pointer;
    font-size: 1.2em;
  }
</style>

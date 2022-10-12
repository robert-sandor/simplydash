<script lang="ts">
    import {afterUpdate, onMount} from "svelte";
    import {Config, ThemeColors} from "./config";

    export let config: Config;
    let currentTheme: ThemeColors = undefined;

    let shortcuts = {
        84: toggleTheme,
    };

    function toggleTheme() {
        currentTheme = (currentTheme == undefined || currentTheme == config.settings.theme.dark) ?
            config.settings.theme.light : config.settings.theme.dark;
        setTheme(currentTheme);
    }

    function setTheme(theme: ThemeColors) {
        document.body.style.setProperty('--background', theme.background);
        document.body.style.setProperty('--element-background', theme.element_background);
        document.body.style.setProperty('--foreground', theme.foreground);
        document.body.style.setProperty('--foreground-secondary', theme.foreground_secondary);
        document.body.style.setProperty('--accent-color', theme.accent_color);
        document.body.style.setProperty('--success-color', theme.success_color);
        document.body.style.setProperty('--warning-color', theme.warning_color);
        document.body.style.setProperty('--error-color', theme.error_color);
    }

    function onKeyDown(e) {
        if (shortcuts[e.keyCode] !== undefined) {
            shortcuts[e.keyCode]();
        }
    }

    onMount(() => {
        currentTheme = config.settings.theme.dark;
    })

    afterUpdate(() => {
        console.log("after update header")
        currentTheme = config.settings.theme.dark;
        setTheme(currentTheme)
    })
</script>

<div class="header container">
    <div class="name unselectable">{config.settings.name}</div>
    <div class="button" on:click={toggleTheme}>
        <svg
            class="icon"
            id="icon"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
        >
            <path
                d="M12,18V6C15.31,6 18,8.69 18,12C18,15.31 15.31,18 12,18M20,15.31L23.31,12L20,8.69V4H15.31L12,0.69L8.69,4H4V8.69L0.69,12L4,15.31V20H8.69L12,23.31L15.31,20H20V15.31Z"
            />
        </svg>
    </div>
</div>

<svelte:window on:keydown={onKeyDown} />

<style>
    .header {
        height: var(--header-height);
        box-sizing: border-box;
    }

    .name {
        float: left;
        height: 100%;
        padding: 8px 8px 8px 16px;
        text-align: left;
        font-size: 1.5em;
    }

    .button {
        float: right;
        height: 100%;
        padding: 12px;
        fill: var(--foreground-secondary);
    }

    .button:hover {
        fill: var(--foreground);
        cursor: pointer;
    }

    .icon {
        height: 24px;
        width: 24px;
    }
</style>

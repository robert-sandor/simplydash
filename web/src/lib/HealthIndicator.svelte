<script lang="ts">
    import { onMount } from "svelte";
    import {ApiClient} from "./api_client";

    export let apiClient: ApiClient;
    export let url;
    let state = 'waiting';

    async function healthcheck() {
        try {
            const status = await apiClient.healthcheck(url);
            state = status >= 500 || status == 404 || status < 0 ? "error" : "ok";
        } catch (e) {
            state = "error";
        }
    }

    onMount(() => {
        healthcheck();
        const interval = setInterval(healthcheck, 60000);
        return () => clearInterval(interval);
    });
</script>

{#if state === "waiting"}
    <div class="waiting"></div>
{:else if state === "error"}
    <div class="error"></div>
{:else}
    <div class="ok"></div>
{/if}

<style>
    div {
        position: relative;
        width: var(--indicator-size);
        height: var(--indicator-size);
        border-radius: calc(var(--indicator-size) / 2);
        left: 100%;
        top: 0%;
        transform: translate(-100%, 0%);
    }

    .waiting {
        background-color: var(--warning-color);
    }

    .ok {
        background-color: var(--success-color);
    }

    .error {
        background-color: var(--error-color);
    }
</style>

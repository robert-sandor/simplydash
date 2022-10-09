<script lang="ts">
    import HealthIndicator from "./HealthIndicator.svelte";
    import {afterUpdate, onMount} from "svelte";
    import {ApiClient} from "./api_client";

    export let apiClient: ApiClient;
    export let enableHealthIndicator: boolean;

    export let name: string
    export let url: string
    export let icon: string
    export let description: string

    function updateIcon() {
        const oldIcon = icon;
        icon = (icon != undefined && icon.trim() !== "") ? icon :
            `https://cdn.jsdelivr.net/gh/walkxcode/dashboard-icons/png/${name.toLowerCase().replaceAll(" ", "-")}.png`;
        if (oldIcon !== icon) {
            showIcon = true;
        }
    }

    afterUpdate(() => {
        description = (description != undefined && description.trim() !== "") ? description : url;
        updateIcon();
    })

    let showIcon: boolean = true;
</script>

<a href={url}>
    <div class="element hover button">
        {#if showIcon}
            <div class="icon-container">
                <img
                    class="icon"
                    src="{apiClient.host}/api/icon?url={encodeURIComponent(icon)}"
                    alt={name}
                    on:error={() => (showIcon = false)}
                />
            </div>
        {/if}
        <div class="button-text" class:text-width-icon={showIcon}>
            <div class="item-name unselectable">{name}</div>
            <div class="item-description unselectable">{description}</div>
        </div>
        {#if enableHealthIndicator}
            <HealthIndicator apiClient="{apiClient}" url={url} />
        {/if}
    </div>
</a>

<style>
    .button {
        padding: 8px;
        border-radius: var(--border-radius);
        height: var(--button-height);
    }

    .icon-container {
        display: flex;
        align-items: center;
        justify-content: center;
        float: left;
        height: var(--icon-size);
        width: var(--icon-size);
    }

    .icon {
        max-width: var(--icon-size);
        max-height: var(--icon-size);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .text-width-icon {
        width: calc(100% - var(--icon-size));
    }

    .button-text {
        float: left;
    }

    .item-name {
        float: left;
        width: 100%;
        padding: 0px 1px 1px 8px;
        font-size: 1.25em;
    }

    .item-description {
        float: left;
        width: 100%;
        padding: 1px 1px 1px 8px;
        font-size: 0.75em;
        color: var(--foreground-secondary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
</style>

<script lang="ts">
    import ItemButton from "./ItemButton.svelte";
    import {Item} from "./category";
    import {ApiClient} from "./api_client";
    import {afterUpdate, onMount} from "svelte";

    export let apiClient: ApiClient;
    export let enableHealthIndicators: boolean;
    export let name: string;
    export let items: Item[];

    let show = true;

    function toggle() {
        show = !show;
    }

</script>

<div class="container">
    <div class="element unselectable hover header" on:click={ toggle }>{name}</div>
    {#if show && items !== undefined}
        <div class="group-container">
            {#each items as item}
                <div class="cell">
                    <ItemButton name="{item.name}" url="{item.url}" icon="{item.icon}" description="{item.description}"
                                enableHealthIndicator="{enableHealthIndicators}" apiClient="{apiClient}"/>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>

    .header {
        padding: 12px;
        border-radius: var(--border-radius);
        margin: 4px;
    }

    .group-container {
        display: inline-block;
        width: 100%;
        border-top: 0;
        border-bottom-left-radius: var(--border-radius);
        border-bottom-right-radius: var(--border-radius);
    }

    .cell {
        width: 100%;
        float: left;
        padding: 4px;
    }

    @media only screen and (min-width: 768px) {
        .cell {
            width: 50%;
        }
    }

    @media only screen and (min-width: 992px) {
        .cell {
            width: 33.33%;
        }
    }

    @media only screen and (min-width: 1200px) {
        .cell {
            width: 25%;
        }
    }
</style>

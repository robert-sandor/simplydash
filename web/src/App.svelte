<script lang="ts">
    import ItemGroup from "./lib/ItemGroup.svelte";
    import Header from "./lib/Header.svelte";
    import {onMount} from "svelte";
    import {ApiClient} from "./lib/api_client.js";
    import {Config} from "./lib/config.js";
    import {Category} from "./lib/category.js";

    const apiHost = "http://localhost:8080"
    const apiClient: ApiClient = new ApiClient(apiHost)

    let categories: Category[] = [];
    let config: Config = Config.default();

    const updateCategories = () => {
        apiClient.getCategories().then((r) => categories = r);
    };
    const updateConfig = () => {
        apiClient.getConfig().then((r) => (config = r));
        updateCategories();
    };

    onMount(() => {
        updateCategories();
        updateConfig();
        apiClient.websocket(
            updateConfig,
            updateCategories
        );
    });
</script>

<svelte:head>
    <title>{config.settings.name}</title>
</svelte:head>

<Header {config}/>
{#each categories as category}
    <ItemGroup apiClient="{apiClient}"
               enableHealthIndicators={config.settings.enable_health_indicators}
               name={category.name}
               items={category.items}/>
{/each}

<style>
</style>

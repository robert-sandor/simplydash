<script lang="ts">
	import { AppGroup, AppSettings } from '$lib/models';
	import AppButton from '$lib/AppButton.svelte';
	import ThemeSwitcher from '$lib/ThemeSwitcher.svelte';
	import Websocket from '$lib/Websocket.svelte';
	import SearchButton from '$lib/SearchButton.svelte';
	import { onMount } from 'svelte';
	import { baseUrl } from '$lib/utils';

	let appGroups: AppGroup[] = [];
	let appSettings = new AppSettings();

	async function getSettings() {
		const response = await fetch(`${baseUrl()}/settings`);
		if (response.status != 200) {
			return;
		}
		appSettings = await response.json();
	}

	onMount(() => getSettings());
</script>

<header class="container flex flex-wrap items-center justify-between mx-auto px-4">
	<Websocket bind:appGroups />
	<h1 class="text-2xl text-neutral-800 dark:text-neutral-200 m-4">{appSettings.name}</h1>
	<div class="flex">
		<SearchButton bind:appGroups />
		<ThemeSwitcher />
	</div>
</header>
<div class="container mx-auto px-4">
	{#each appGroups as appGroup}
		{#if appGroup.apps.length > 0}
			<div class="m-2 text-neutral-800 dark:text-neutral-200">{appGroup.name}</div>
			<div class="grid md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 grid-flow-row auto-cols gap-4">
				{#each appGroup.apps as app}
					<AppButton {app} />
				{/each}
			</div>
		{/if}
	{/each}
</div>
<footer class="container flex flex-wrap items-center justify-end mx-auto p-4">
	<div class="text-sm text-neutral-500">simplydash version: dev</div>
</footer>

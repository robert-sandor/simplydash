<script lang="ts">
	import Fuse from 'fuse.js';
	import { blurOnEscape } from './directives';
	import type { App } from './models';
	import AppButton from './AppButton.svelte';
	import { onDestroy, onMount } from 'svelte';

	export let apps: App[] = [];

	let searchInput: HTMLInputElement;
	let searchString: string;

	let fuse: Fuse<App> = new Fuse(apps, { keys: ['name'] });
	let fuseResults: Fuse.FuseResult<App>[] = [];

	$: fuseResults = fuse.search(searchString ? searchString : '');

	onMount(() => {
		searchInput.focus();
		window.document.addEventListener('keyup', handleKeys);
	});

	onDestroy(() => {
		window.document.removeEventListener('keyup', handleKeys);
	});

	let selectedIndex = -1;
	function handleKeys(event: KeyboardEvent) {
		if (event.key === 'ArrowUp') {
			selectedIndex = selectedIndex <= 0 ? fuseResults.length - 1 : selectedIndex - 1;
			return;
		}
		if (event.key === 'ArrowDown') {
			selectedIndex =
				selectedIndex == fuseResults.length - 1 || selectedIndex < 0 ? 0 : selectedIndex + 1;
			return;
		}
		if (event.key === 'Enter') {
			window.location.assign(fuseResults[selectedIndex < 0 ? 0 : selectedIndex].item.link);
			return;
		}
	}
</script>

<div
	class="bg-neutral-200 dark:bg-neutral-800 w-5/6 md:w-[32rem] mx-auto drop-shadow-xl border-neutral-500 border-2 rounded-lg top-8 absolute left-1/2 -translate-x-1/2"
>
	<input
		bind:this={searchInput}
		bind:value={searchString}
		use:blurOnEscape
		autocomplete="off"
		placeholder=""
		class="w-full h-12 bg-transparent focus:outline-none text-neutral-800 dark:text-neutral-200 text-center my-2"
	/>
	{#if fuseResults.length > 0}
		<div class="grid grid-flow-row auto-cols gap-4 m-4">
			{#each fuseResults as fr, i}
				<AppButton app={fr.item} showSelected={i === selectedIndex} />
			{/each}
		</div>
	{/if}
</div>

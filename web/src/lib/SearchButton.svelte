<script lang="ts">
	import Search from './Search.svelte';
	import type { App, AppGroup } from './models';

	export let appGroups: AppGroup[] = [];

	let apps: App[] = [];
	$: apps = appGroups.flatMap((ag) => ag.apps);

	let showSearch = false;
	function show() {
		showSearch = true;
	}

	function hide() {
		showSearch = false;
	}

	function toggle() {
		showSearch ? hide() : show();
	}

	window.document.addEventListener('keyup', handleToggleKey);
	window.document.addEventListener('keyup', handleCloseKey);

	function handleToggleKey(event: KeyboardEvent) {
		if (event.ctrlKey && event.key === 'k') {
			toggle();
		}
	}

	function handleCloseKey(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			hide();
		}
	}
</script>

<button class="focus:outline-neutral-500" on:click={show}>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		viewBox="0 0 24 24"
		fill="currentColor"
		class="w-6 h-6 m-4 fill-neutral-800 dark:fill-neutral-200"
	>
		<path
			fill-rule="evenodd"
			d="M10.5 3.75a6.75 6.75 0 100 13.5 6.75 6.75 0 000-13.5zM2.25 10.5a8.25 8.25 0 1114.59 5.28l4.69 4.69a.75.75 0 11-1.06 1.06l-4.69-4.69A8.25 8.25 0 012.25 10.5z"
			clip-rule="evenodd"
		/>
	</svg>
</button>
{#if showSearch}
	<div
		class="absolute bg-neutral-200/30 dark:bg-neutral-800/30 w-full h-full top-0 left-0 z-10 backdrop-blur-sm"
		role="none"
		on:click={hide}
		on:keyup={handleCloseKey}
	>
		<Search {apps} />
	</div>
{/if}

<script lang="ts">
	import type { App } from './models';
	import { baseUrl } from './utils';

	export let app: App;
	export let showSelected = false;
</script>

<a href={app.link} class="focus:outline-none focus:ring outline-sky-500 rounded-lg ring-sky-500" class:ring={showSelected}>
	<div
		class="rounded-lg shadow-inner flex justify-end"
		class:bg-neutral-300={app.healthcheck.health === 'unknown'}
		class:dark:bg-neutral-600={app.healthcheck.health === 'unknown'}
		class:bg-emerald-400={app.healthcheck.health === 'healthy'}
		class:dark:bg-emerald-600={app.healthcheck.health === 'healthy'}
		class:bg-amber-400={app.healthcheck.health === 'timeout'}
		class:dark:bg-amber-600={app.healthcheck.health === 'timeout'}
		class:bg-orange-400={app.healthcheck.health === 'warning'}
		class:dark:bg-orange-600={app.healthcheck.health === 'warning'}
		class:bg-rose-400={app.healthcheck.health === 'error'}
		class:dark:bg-rose-600={app.healthcheck.health === 'error'}
	>
		<div
			class="bg-neutral-200 dark:bg-neutral-800 rounded-md flex items-center justify-between drop-shadow-lg m-1"
			style="width: 95%"
		>
			<div class="w-52 m-4">
				<div class="text-lg text-neutral-800 dark:text-neutral-200">{app.name}</div>
				<div
					class="text-xs text-neutral-600 dark:text-neutral-400 whitespace-nowrap overflow-hidden text-ellipsis"
				>
					{app.description}
				</div>
			</div>
			<img
				class="object-contain w-12 h-12 bg-none m-4"
				src="{baseUrl()}/image?url={app.icon}"
				alt={app.name}
			/>
		</div>
	</div>
</a>

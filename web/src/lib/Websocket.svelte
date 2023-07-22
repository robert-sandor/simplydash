<script lang="ts">
	import type { AppGroup } from './models';
	import { websocketUrl } from '$lib/utils';

	const reconnectAfterSeconds = 5;
	const wsUrl = websocketUrl();

	export let appGroups: AppGroup[] = [];

	let disconnected = false;
	let reconnected = false;
	let timeLeftInSeconds = 0;
	let updateTimeLeftHandle = 0;

	function newWebsocketConnection() {
		let websocket = new WebSocket(wsUrl);
		websocket.onclose = onClose;
		websocket.onmessage = onMessage;

		clearInterval(updateTimeLeftHandle);
	}

	function onMessage(event: MessageEvent<string>) {
		if (disconnected) {
			disconnected = false;
			reconnected = true;
			setTimeout(() => (reconnected = false), 1000);
		}

		appGroups = JSON.parse(event.data);
	}

	function onClose() {
		disconnected = true;
		setTimeout(newWebsocketConnection, reconnectAfterSeconds * 1000);

		timeLeftInSeconds = reconnectAfterSeconds;
		updateTimeLeftHandle = setInterval(() => (timeLeftInSeconds -= 1), 1000);
	}

	newWebsocketConnection();
</script>

<div
	class="absolute w-11/12 sm:w-auto text-center p-4 rounded-lg top-4 left-1/2 -translate-x-1/2 text-neutral-800 dark:text-neutral-200 shadow-xl transition-all ease-linear"
	class:bg-rose-400={disconnected}
	class:dark:bg-rose-600={disconnected}
	class:bg-emerald-400={reconnected}
	class:dark:bg-emerald-600={reconnected}
	class:hidden={!disconnected && !reconnected}
>
	{#if disconnected}
		connection lost - retrying in {timeLeftInSeconds}s
	{/if}
	{#if reconnected}
		reconnected
	{/if}
</div>

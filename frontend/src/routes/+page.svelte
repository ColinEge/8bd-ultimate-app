<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { setKeySelectionState } from "$lib/state/selection-state.svelte";
	import { setKeyboardState } from "$lib/state/keyboard-state.svelte";
	import {
		getNavigationState,
		setNavigationState,
		type NavigationEntry,
	} from "$lib/state/navigation-state.svelte";
	import { onMount } from "svelte";
	import Devices from "$lib/components/el/Devices.svelte";

	setKeySelectionState();
	setKeyboardState();

	setNavigationState();

	const nav = getNavigationState();
	const current: NavigationEntry = $derived.by(() => nav.peek());

	function goBack() {
		nav.pop();
	}

	onMount(() => {
		nav.push({
			component: Devices,
		});
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<header>
	{#if nav.length() > 1}
		<button class="button button-secondary ml-5 mt-5" onclick={goBack}> Back </button>
	{/if}
</header>

<div class="flex justify-center min-h-dvh w-full overflow-x-auto {current?.component === Devices ? 'items-center' : 'items-start'}">
	<div class="flex flex-col gap-2">
		{#if nav.length() >= 1}
			<current.component {...current.props} />
		{/if}
	</div>
</div>

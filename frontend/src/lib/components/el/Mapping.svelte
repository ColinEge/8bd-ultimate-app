<script lang="ts">
    import { onMount } from "svelte";
    import { GetProfileLayout, GetProfileName } from "$lib/wailsjs/go/app/App";
    import { getKeyboardState } from "$lib/state/keyboard-state.svelte";
    import { getKeySelectionState } from "$lib/state/selection-state.svelte";
    import type { service } from "$lib/wailsjs/go/models";
    import ActionSheet from "../ActionSheet.svelte";
    import Keyboard from "../Keyboard.svelte";

    let { title }: { title: string } = $props();
    const keyboard = getKeyboardState();
    const selection = getKeySelectionState();
    let keys: Record<string, service.Key> = $state({});
    let active = $state("");
    let error = $state("");
    let loading = $state(true);

    async function refresh() {
        loading = true;
        error = "";
        selection.deselect();
        try {
            const [layout, name] = await Promise.all([GetProfileLayout(title), GetProfileName()]);
            keys = layout;
            active = name;
        } catch (e) { error = String(e); }
        finally { loading = false; }
    }
    onMount(() => { void refresh(); return () => selection.deselect(); });
</script>

<section class="flex flex-col gap-3 p-6">
    <h1 class="text-2xl">{title} {#if active === title}<span class="text-sm text-primary">Active</span>{/if}</h1>
    <p>{active === title
        ? "Changes save to this profile and are also applied to the keyboard."
        : "Changes save to this profile only. The keyboard will not be changed."}</p>
    {#if error}<p role="alert">{error}</p>{/if}
    <button class="button button-secondary" disabled={loading} onclick={refresh}>Refresh profile</button>
    {#if loading}<p>Loading profile...</p>
    {:else}
        <Keyboard {keys} hasKeypad={keyboard.connection?.HasKeypad} />
        <ActionSheet visible={selection.selectedKey !== null} profileTitle={title} profileKeys={keys}
            onUpdate={(layout) => { keys = layout; }} />
    {/if}
</section>

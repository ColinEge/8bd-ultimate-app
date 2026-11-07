<script lang="ts">
    import { ScanDevices } from "$lib/wailsjs/go/app/App";
    import type { service } from "$lib/wailsjs/go/models";
    import { getNavigationState } from "$lib/state/navigation-state.svelte";
    import { getKeyboardState } from "$lib/state/keyboard-state.svelte";
    import { getKeySelectionState } from "$lib/state/selection-state.svelte";
    import Profiles from "./Profiles.svelte";

    const keyboardState = getKeyboardState();
    const navigationState = getNavigationState();
    const selectionState = getKeySelectionState();

    let devicesPromise: Promise<service.Device[]> | null = $state(null);
    let error = $state("");
    let selecting = $state(false);

    function scanDevices() {
        devicesPromise = ScanDevices();
    }

    async function selectDevice(device: service.Device) {
        selecting = true;
        error = "";
        try {
            await keyboardState.select(device.DeviceID);
            selectionState.deselect();
            navigationState.push({ component: Profiles });
        } catch (e) {
            error = String(e);
        } finally {
            selecting = false;
        }
    }

    $effect(() => {
        scanDevices();
    });
</script>

{#if error}<p role="alert">{error}</p>{/if}
{#await devicesPromise}
    <p>Scanning for keyboards...</p>
{:then devices}
    {#if devices && devices.length > 0}
        <p>Found some devices, select one to interact with</p>
        {#each devices as device}
            <button
                onclick={() => selectDevice(device)}
                disabled={selecting}
                class="button button-primary">{device.ProductName} ({device.DeviceID})</button
            >
        {/each}
    {:else}
        <p>No keyboards found. Connect via USB with the wireless selector off, then rescan.</p>
    {/if}
    <button class="button button-secondary" onclick={scanDevices} disabled={selecting}>Rescan</button>
{:catch e}
    <p>{e}</p>
    <button class="button button-secondary" onclick={scanDevices}>Rescan</button
    >
{/await}

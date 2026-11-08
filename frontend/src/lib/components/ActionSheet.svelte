<script lang="ts">
  import { fly } from "svelte/transition";
  import { GetActions, UpdateProfileBinding } from "$lib/wailsjs/go/app/App";
  import { getKeySelectionState } from "$lib/state/selection-state.svelte";
  import type { service } from "$lib/wailsjs/go/models";
  import { onMount } from "svelte";

  let { visible = false, profileTitle, profileKeys, onUpdate }: {
    visible: boolean;
    profileTitle: string;
    profileKeys: Record<string, service.Key>;
    onUpdate: (keys: Record<string, service.Key>) => void;
  } = $props();
  const selection = getKeySelectionState();
  const selectedID = $derived(selection.selectedKey?.physicalKey);
  const currentKey = $derived(selection.selectedKey
    ? profileKeys[selection.selectedKey.id] : null);

  let options: service.ActionCategory[] = $state([]);
  let category = $state("");
  let action = $state("");
  let error = $state("");
  let busy = $state(false);
  let progress = $state("");
  const actions = $derived(options.find((entry) => entry.Name === category)?.Actions ?? []);

  onMount(() => {
    GetActions().then((result) => {
      options = result;
      category = result[0]?.Name ?? "";
    }).catch((e: unknown) => { error = String(e); });
  });

  async function change(operation: "bind" | "disable" | "reset-to-default") {
    if (selectedID === undefined || busy) return;
    const id = selectedID;
    busy = true;
    progress = `${operation === "bind" ? "Rebinding" : operation === "disable" ? "Disabling" : "Resetting to default:"} ${currentKey?.keyLabel ?? "key"}... Saving profile changes.`;
    error = "";
    try {
      const layout = await UpdateProfileBinding(profileTitle, id,
        operation === "disable" ? "07-00-00" : action, operation === "reset-to-default");
      onUpdate(layout);
    } catch (e) {
      error = String(e);
    } finally {
      busy = false;
    }
  }

  function close() {
    selection.deselect();
    error = "";
  }
</script>

{#if visible}
  <div
    class="fixed bottom-0 left-0 right-0 h-fit bg-zinc-800 text-white rounded-t-2xl pt-10 p-4 z-50"
    aria-busy={busy}
    transition:fly={{ y: 100, duration: 200 }}
  >
    <button class="absolute top-3 right-4 text-xl" aria-label="Close" onclick={close} disabled={busy}>
      &times;
    </button>
    <p>Key: {currentKey?.keyLabel}</p>
    <p>Binding: {currentKey?.actionLabel}</p>
    {#if error}<p role="alert">{error}</p>{/if}
    {#if busy}
      <div role="status" class="flex items-center gap-3 mt-4 rounded-lg bg-zinc-700 p-3">
        <span aria-hidden="true" class="h-5 w-5 rounded-full border-2 border-white/30 border-t-white motion-safe:animate-spin"></span>
        <span>{progress}</span>
      </div>
    {/if}
    <div class="flex flex-row gap-4 mt-6">
      <label for="categories">Categories</label>
      <select class="select" id="categories" bind:value={category} onchange={() => { action = ""; }} disabled={busy}>
        {#each options as option}<option value={option.Name}>{option.Name}</option>{/each}
      </select>
      <label for="actions">Actions</label>
      <select class="select" id="actions" bind:value={action} disabled={busy}>
        <option value="">Select an action</option>
        {#each actions as entry}<option value={entry.actionCode}>{entry.actionLabel}</option>{/each}
      </select>
      <button class="button button-primary" onclick={() => change("bind")} disabled={busy || !action}>Rebind</button>
      <button class="button button-destructive" onclick={() => change("disable")} disabled={busy}>Disable</button>
      <button class="button button-secondary" onclick={() => change("reset-to-default")} disabled={busy}>Reset to default</button>
    </div>
  </div>
{/if}

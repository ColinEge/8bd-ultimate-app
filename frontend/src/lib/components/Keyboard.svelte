<script lang="ts">
  import { service } from "$lib/wailsjs/go/models";
  import Key from "./Key.svelte";
  import { getKeySelectionState } from "$lib/state/selection-state.svelte";
  import {
    fRow,
    typingRows,
    navTop,
    navMain,
    arrows,
    numpadTop,
    numpad,
  } from "./KeyboardLayout";
  
  interface Props {
    keys: Record<string,service.Key>;
    hasKeypad?: boolean;
  }
  
  let { keys, hasKeypad = false }: Props = $props();
  const selectionState = getKeySelectionState();
</script>

<!-- KEYBOARD WRAPPER -->
{#if keys}
  <div class="flex gap-6 p-6 bg-card rounded-xl mx-auto w-fit">
    <!-- LEFT -->
    <div class="flex flex-col gap-4">
      <!-- F ROW -->
      <div class="grid grid-cols-13 gap-1.5">
        {#each fRow as key}
          <Key
            {...key}
            key={keys[key.id]}
            onclick={() => selectionState.select(keys[key.id])}
          ></Key>
        {/each}
      </div>

      <!-- MAIN TYPING AREA -->
      <div class="flex flex-col gap-1.5">
        {#each typingRows as row}
          <div class="grid grid-cols-180 gap-1.5">
            {#each row as key}
              <Key
                {...key}
                key={keys[key.id]}
                onclick={() => selectionState.select(keys[key.id])}
                colSpan={key.colSpan ? key.colSpan * 12 : 12}
              ></Key>
            {/each}
          </div>
        {/each}
      </div>
    </div>

    <!-- MIDDLE -->
    <div class="flex flex-col gap-4">
      <!-- TOP CLUSTER -->
      <div class="flex flex-col gap-1.5">
        {#each navTop as row}
          <div class="grid grid-cols-3 gap-1.5">
            {#each row as key}
              <Key
                {...key}
                key={keys[key.id]}
                onclick={() => selectionState.select(keys[key.id])}
              ></Key>
            {/each}
          </div>
        {/each}
      </div>

      <!-- NAV CLUSTER -->
      <div class="flex flex-col gap-1.5">
        {#each navMain as row}
          <div class="grid grid-cols-3 gap-1.5">
            {#each row as key}
              <Key
                {...key}
                key={keys[key.id]}
                onclick={() => selectionState.select(keys[key.id])}
              ></Key>
            {/each}
          </div>
        {/each}
      </div>

      <!-- ARROWS -->
      <div class="flex flex-col gap-1.5 mt-11">
        {#each arrows as row}
          <div class="grid grid-cols-3 gap-1.5">
            {#each row as key}
              <Key
                {...key}
                key={keys[key.id]}
                onclick={() => selectionState.select(keys[key.id])}
              ></Key>
            {/each}
          </div>
        {/each}
      </div>
    </div>

    <!-- RIGHT (NUMPAD) -->
    {#if hasKeypad}
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-1.5">
          {#each numpadTop as row}
            <div class="grid grid-cols-4 gap-1.5">
              {#each row as key}
                <Key
                  {...key}
                  key={keys[key.id]}
                  onclick={() => selectionState.select(keys[key.id])}
                ></Key>
              {/each}
            </div>
          {/each}
        </div>

        <div class="flex flex-col gap-1.5">
          <div class="grid gap-1.5 grid-cols-4">
            {#each numpad as key}
              <Key
                {...key}
                key={keys[key.id]}
                onclick={() => selectionState.select(keys[key.id])}
              ></Key>
            {/each}
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}

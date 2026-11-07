<script module>
  import type { service } from "$lib/wailsjs/go/models";

  export interface Props {
    key?: service.Key;
    colSpan?: number;
    rowSpan?: number;
    onclick?: () => void;
  }
</script>

<script lang="ts">
  let {
    key,
    colSpan = 1, // grid columns, fractional units
    rowSpan = 1, // grid rows, fractional units
    onclick = () => {},
  }: Props = $props();

  let color = "bg-secondary hover:bg-secondary/20";
  let remappedColor = "bg-primary hover:bg-primary/20";
</script>

{#if key !== undefined}
  <button
    {onclick}
    id={key.keyLabel}
    title={`${key.bindingState}: ${key.actionLabel}`}
    style="grid-column: span {colSpan}; grid-row: span {rowSpan};"
    class="min-h-12 min-w-12 flex items-center justify-center text-sm text-white rounded-md shadow {key.bindingState !== "default"
      ? remappedColor
      : color}"
  >
    {key.keyLabel}
  </button>
{:else}
  <div
    style="grid-column: span {colSpan}; grid-row: span {rowSpan};"
    class="min-h-12 min-w-12 flex"
  ></div>
{/if}



<!-- <button
  onclick={onclick}
  {id}
  title={usage ? `Remapped: ${usage.Desc}` : label}
  style="grid-column: span {colSpan}; grid-row: span {rowSpan};"
  class="min-h-12 min-w-12 flex items-center justify-center text-sm
 {transparent
    ? ''
    : 'text-white rounded-md shadow ' + (usage ? remappedColor : color)}"
>
  {label}
</button> -->

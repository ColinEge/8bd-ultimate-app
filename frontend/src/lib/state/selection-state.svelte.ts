import type { service } from "$lib/wailsjs/go/models";
import { getContext, setContext } from "svelte";

export class KeySelectionState {
    selectedKey = $state<service.Key | null>(null);

    constructor() { }

    select(key: service.Key) {
        console.log("Selected", key)
        this.selectedKey = key;
    }

    deselect() {
        this.selectedKey = null;
    }

    isSelected(keyID: string) {
        return this.selectedKey?.id === keyID;
    }
}

const SELECTION_KEY = Symbol('KEY_SELECTION')

export function setKeySelectionState() {
    return setContext(SELECTION_KEY, new KeySelectionState());
}

export function getKeySelectionState() {
    return getContext<ReturnType<typeof setKeySelectionState>>(SELECTION_KEY);
}

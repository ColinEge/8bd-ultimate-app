import type { service } from "$lib/wailsjs/go/models";
import { getContext, setContext } from "svelte";
import { SelectDevice, GetSelectedDevice } from "$lib/wailsjs/go/app/App";

export class KeyboardState {
    connection = $state<service.Keyboard | null>(null);

    async select(deviceID: string): Promise<void> {
        this.connection = await SelectDevice(deviceID);
    }

    async refresh(): Promise<void> {
        this.connection = await GetSelectedDevice();
    }

    updateKey(key: service.Key) {
        if (this.connection) {
            this.connection.Keys[key.id] = key;
        }
    }

    unload() {
        this.connection = null;
    }
}

const KEYBOARD_KEY = Symbol('KEYBOARD');

export function setKeyboardState() {
    return setContext(KEYBOARD_KEY, new KeyboardState());
}

export function getKeyboardState() {
    return getContext<KeyboardState>(KEYBOARD_KEY);
}

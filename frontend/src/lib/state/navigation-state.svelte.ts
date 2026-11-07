import { getContext, setContext } from "svelte";

export type NavigationEntry = {
  component: any;
  props?: Record<string, any>;
};

export class NavigationState {
    stack = $state<NavigationEntry[]>([]);

    constructor() { }

    length(): number {
        return this.stack.length
    }

    peek(): NavigationEntry {
        return this.stack[this.length() - 1]
    }

    push(nav: NavigationEntry) {
        this.stack = [...this.stack, nav]
    }

    pop() {
        this.stack = this.stack.slice(0, -1);
    }

    reset(entry: NavigationEntry) {
        this.stack = [entry]
    }
}

const NAVIGATION_KEY = Symbol('KEY_NAVIGATION')

export function setNavigationState() {
    return setContext(NAVIGATION_KEY, new NavigationState());
}

export function getNavigationState() {
    return getContext<ReturnType<typeof setNavigationState>>(NAVIGATION_KEY);
}

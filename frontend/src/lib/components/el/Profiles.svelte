<script lang="ts">
    import { onMount } from "svelte";
    import {
        ListProfiles, CreateProfile, RenameProfile, DeleteProfile, CaptureProfile,
        ApplyProfile, GetProfileDirectory, ChooseProfileDirectory, ResetProfileDirectory,
        GetProfileName, ImportSharedProfiles,
    } from "$lib/wailsjs/go/app/App";
    import type { service } from "$lib/wailsjs/go/models";
    import { getKeyboardState } from "$lib/state/keyboard-state.svelte";
    import { getNavigationState } from "$lib/state/navigation-state.svelte";
    import Mapping from "./Mapping.svelte";

    const keyboard = getKeyboardState();
    const navigation = getNavigationState();
    let profiles: service.Profile[] = $state([]);
    let profilesPromise: Promise<void> | null = $state(null);
    let active = $state("");
    let directory = $state("");
    let busy = $state(false);
    let error = $state("");
    let menu = $state("");
    let renaming = $state("");
    let newTitle = $state("");
    let deleting = $state("");
    let creating = $state(false);
    let title = $state("");
    const ordered = $derived([...profiles].sort((a, b) =>
        Number(b.title === active) - Number(a.title === active) || a.title.localeCompare(b.title)));

    function load(): Promise<void> {
        profilesPromise = Promise.all([ListProfiles(), GetProfileName(), GetProfileDirectory()])
            .then(([list, name, path]) => {
                profiles = list;
                active = name;
                directory = path;
            });
        return profilesPromise;
    }

    async function run(operation: () => Promise<void>) {
        if (busy) return;
        busy = true;
        error = "";
        menu = "";
        try { await operation(); }
        catch (e) { error = String(e); }
        finally { busy = false; }
    }

    onMount(() => { void run(load); });

    async function create(capture: boolean) {
        if (capture) await CaptureProfile(title);
        else await CreateProfile(title);
        title = "";
        creating = false;
        await load();
    }
</script>

{#snippet loadingProfiles()}
    <div role="status" aria-label="Loading profiles" class="flex flex-col gap-3">
        <span class="sr-only">Loading profiles...</span>
        {#each [1, 2, 3] as placeholder (placeholder)}
            <div aria-hidden="true" class="rounded-xl bg-card p-4 flex items-center justify-between gap-4">
                <div class="flex flex-col gap-3 motion-safe:animate-pulse">
                    <div class="h-5 w-40 rounded bg-secondary"></div>
                    <div class="h-4 w-24 rounded bg-secondary"></div>
                </div>
                <div class="h-9 w-9 rounded bg-secondary motion-safe:animate-pulse"></div>
            </div>
        {/each}
    </div>
{/snippet}

<section class="flex flex-col gap-5 p-6 w-full max-w-2xl">
    <div>
        <h1 class="text-2xl">{keyboard.connection?.ProductName}</h1>
        <p class="text-sm">Choose a profile to edit or make active.</p>
    </div>
    {#if error}<p role="alert">{error}</p>{/if}
    <div class="flex gap-2">
        <button class="button button-primary" disabled={busy} onclick={() => { creating = !creating; }}>New profile</button>
        <button class="button button-secondary" disabled={busy} onclick={() => run(load)}>Refresh</button>
    </div>
    {#if creating}
        <fieldset disabled={busy} class="flex flex-col gap-2">
            <label for="profile-title">Profile title</label>
            <input class="input" id="profile-title" bind:value={title} />
            <div class="flex gap-2">
                <button class="button button-primary" disabled={!title.trim()} onclick={() => run(() => create(false))}>Create empty</button>
                <button class="button button-secondary" disabled={!title.trim()} onclick={() => run(() => create(true))}>Copy keyboard bindings</button>
            </div>
        </fieldset>
    {/if}
    {#if profilesPromise === null}
        {@render loadingProfiles()}
    {:else}
    {#await profilesPromise}
        {@render loadingProfiles()}
    {:then}
    {#if active && !profiles.some((p) => p.title === active)}
        <p>The keyboard uses "{active}", which is not saved for this model. Copy its bindings into a new profile to keep them.</p>
    {/if}
    {#if profiles.length === 0}<p>No profiles for this keyboard model yet.</p>{/if}
    <ul class="flex flex-col gap-3">
        {#each ordered as profile (profile.title)}
            <li class="rounded-xl bg-card p-4 flex flex-col gap-3">
                <div class="flex items-center justify-between gap-4">
                    <div>
                        <h2 class="text-lg">{profile.title} {#if profile.title === active}<span class="text-sm text-primary">Active</span>{/if}</h2>
                        <p class="text-sm">{profile.keys.length} rebound key{profile.keys.length === 1 ? "" : "s"}</p>
                    </div>
                    <button class="button button-secondary" aria-label={`Options for ${profile.title}`} aria-expanded={menu === profile.title}
                        disabled={busy} onclick={() => { menu = menu === profile.title ? "" : profile.title; }}>&#8942;</button>
                </div>
                {#if menu === profile.title}
                    <div class="flex gap-2 flex-wrap">
                        <button class="button button-secondary" onclick={() => { renaming = profile.title; newTitle = profile.title; menu = ""; }}>Rename</button>
                        {#if profile.title !== active}
                        <button class="button button-primary" onclick={() => run(async () => {
                            await ApplyProfile(profile.title); await load();
                        })}>Set as active</button>
                        {/if}
                        <button class="button button-secondary" onclick={() => navigation.push({ component: Mapping, props: { title: profile.title } })}>Edit profile</button>
                        <button class="button button-destructive" onclick={() => { deleting = profile.title; menu = ""; }}>Delete</button>
                    </div>
                {/if}
                {#if renaming === profile.title}
                    <fieldset disabled={busy} class="flex gap-2">
                        <label for="rename-title">New title</label>
                        <input class="input" id="rename-title" bind:value={newTitle} />
                        <button class="button button-primary" disabled={!newTitle.trim() || newTitle === profile.title} onclick={() => run(async () => {
                            await RenameProfile(profile.title, newTitle); renaming = ""; await load();
                        })}>Save name</button>
                        <button class="button button-secondary" onclick={() => { renaming = ""; }}>Cancel</button>
                    </fieldset>
                {/if}
                {#if deleting === profile.title}
                    <p>Delete "{profile.title}"? {profile.title === active
                        ? "This also clears the active profile on the keyboard."
                        : "Only the local file will be deleted."}</p>
                    <button class="button button-destructive" disabled={busy} onclick={() => run(async () => {
                        await DeleteProfile(profile.title); deleting = ""; await load();
                    })}>Confirm delete</button>
                    <button class="button button-secondary" disabled={busy} onclick={() => { deleting = ""; }}>Cancel</button>
                {/if}
            </li>
        {/each}
    </ul>
    {:catch}
        <p>Profiles could not be loaded. Use Refresh to retry.</p>
    {/await}
    {/if}
    <details>
        <summary>Profile storage</summary>
        <p class="break-all text-sm">{directory}</p>
        <p class="text-sm">Profiles are separated by keyboard model. Changing folders does not move files.</p>
        <div class="flex gap-2 flex-wrap mt-2">
            <button class="button button-secondary" disabled={busy} onclick={() => run(async () => { await ChooseProfileDirectory(); await load(); })}>Choose folder</button>
            <button class="button button-secondary" disabled={busy} onclick={() => run(async () => { await ResetProfileDirectory(); await load(); })}>Use default</button>
            <button class="button button-secondary" disabled={busy} onclick={() => run(async () => { await ImportSharedProfiles(); await load(); })}>Import old shared profiles</button>
        </div>
    </details>
</section>

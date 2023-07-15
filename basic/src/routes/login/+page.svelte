<script lang="ts">
    import {onMount} from "svelte";
    import {login} from "$lib/api";
    import {goto} from "$app/navigation";

    let userId = "";
    let password = "";

    async function _page_login() {
        try {
            await login(userId, password)
            alert(`Welcome ${userId}`)
            await goto("/")
        } catch (e) {
            alert(e)
        }
    }

    onMount(async () => {
    });
</script>

<div class="container h-full mx-auto flex flex-col justify-center items-center">
    <div class="card p-10 space-y-5 flex flex-col justify-center items-center">
        <input type="text" placeholder="ID" bind:value={userId} class="border-2 border-gray-300 rounded-md p-2 w-96 input" />
        <input type="password" placeholder="Password" bind:value={password} class="border-2 border-gray-300 rounded-md p-2 w-96 input" />
        <button class="btn " on:click={async () => {await _page_login()}}>Login</button>
        <span class="text-sm">Don’t have an account? <a href="/register" class="text-primary-400">Sign Up</a></span>
    </div>
</div>

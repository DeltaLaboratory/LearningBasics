<script lang="ts">
    import {onMount} from "svelte";
    import {goto} from "$app/navigation";
    import {parseJwt} from "$lib/api";

    onMount(async () => {
        let accessToken = localStorage.getItem("agnosco")
        if (accessToken == null) {
            await goto("/login")
        }
        try {
            let claims = parseJwt(accessToken)
            if (claims['exp'] * 1000 < Date.now()) {
                await goto("/login")
            }
            await goto("/articles")
        } catch (e) {
            await goto("/login")
        }
    })
</script>
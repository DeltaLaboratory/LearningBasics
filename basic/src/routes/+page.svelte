<script lang="ts">
    import {onMount} from "svelte";
    import {goto} from "$app/navigation";
    import {parseJwt} from "$lib/api";

    onMount(async () => {
        await goto("/articles")
    })
</script>
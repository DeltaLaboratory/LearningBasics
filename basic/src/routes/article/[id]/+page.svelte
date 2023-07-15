<script lang="ts">
    import {onMount} from "svelte";
    import {get_article} from "$lib/api";

    export let data;
    let title;
    let author;
    let content;
    let created;

    onMount(async () => {
        const resp = (await get_article(data.article_id));
        title = resp['title'];
        author = resp['author'];
        created = resp['created'];
        content = resp['content'];
    })
</script>

<nav class="shadow p-4 flex gap-6 items-center align-middle">
    <a href="/articles"><span class="material-symbols-outlined align-middle">arrow_back</span></a>
    <h1 class="text-2xl font-bold align-middle">Article</h1>
</nav>

<main class="p-4 justify-center items-center flex flex-col">
    <h1 class="text-4xl mb-4">{title}</h1>
    <p class="text-sm text-gray-500 mb-4">{author} | {created}</p>
    <p>{content}</p>
</main>
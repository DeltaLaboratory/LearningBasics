<script lang="ts">
    import {list_articles, sessionAud, sessionState} from "$lib/api";
    import {onMount} from "svelte";

    let loggedIn = false
    let loggedAud = ''

    let articles = []

    onMount(async () => {
        loggedIn = sessionState()
        loggedAud = sessionAud()

        articles = (await list_articles())['articles']
    })
</script>

<nav class="shadow p-4 flex justify-between items-center">
    <a href="/" class="text-2xl font-bold">Basic Basic Basic</a>
    <div>
        {#if loggedIn}
            <span class="font-bold text-xl">{loggedAud}</span>
            <a href="/create" class="btn font-bold">Create Article</a>
            <a href="/logout" class="btn font-bold">Logout</a>
        {:else}
            <a href="/login" class="btn font-bold">Login</a>
            <a href="/register" class="btn font-bold">Register</a>
        {/if}
    </div>
</nav>

<section class="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 w-full">
    {#each articles as article}
        <article class="card shadow-md rounded-md overflow-hidden">
            <div class="p-4">
                <h2 class="text-xl font-bold"><a href={`/article/${article.id}`}>{article.title}</a></h2>
                <p class="text-gray-500 text-sm">{article.author}</p>
                <p class="mt-2 text-gray-600">{article.content}</p>
                <p class="text-gray-500 text-xs mt-2 font-bold">{article.created}</p>
            </div>
        </article>
    {/each}
</section>
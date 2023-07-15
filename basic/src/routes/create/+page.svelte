<script>
    import {create_article} from "$lib/api";
    import {goto} from "$app/navigation";

    let title = '';
    let content = '';

    const submitArticle = async () => {
        try {
            let response = await create_article(title, content)
            await goto(`/article/${response['article']['id']}`)
        } catch (e) {
            console.log(e)
        }
    };
</script>

<nav class="shadow p-4 flex gap-6 items-center align-middle">
    <a href="/articles"><span class="material-symbols-outlined align-middle">arrow_back</span></a>
    <h1 class="text-2xl font-bold align-middle">New Article</h1>
</nav>


<main class="p-4">
    <form on:submit={async () => {await submitArticle()}} class="space-y-4 h-screen">
        <div>
            <label for="title" class="block text-sm font-medium">Title</label>
            <input id="title" bind:value={title} type="text" required class="input mt-2 block w-full rounded-md shadow-sm" />
        </div>
        <div>
            <label for="content" class="block text-sm font-medium">Content</label>
            <textarea id="content" bind:value={content} required rows="16" class="textarea mt-2 block w-full rounded-md shadow-sm"></textarea>
        </div>

        <button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Submit Article</button>
    </form>
</main>

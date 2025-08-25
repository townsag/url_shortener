<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	const buttons = [{name: "Home", href: "/"}, {name: "About", href: "/about/"}]

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex flex-col min-h-screen bg-jumbo-orange">
	<!-- divide-x-1 divide-slate-400 border-slate-400 border-1 -->
	<nav class="bg-slate-100 flex flex-row m-4 rounded-md w-min ">
		{#each buttons as { name, href }}
			<a href="{href}" class="px-4 py-3 hover:bg-slate-200 active:bg-slate-300 rounded-md text-slate-800">{name}</a>
		{/each}
	</nav>

	{@render children?.()}
</div>


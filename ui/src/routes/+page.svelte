<script lang="ts">
    import { isValidUrl } from "$lib/utils/validate";
    import { page } from "$app/state";

    type StatusIcon = "!" | "✓"
    interface Message {
        isVisible: boolean;
        contents: string;
        status?: StatusIcon;
    }

    interface Mapping {
        shortUrlId?: string;
        shortUrl?: string;
        longUrl?: string;
    }

	let longUrl: string = $state('');
    let message: Message = $state({ isVisible: false, contents: "" });
    let mapping: Mapping = $state({});

	async function createMappingPost(longUrl: string): Promise<any> {
		try {
			const response = await fetch('/api/mapping', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ longUrl: longUrl })
			});
            console.log(response)
			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(`HTTP error! status: ${response.status}, message: ${errorData}`);
			}
            const data = await response.json();
            return data;
		} catch (error) {
			console.error('Error creating mapping: ', error);
			throw error;
		}
	}
	
    async function createMapping(): Promise<void> {
        // if the long url in the input field is the same as the one we have just processed, dont make a new mapping
        if (mapping.longUrl == longUrl) {
            return;
        }
        // validate the long url in the input field
        // if it is invalid, write to the warning message and make the warning message visible
        // else make the warning message invisible
        if (!isValidUrl(longUrl)) {
            message.isVisible = true;
            message.contents = "please enter a valid url";
            return;
        }
        // call the create mapping api endpoint
        try {
            console.log("calling create mapping api")
            const result = await createMappingPost(longUrl=longUrl);
            message.isVisible = false;
            console.log(page.url.origin);
            mapping = {
                shortUrlId: result.shortUrlId,
                shortUrl: new URL(`/api/${result.shortUrl}`, page.url.origin).href,
                longUrl: longUrl,
            };
            console.log(`updated generated short url to be: ${mapping.shortUrl}`);
        } catch (error) {
            console.log(error)
            message.isVisible = true;
            message.contents = "There was an error creating your short url";
            return; 
        }
    }
    // TODO: add some logic to make the button inactive while we are processing a request
    // TODO: add a loading animation for while we are processing a request
    // TODO: add some logic to prevent generating the same url twice in this session
    //          - the user can generate multiple mappings to the same long url but use
    //            client side caching to prevent this if possible
</script>

<div class="flex flex-col items-center">
	<h1 class="font-abril text-[178px] font-bold text-white">Jumbo</h1>
	<p class="text-[32px] text-white">Short Urls Big Dreams</p>
    <div class="flex flex-col items-start space-y-4 w-min">
        <form class="flex flex-row space-x-4 pt-2">
            <input
                bind:value={longUrl}
                type="ur"
                placeholder="https://example.com"
                class="focus:outline-jumbo-orange w-80 rounded-md text-slate-800 bg-slate-100 px-3 py-2"
            />
            <button
                class="rounded-md bg-slate-100 px-3 py-2 text-slate-800 hover:bg-slate-200 active:bg-slate-300 h-min"
                onclick={createMapping}
                type="submit"
            >
                Submit
            </button>
        </form>
        {#if message.isVisible}
            <div class="px-3 py-2 bg-slate-100 text-slate-800 rounded-md w-fit flex flex-row space-x-2 items-center">
                <p class="rounded-full border-2 border-jumbo-orange text-jumbo-orange w-7 h-7 text-center">!</p>
                <p>{message.contents}</p>
                <button 
                    onclick={() => { message.isVisible = false; }}
                    class="text-jumbo-orange self-center pl-2"
                >X</button>
            </div>
        {/if}
        {#if mapping.shortUrl}
            <div class="flex flex-row bg-slate-100 rounded-md w-fit space-x-2 px-3 py-2">
                <p>Short Url:</p>
                <a 
                    href={mapping.shortUrl}
                    class="text-slate-800"
                >{mapping.shortUrl}</a>
                <!-- TODO: add a button that allows the user to copy the short url to their clipboard -->
                <!-- https://developer.mozilla.org/en-US/docs/Web/API/Clipboard/write -->
            </div>
        {/if}
    </div>
</div>
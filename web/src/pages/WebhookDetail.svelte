<script lang="ts">
  import { query } from "svelte-apollo";
  import { FIND_WEBHOOK } from "../queries";
  import type { Webhook } from "../queries";

  export let params: { id: string };

  const { id } = params;

  const webhook = query<Webhook>(FIND_WEBHOOK, {
    variables: {
      id,
    },
  });

  $: console.log($webhook);
</script>

TESTING

<!-- markup (zero or more items) goes here -->
<style>
  /* your styles go here */
</style>

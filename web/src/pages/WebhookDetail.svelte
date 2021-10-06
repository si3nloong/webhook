<script lang="ts">
  import dayjs from "dayjs";
  import { query } from "svelte-apollo";
  import { FIND_WEBHOOK } from "../queries";
  import type { Webhook } from "../queries";

  export let params: { id: string };

  const { id } = params;

  const webhook = query<{ webhook: Webhook }>(FIND_WEBHOOK, {
    variables: {
      id,
    },
  });

  $: console.log($webhook);
</script>

<div>
  {#if $webhook.data}
    <header>{$webhook.data.webhook.method} {$webhook.data.webhook.url}</header>
    <div>{$webhook.data.webhook.body}</div>
    <div>{$webhook.data.webhook.timeout}ms (Milliseconds)</div>
    <div>{$webhook.data.webhook.noOfRetries}</div>
    {#each $webhook.data.webhook.headers as item}
      <div>{item.key}={item.value}</div>
    {/each}
    <div>
      {dayjs($webhook.data.webhook.createdAt).format("DD MMM YYYY HH:mmA")}
    </div>
    <div>
      {dayjs($webhook.data.webhook.updatedAt).format("DD MMM YYYY HH:mmA")}
    </div>
    {#each $webhook.data.webhook.attempts as item}
      <code>{item.body}</code>
      <div>{item.createdAt}</div>
    {/each}
  {:else if $webhook.error}
    <div>Webhook not found</div>
  {/if}
</div>

<style lang="scss">
  code {
    display: block;
    padding: 10px;
    background: #dcdcdc;
  }

  header {
    font-weight: 600;
  }
</style>

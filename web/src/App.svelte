<script lang="ts">
  import { ApolloClient, InMemoryCache } from "@apollo/client";
  import { setClient, query } from "svelte-apollo";
  import dayjs from "dayjs";
  import { GET_WEBHOOKS } from "./queries";
  import type { Webhook } from "./queries";

  const cache = new InMemoryCache();

  const client = new ApolloClient({
    // Provide required constructor fields
    cache: cache,
    uri: "http://localhost:8080",
    name: "webhook-client",
    version: "1.0",
    queryDeduplication: false,
    // defaultOptions: {
    //   watchQuery: {
    //     fetchPolicy: "cache-and-network",
    //   },
    // },
  });

  setClient(client);

  const webhooks = query<{ webhooks: { nodes: [Webhook] } }>(GET_WEBHOOKS);

  console.log($webhooks);

  console.log(client);
  export let name: string;
</script>

<main>
  {#if $webhooks.loading}
    <div>Loading...</div>
  {:else if $webhooks.data}
    <table>
      {#each $webhooks.data.webhooks.nodes as item}
        <tr>
          <td>{item.id}</td>
          <td>{item.method}</td>
          <td>{item.url}</td>
          <td>{dayjs(item.createdAt).format("YYYY MMM DD")}</td>
        </tr>
      {/each}
    </table>
  {/if}
</main>

<style>
  main {
    text-align: center;
    padding: 1em;
    max-width: 240px;
    margin: 0 auto;
  }

  h1 {
    color: #ff3e00;
    text-transform: uppercase;
    font-size: 4em;
    font-weight: 100;
  }

  @media (min-width: 640px) {
    main {
      max-width: none;
    }
  }
</style>

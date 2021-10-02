<script lang="ts">
  import { ApolloClient, InMemoryCache } from "@apollo/client";
  import { setClient, query } from "svelte-apollo";
  import { GET_WEBHOOKS } from "./queries";

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

  const books = query(GET_WEBHOOKS);

  console.log($books);

  console.log(client);
  export let name: string;
</script>

<main>
  <h1>Hello {name}!</h1>
  <p>
    Visit the <a href="https://svelte.dev/tutorial">Svelte tutorial</a> to learn
    how to build Svelte apps.
  </p>
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

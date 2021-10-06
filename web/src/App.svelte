<script lang="ts">
  import { ApolloClient, InMemoryCache } from "@apollo/client";
  import { setClient } from "svelte-apollo";
  import { push, Router } from "svelte-history-router";
  import routes from "./routes";

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
</script>

<header id="appbar">
  <span id="appname" on:click={() => push("/")}>Webhook UI</span>
</header>
<main>
  <Router {routes} />
</main>

<style>
  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(html, body) {
    width: 100%;
    height: 100%;
    font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
    font-size: 14px;
  }

  #appbar {
    position: fixed;
    top: 0;
    left: 0;
    background: #fff;
    width: 100%;
    padding: 1rem;
    box-shadow: 0 2px 3px rgba(0, 0, 0, 0.2);
  }

  #appname {
    cursor: pointer;
  }

  main {
    display: flex;
    padding-top: 50px;
    max-width: 100%;
    min-height: 100%;
    margin: 0 auto;
  }
</style>

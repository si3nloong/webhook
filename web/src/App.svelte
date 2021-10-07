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

<!-- <head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link
    href="https://fonts.googleapis.com/css2?family=Bitter&display=swap"
    rel="stylesheet"
  />
</head> -->

<header id="appbar">
  <span id="appname" on:click={() => push("/")}>Webhook UI</span>
</header>
<main>
  <Router {routes} />
</main>

<style lang="scss">
  @import url("https://fonts.googleapis.com/css2?family=Bitter&display=swap");
  @import url("https://fonts.googleapis.com/css2?family=Open+Sans&display=swap");

  :global(*) {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  :global(html, body) {
    position: relative;
    width: 100%;
    height: 100%;
    font-family: "Open Sans", sans-serif;
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

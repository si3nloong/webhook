<script lang="ts">
  import { query } from "svelte-apollo";
  import dayjs from "dayjs";
  import { push } from "svelte-history-router";
  import { GET_WEBHOOKS } from "../queries";
  import type { Webhook } from "../queries";

  const webhooks = query<{ webhooks: { nodes: [Webhook] } }>(GET_WEBHOOKS);

  console.log($webhooks);

  const onFormSubmit = (e: Event) => {
    console.log((<HTMLFormElement>e.currentTarget).elements);
  };
</script>

<section id="left-pane">
  <div>Webhooks</div>
  <form on:submit|preventDefault={onFormSubmit}>
    <div>Method</div>
    <div>
      <select multiple>
        <option>All</option>
        <option>GET</option>
        <option>POST</option>
      </select>
    </div>
    <div>Status</div>
    <div>
      <select multiple>
        <option>All</option>
        <option>Success</option>
        <option>Failed</option>
        <option>Expired</option>
      </select>
    </div>
    <div>Start Time</div>
    <div><input name="startDate" type="date" /></div>
    <div><input name="startTime" type="time" /></div>
    <div>End Time</div>
    <div><input name="endDate" type="date" /></div>
    <div><input name="endTime" type="time" /></div>
    <div>Limit Results</div>
    <div><input name="limit" type="number" pattern="[0-9]+" value="50" /></div>
    <button type="submit">Filter</button>
  </form>
</section>
<section id="right-pane">
  {#if $webhooks.loading}
    <div>Loading...</div>
  {:else if $webhooks.data}
    {#if $webhooks.data.webhooks.nodes.length > 0}
      <table>
        {#each $webhooks.data.webhooks.nodes as item}
          <tr on:click={() => push(`/webhook/${item.id}`)}>
            <td>{item.id}</td>
            <td>{item.method}</td>
            <td>{item.url}</td>
            <td>{item.retries}</td>
            <td>{dayjs(item.createdAt).format("YYYY MMM DD HH:mm:ss A")}</td>
          </tr>
        {/each}
      </table>
    {:else}
      <div>No Record</div>
    {/if}
  {/if}
</section>

<style>
  #left-pane {
    text-align: left;
    border-right: 1px solid red;
    padding: 1rem;
  }

  #right-pane {
    padding: 1rem;
  }
</style>

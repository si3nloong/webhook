<script lang="ts">
  import { query } from "svelte-apollo";
  import dayjs from "dayjs";
  import { push } from "svelte-history-router";
  import { GET_WEBHOOKS } from "../queries";
  import type { GetWebhooks } from "../queries";
  import Status from "../components/Status.svelte";
  import Button from "../components/Button.svelte";

  const webhooks = query<GetWebhooks>(GET_WEBHOOKS);

  console.log($webhooks);

  const onFormSubmit = (e: Event) => {
    console.log((<HTMLFormElement>e.currentTarget).elements);
  };

  const onGoTo = (cursor?: string) => (e: Event) => {
    console.log(cursor);
  };
</script>

<section id="left-pane">
  <form on:submit|preventDefault={onFormSubmit}>
    <section>
      <div>Method</div>
      <div>
        <select multiple>
          <option>All</option>
          <option>GET</option>
          <option>POST</option>
        </select>
      </div>
    </section>
    <section>
      <div>Status</div>
      <div>
        <select multiple>
          <option>All</option>
          <option>Success</option>
          <option>Failed</option>
          <option>Expired</option>
        </select>
      </div>
    </section>
    <section>
      <div>Start Time</div>
      <div><input name="startDate" type="date" /></div>
      <div><input name="startTime" type="time" /></div>
    </section>
    <section>
      <div>End Time</div>
      <div><input name="endDate" type="date" /></div>
      <div><input name="endTime" type="time" /></div>
    </section>
    <div>Limit Results</div>
    <div><input name="limit" type="number" pattern="[0-9]+" value="50" /></div>
    <Button
      type="submit"
      style="margin-top: 1rem; background: #7C3AED; color: #fff; width: 100%"
      >Filter</Button
    >
  </form>
</section>
<section id="right-pane">
  <div id="header">
    <div class="caption">Webhooks</div>
    <div class="option-tabs">
      <span>All</span><!----><span>Succeeded</span><!----><span>Failed</span>
    </div>
  </div>
  {#if $webhooks.loading}
    <div>Loading...</div>
  {:else if $webhooks.data}
    {#if $webhooks.data.webhooks.nodes.length > 0}
      <table class="datatable">
        <tr class="columns">
          <td style="width: 10%">Status</td>
          <td style="width: 20%">Webhook ID</td>
          <td>Method</td>
          <td>URL</td>
          <td>Date</td>
        </tr>
        {#each $webhooks.data.webhooks.nodes as item}
          <tr>
            <td><Status value={item.latestStatusCode.toString()} /></td>
            <td
              ><span class="link" on:click={() => push(`/webhook/${item.id}`)}
                >{item.id}</span
              ></td
            >
            <td>{item.method}</td>
            <td>{item.url}</td>
            <td>{dayjs(item.createdAt).format("DD MMM YYYY, HH:mm:ss A")}</td>
          </tr>
        {/each}
      </table>
      <footer id="footer">
        <div>Result {$webhooks.data.webhooks.totalCount}</div>
        <div>
          <Button
            disabled={!$webhooks.data.webhooks.pageInfo.hasPreviousPage}
            on:click={onGoTo($webhooks.data.webhooks.pageInfo.startCursor)}
            >Previous</Button
          >
          <Button
            disabled={!$webhooks.data.webhooks.pageInfo.hasNextPage}
            on:click={onGoTo($webhooks.data.webhooks.pageInfo.endCursor)}
            >Next</Button
          >
        </div>
      </footer>
    {:else}
      <div>No Record</div>
    {/if}
  {/if}
</section>

<style lang="scss">
  #left-pane {
    text-align: left;
    border-right: 1px solid #dcdcdc;
    background: #f5f5f5;
    padding: 1rem;
    width: 260px;

    section {
      padding: 1rem 0;
    }
  }

  input[type="date"],
  input[type="time"] {
    width: 100%;
  }

  #right-pane {
    padding: 2rem;
    flex-grow: 1;
  }

  #header {
    display: flex;
    justify-content: space-between;
    padding: 0 0 2rem;
    align-items: center;
  }

  .caption {
    font-size: 22px;
    font-weight: 800;
  }

  .option-tabs {
    color: #171717;
    border: none;
    border-radius: 5px;

    span {
      cursor: pointer;
      background: #f5f5f5;
      display: inline-block;
      height: 30px;
      line-height: 30px;
      padding: 0 10px;
    }
  }

  table,
  th,
  td {
    border-collapse: collapse;
    padding: 5px 10px;
  }

  th,
  td {
    white-space: nowrap;
    border-top: 1px solid #dcdcdc;
    border-bottom: 1px solid #dcdcdc;
  }

  tr:hover {
    background: #f5f5f5;
  }

  .datatable {
    border-radius: 6px;
    width: 100%;
    border: none;
    // box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
    //   0 4px 6px -2px rgba(0, 0, 0, 0.05) !important;

    .columns {
      font-weight: 800;
      background: #f5f5f5;
    }
  }

  th {
    background: #f5f5f5;
    border-bottom: 2px solid #dcdcdc;
  }

  .link {
    color: #0969da;
    cursor: pointer;
    font-weight: 600;

    &:hover {
      text-decoration: underline;
    }
  }

  #footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 0;
  }
</style>

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

  let selectedIndex: number[] = [];
  const onClick = (idx: number) => (e: Event) => {
    selectedIndex = [...selectedIndex, idx];
  };
</script>

<div style="padding: 2rem;">
  {#if $webhook.data}
    <h2>{$webhook.data.webhook.method} {$webhook.data.webhook.url}</h2>
    <div>{$webhook.data.webhook.timeout}ms (Milliseconds)</div>
    <!-- <div>{$webhook.data.webhook.noOfRetries}</div> -->
    <table>
      <tr>
        <td>Headers</td>
        <td>
          {#each $webhook.data.webhook.headers as item}
            <div>{item.key} {item.value}</div>
          {/each}
        </td>
      </tr>
      <tr>
        <td>Body</td>
        <td>{$webhook.data.webhook.body}</td>
      </tr>
    </table>
    <div>
      {dayjs($webhook.data.webhook.createdAt).format("DD MMM YYYY HH:mmA")}
    </div>
    <div>
      {dayjs($webhook.data.webhook.updatedAt).format("DD MMM YYYY HH:mmA")}
    </div>
    {#each $webhook.data.webhook.attempts as item, i}
      <div class="attempt">
        <header on:click={onClick(i)}>
          <div>Attempt {i + 1}</div>
          <div>
            ({item.elapsedTime}ms) {dayjs(item.createdAt).format(
              "DD MMM YYYY, HH:mm:ssA"
            )}
          </div>
        </header>
        <div class="content" class:open={selectedIndex.includes(i)}>
          <table>
            <tr>
              <td>Headers</td>
              <td>
                {#each item.headers as item}
                  <div>{item.key} {item.value}</div>
                {/each}
              </td>
            </tr>
            <tr>
              <td>Body</td>
              <td><code>{item.body}</code></td>
            </tr>
          </table>
        </div>
      </div>
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

  .attempt {
    header {
      cursor: pointer;
      display: flex;
      justify-content: space-between;
    }

    .content {
      overflow: hidden;
      height: 0;
      transition: all 0.5s;

      &.open {
        height: auto;
      }
    }
  }
</style>

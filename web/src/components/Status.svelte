<script lang="ts">
  export let value = 0;

  let paths: [string, string] = ["unknown", "Unknown"];
  if (value <= 200) {
    paths = ["succeeded", value.toString()];
  } else if (value >= 500) {
    paths = ["failed", value.toString()];
  }
</script>

<span class="status {paths[0]}">{paths[1]}</span>

<style lang="scss">
  .status {
    font-family: "Bitter", serif;
    display: inline-block;
    color: #1d4ed8;
    background: #dbeafe;
    // width: 50px;
    text-align: center;
    padding: 3px 8px;
    border-radius: 10px;

    &.succeeded {
      color: #065f46;
      background: #d1fae5;
    }

    &.failed {
      color: #991b1b;
      background: #fecaca;
    }
  }
</style>

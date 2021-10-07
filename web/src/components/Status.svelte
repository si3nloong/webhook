<script lang="ts">
  export let value = "";
</script>

<span class="status success">Success</span>

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

    &.success {
      background: #d1fae5;
      color: #065f46;
    }
  }
</style>

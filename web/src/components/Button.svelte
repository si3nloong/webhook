<script lang="ts">
  export let type = "submit";
  export let style = "";
  export let disabled = false;
</script>

<button {type} {disabled} {style} on:click><slot /></button>

<style lang="scss">
  button {
    cursor: pointer;
    display: inline-block;
    height: 30px;
    line-height: 30px;
    background: #f5f5f5;
    color: #171717;
    border: none;
    border-radius: 5px;
    width: 80px;
    box-shadow: 0 0 3px -2px rgba(0, 0, 0, 0.3);

    &:disabled {
      cursor: not-allowed;
      background: #bebebe;
    }
  }
</style>

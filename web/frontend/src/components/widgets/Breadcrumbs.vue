<template>
  <div>
    <v-breadcrumbs :items="this.items" divider=">" class="pl-2 py-1"></v-breadcrumbs>
  </div>
</template>
<script>
  export default {
    props: {
      items:Array,
    },
}
  

</script>
<template>
    <v-card max-width="500" class="mx-auto">
        <v-toolbar
            color="pink"
            dark
        >
            <v-app-bar-nav-icon></v-app-bar-nav-icon>

            <v-toolbar-title>Категории</v-toolbar-title>

            <div class="flex-grow-1"></div>
            
        </v-toolbar>

        <v-list>
          <v-list-item
            v-for="item in items"
            :key="item.Slug"
            :to="getLink(item)"
            @click="$emit('click',item.Slug)"
          >
            <v-list-item-content>
                <v-list-item-title v-text="item.Name"></v-list-item-title>
            </v-list-item-content>

          
          </v-list-item>
        </v-list>
    </v-card>
</template>
<script>
  export default {
    props: {
        items:Array,
    },
    data: function () {
      return {
        
      }
    },
    methods:{
      getLink:function(item){
        return "/category/"+item.Slug;
      }
    }
  }
</script>

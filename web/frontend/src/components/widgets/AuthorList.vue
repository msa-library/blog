<template>
    <v-card ma-2 max-width="500" class="mx-auto">
        <v-toolbar
            color="deep-purple accent-4"
            dark
        >
            <v-app-bar-nav-icon></v-app-bar-nav-icon>

            <v-toolbar-title>Авторы</v-toolbar-title>

            <div class="flex-grow-1"></div>

            
        </v-toolbar>

        <v-list subheader>
            <v-list-item
              v-for="item in items"
              :key="item.title"
              :to="getLink(item)"
              @click="$emit('click',item.Slug)"
            >
            <v-list-item-avatar>
                <v-img :src="item.SrcAvatar"></v-img>
            </v-list-item-avatar>

            <v-list-item-content>
                <v-list-item-title v-text="getFullName(item)"></v-list-item-title>
            </v-list-item-content>

            <v-list-item-icon>
                <v-icon :color="item.active ? 'deep-purple accent-4' : 'grey'">chat_bubble</v-icon>
            </v-list-item-icon>
            </v-list-item>
        </v-list>

        
    </v-card>    
</template>
<script>
  export default {
    props: {
        items:Array,
    },
    methods:{
      getFullName:function(item){
        return item.FirstName+' '+item.LastName;
      },
      getLink:function(item){
        return "/author/"+item.Slug;
      }
    }
    
  }
</script>


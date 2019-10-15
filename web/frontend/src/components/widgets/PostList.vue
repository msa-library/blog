<template>
    <v-layout text-xs-center wrap>
        <v-flex 
          xs12
          md6
          v-for="post in items"
          :key="post.Slug"
          pa-2
        >
        <v-card :to="getLink(post)">
            <v-img
            :src="post.Src"
            class="white--text"
            height="200px"
            gradient="to bottom, rgba(0,0,0,.1), rgba(0,0,0,.5)"
            >
            <v-card-title
                class="fill-height align-end"
                v-text="post.Title"
            ></v-card-title>
            </v-img>

            <v-card-actions>
            <div class="flex-grow-1"></div>

            <v-btn icon>
                <v-icon>mdi-heart</v-icon>
            </v-btn>

            <v-btn icon>
                <v-icon>mdi-bookmark</v-icon>
            </v-btn>

            <v-btn icon>
                <v-icon>mdi-share-variant</v-icon>
            </v-btn>
            </v-card-actions>
        </v-card>
        </v-flex>
    </v-layout>
</template>
<script>
  export default {
    props: {
        items:Array,
    },
    methods:{
        getLink:function(item){
            return "/post/"+item.Slug;
        }
    },
    
  }
</script>

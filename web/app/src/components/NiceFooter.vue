<template>
  <v-footer dark>
    <v-row class="ma-0">
      <v-col class="text-left" cols="4"> </v-col>
      <v-col class="text-center" cols="4">
        {{ new Date().getFullYear() }} — <strong>Recipe Hunter</strong>
      </v-col>
      <v-col class="text-right" cols="4">
        <v-tooltip top>
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              dark
              v-bind="attrs"
              v-on="on"
              elevation="3"
              icon
              tile
              small
              href="https://github.com/reaper47/recipe-hunter"
              target="_blank"
            >
              <v-icon>mdi-github</v-icon>
            </v-btn>
          </template>
          <span>Made with <v-icon color="red">mdi-heart</v-icon></span>
        </v-tooltip>
      </v-col>
    </v-row>
  </v-footer>
</template>

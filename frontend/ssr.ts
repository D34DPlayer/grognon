import { createInertiaApp } from '@inertiajs/vue3'
import createServer from '@inertiajs/vue3/server'
import { renderToString } from '@vue/server-renderer'
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers'
import { createSSRApp, h } from 'vue'
import vuetify from './vuetify'

createServer(page =>
  createInertiaApp({
    page,
    render: renderToString,
    // @ts-expect-error Expected error
    resolve: name => resolvePageComponent(`./pages/${name}.vue`, import.meta.glob('./pages/**/*.vue')),
    setup({ App, props, plugin }) {
      return createSSRApp({ render: () => h(App, props) })
        .use(vuetify)
        .use(plugin)
    },
  }),
)

import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../pages/Dashboard.vue'
import Reportes from '../pages/Reportes.vue'
import Crear from '../pages/Crear.vue'
import Analisis from '../pages/Analisis.vue'
import Agents from '../pages/Agents.vue'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
  },
  {
    path: '/reportes',
    name: 'Reportes',
    component: Reportes,
  },
  {
    path: '/crear',
    name: 'Crear',
    component: Crear,
  },
  {
    path: '/agents',
    name: 'Agents',
    component: Agents,
  },
  {
    path: '/analisis',
    name: 'Analisis',
    component: Analisis,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    redirectTo: '/dashboard',
    pathMatch: 'full'
  },
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.module').then(m => m.AuthModule)
  },
  {
    path: 'dashboard',
    loadChildren: () => import('./features/dashboard/dashboard.module').then(m => m.DashboardModule)
  },
  {
    path: 'api-discovery',
    loadChildren: () => import('./features/api-discovery/api-discovery.module').then(m => m.ApiDiscoveryModule)
  },
  {
    path: 'threat-detection',
    loadChildren: () => import('./features/threat-detection/threat-detection.module').then(m => m.ThreatDetectionModule)
  },
  {
    path: 'data-protection',
    loadChildren: () => import('./features/data-protection/data-protection.module').then(m => m.DataProtectionModule)
  },
  {
    path: 'attack-protection',
    loadChildren: () => import('./features/attack-protection/attack-protection.module').then(m => m.AttackProtectionModule)
  },
  {
    path: 'gateway-integration',
    loadChildren: () => import('./features/gateway-integration/gateway-integration.module').then(m => m.GatewayIntegrationModule)
  },
  {
    path: '**',
    redirectTo: '/dashboard'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

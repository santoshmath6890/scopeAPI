import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { GatewayIntegrationOverviewComponent } from './components/gateway-integration-overview/gateway-integration-overview.component';
import { IntegrationListComponent } from './components/integration-list/integration-list.component';
import { IntegrationFormComponent } from './components/integration-form/integration-form.component';
import { IntegrationDetailsComponent } from './components/integration-details/integration-details.component';
import { KongIntegrationComponent } from './components/kong-integration/kong-integration.component';
import { NginxIntegrationComponent } from './components/nginx-integration/nginx-integration.component';
import { TraefikIntegrationComponent } from './components/traefik-integration/traefik-integration.component';
import { EnvoyIntegrationComponent } from './components/envoy-integration/envoy-integration.component';
import { HAProxyIntegrationComponent } from './components/haproxy-integration/haproxy-integration.component';

const routes: Routes = [
  {
    path: '',
    component: GatewayIntegrationOverviewComponent,
    children: [
      { path: '', redirectTo: 'integrations', pathMatch: 'full' },
      { path: 'integrations', component: IntegrationListComponent },
      { path: 'integrations/new', component: IntegrationFormComponent },
      { path: 'integrations/:id', component: IntegrationDetailsComponent },
      { path: 'integrations/:id/edit', component: IntegrationFormComponent },
      { path: 'kong', component: KongIntegrationComponent },
      { path: 'nginx', component: NginxIntegrationComponent },
      { path: 'traefik', component: TraefikIntegrationComponent },
      { path: 'envoy', component: EnvoyIntegrationComponent },
      { path: 'haproxy', component: HAProxyIntegrationComponent }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class GatewayIntegrationRoutingModule { } 
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { GatewayIntegrationRoutingModule } from './gateway-integration-routing.module';

// Components
import { GatewayIntegrationOverviewComponent } from './components/gateway-integration-overview/gateway-integration-overview.component';
import { IntegrationListComponent } from './components/integration-list/integration-list.component';
import { IntegrationFormComponent } from './components/integration-form/integration-form.component';
import { IntegrationDetailsComponent } from './components/integration-details/integration-details.component';
import { KongIntegrationComponent } from './components/kong-integration/kong-integration.component';
import { NginxIntegrationComponent } from './components/nginx-integration/nginx-integration.component';
import { TraefikIntegrationComponent } from './components/traefik-integration/traefik-integration.component';
import { EnvoyIntegrationComponent } from './components/envoy-integration/envoy-integration.component';
import { HAProxyIntegrationComponent } from './components/haproxy-integration/haproxy-integration.component';

// Services
import { GatewayIntegrationService } from './services/gateway-integration.service';
import { KongService } from './services/kong.service';
import { NginxService } from './services/nginx.service';
import { TraefikService } from './services/traefik.service';
import { EnvoyService } from './services/envoy.service';
import { HAProxyService } from './services/haproxy.service';

// Shared components
import { SharedModule } from '../../shared/shared.module';

@NgModule({
  declarations: [
    GatewayIntegrationOverviewComponent,
    IntegrationListComponent,
    IntegrationFormComponent,
    IntegrationDetailsComponent,
    KongIntegrationComponent,
    NginxIntegrationComponent,
    TraefikIntegrationComponent,
    EnvoyIntegrationComponent,
    HAProxyIntegrationComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    RouterModule,
    GatewayIntegrationRoutingModule,
    SharedModule
  ],
  providers: [
    GatewayIntegrationService,
    KongService,
    NginxService,
    TraefikService,
    EnvoyService,
    HAProxyService
  ]
})
export class GatewayIntegrationModule { } 